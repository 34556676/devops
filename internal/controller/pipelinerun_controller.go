package controller

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	log "sigs.k8s.io/controller-runtime/pkg/log"

	"time"

	devopsv1 "github.com/lsym/devops/api/v1"
)

const (
	WorkerSpacePrefix     = "/internal/workspaces"
	ResultPrefix          = "/internal/results"
	ScriptPrefix          = "/internal/scripts"
	StepPrefix            = "/internal/steps"
	InternalWorkspaceName = "internal-workspaces"
	InternalScriptName    = "internal-scripts"
	InternalResultName    = "internal-results"
	InternalStepName      = "internal-steps"
)

var (
	internalVolume = []corev1.Volume{
		{
			Name:         InternalStepName,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		},
	}
	internalVolumeMount = []corev1.VolumeMount{
		{
			Name:      InternalStepName,
			MountPath: StepPrefix,
			ReadOnly:  true,
		},
	}
)

// PipelineRunReconciler reconciles a PipelineRun object
type PipelineRunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=devops.lsym.org,resources=pipelineruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.lsym.org,resources=pipelineruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=devops.lsym.org,resources=pipelineruns/finalizers,verbs=update

func (r *PipelineRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling PipelineRun", "name", req.Name)

	// 获取 PipelineRun 实例
	pipelineRun := &devopsv1.PipelineRun{}
	if err := r.Get(ctx, req.NamespacedName, pipelineRun); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 如果 PipelineRun 已经完成，则不再处理
	if pipelineRun.Status.OverallStatus == "Succeeded" || pipelineRun.Status.OverallStatus == "Failed" {
		logger.Info("PipelineRun already completed", "status", pipelineRun.Status.OverallStatus)
		return ctrl.Result{}, nil
	}

	// 获取关联的 Pipeline
	pipeline := &devopsv1.Pipeline{}
	if err := r.Get(ctx, types.NamespacedName{Name: pipelineRun.Spec.PipelineRef, Namespace: req.Namespace}, pipeline); err != nil {
		if errors.IsNotFound(err) {
			r.updateStatus(ctx, pipelineRun, "Failed", fmt.Sprintf("Pipeline %s not found", pipelineRun.Spec.PipelineRef))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger.Info("Pipeline found", "Script", pipeline.Spec.Stages[0].Tasks[0].Script)

	// 初始化状态
	if pipelineRun.Status.StageStatuses == nil {
		pipelineRun.Status.StageStatuses = make([]devopsv1.StageStatus, len(pipeline.Spec.Stages))
		for i, stage := range pipeline.Spec.Stages {
			pipelineRun.Status.StageStatuses[i] = devopsv1.StageStatus{
				Name:  stage.Name,
				State: "Pending",
			}
		}
		pipelineRun.Status.OverallStatus = "Pending"
		pipelineRun.Status.StartTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, pipelineRun); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("Initialized PipelineRun status")
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	// 找到下一个需要处理的阶段
	currentStageIndex := -1
	for i, stageStatus := range pipelineRun.Status.StageStatuses {
		if stageStatus.State == "Pending" || stageStatus.State == "Running" {
			currentStageIndex = i
			break
		}
	}

	// 如果没有找到需要处理的阶段，检查是否所有阶段都已完成
	if currentStageIndex == -1 {
		logger.Info("No more stages to process, checking final status")
		// 检查是否有失败的阶段
		for _, stageStatus := range pipelineRun.Status.StageStatuses {
			if stageStatus.State == "Failed" {
				pipelineRun.Status.OverallStatus = "Failed"
				pipelineRun.Status.CompletionTime = &metav1.Time{Time: time.Now()}
				if err := r.Status().Update(ctx, pipelineRun); err != nil {
					return ctrl.Result{}, err
				}
				logger.Info("PipelineRun failed due to stage failure")
				return ctrl.Result{}, nil
			}
		}

		// 所有阶段都成功完成
		pipelineRun.Status.OverallStatus = "Succeeded"
		pipelineRun.Status.CompletionTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, pipelineRun); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("PipelineRun completed successfully")
		return ctrl.Result{}, nil
	}

	// 处理当前阶段
	currentStage := &pipeline.Spec.Stages[currentStageIndex]
	currentStageStatus := &pipelineRun.Status.StageStatuses[currentStageIndex]

	// 更新当前步骤显示
	pipelineRun.Status.CurrentStep = currentStage.Name
	if err := r.Status().Update(ctx, pipelineRun); err != nil {
		return ctrl.Result{}, err
	}

	// 检查前序阶段是否完成（如果不是第一个阶段）
	if currentStageIndex > 0 {
		previousStageStatus := pipelineRun.Status.StageStatuses[currentStageIndex-1]
		if previousStageStatus.State != "Succeeded" {
			if previousStageStatus.State == "Failed" {
				// 前序阶段失败，当前阶段也标记为失败
				currentStageStatus.State = "Failed"
				currentStageStatus.Message = "Previous stage failed"
				r.updateStatus(ctx, pipelineRun, "Failed", "A previous stage failed")
				logger.Info("Stage failed due to previous stage failure", "stage", currentStage.Name)
				return ctrl.Result{}, nil
			} else {
				// 前序阶段未完成，等待
				currentStageStatus.Message = "Waiting for previous stage to complete"
				if err := r.Status().Update(ctx, pipelineRun); err != nil {
					return ctrl.Result{}, err
				}
				logger.Info("Waiting for previous stage to complete", "stage", currentStage.Name)
				return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
			}
		}
	}

	// 如果阶段未开始，则创建 Jobs
	if currentStageStatus.State == "Pending" {
		logger.Info("Starting stage", "stage", currentStage.Name)
		currentStageStatus.State = "Running"
		currentStageStatus.StartTime = &metav1.Time{Time: time.Now()}
		pipelineRun.Status.OverallStatus = "Running"
		pipelineRun.Status.StageStatuses[currentStageIndex].State = "Running"
		if err := r.Status().Update(ctx, pipelineRun); err != nil {
			return ctrl.Result{}, err
		}

		// 创建该阶段的所有任务 Job
		for _, task := range currentStage.Tasks {
			if err := r.createJobForTask(ctx, pipelineRun, currentStage, &task); err != nil {
				currentStageStatus.State = "Failed"
				currentStageStatus.Message = fmt.Sprintf("Failed to create job: %v", err)
				r.updateStatus(ctx, pipelineRun, "Failed", fmt.Sprintf("Stage %s failed: %v", currentStage.Name, err))
				logger.Error(err, "Failed to create job for task", "stage", currentStage.Name, "task", task.Name)
				return ctrl.Result{}, err
			}
		}
		logger.Info("Created all jobs for stage", "stage", currentStage.Name)
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	// 检查该阶段的所有 Job 状态
	logger.Info("Checking stage status", "stage", currentStage.Name)
	stageCompleted := true
	stageFailed := false
	var failedTask string

	for _, task := range currentStage.Tasks {
		jobName := generateJobName(pipelineRun.Name, currentStage.Name, task.Name)
		job := &batchv1.Job{}
		if err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: pipelineRun.Namespace}, job); err != nil {
			if errors.IsNotFound(err) {
				stageCompleted = false
				logger.Info("Job not found", "job", jobName)
				continue
			}
			logger.Error(err, "Failed to get job", "job", jobName)
			return ctrl.Result{}, err
		}

		// 检查 Job 状态 - 修复这里的逻辑
		logger.Info("Job status", "job", jobName, "succeeded", job.Status.Succeeded, "failed", job.Status.Failed, "active", job.Status.Active)

		// Job 成功完成的条件：Succeeded > 0 且 Active == 0
		if isJobCompleted(job) {
			// Job 成功完成
			logger.Info("Job succeeded", "job", jobName)
		} else if isJobFailed(job) {
			// Job 失败
			stageFailed = true
			failedTask = task.Name
			logger.Info("Job failed", "job", jobName)
			break
		} else {
			// Job 仍在运行或等待中
			stageCompleted = false
			logger.Info("Job still running or pending", "job", jobName,
				"succeeded", job.Status.Succeeded, "failed", job.Status.Failed, "active", job.Status.Active)
		}
	}

	if stageFailed {
		currentStageStatus.State = "Failed"
		currentStageStatus.CompletionTime = &metav1.Time{Time: time.Now()}
		currentStageStatus.Message = fmt.Sprintf("Task %s failed", failedTask)
		r.updateStatus(ctx, pipelineRun, "Failed", fmt.Sprintf("Stage %s failed: task %s failed", currentStage.Name, failedTask))
		logger.Info("Stage failed", "stage", currentStage.Name, "task", failedTask)
		return ctrl.Result{}, nil
	} else if stageCompleted {
		currentStageStatus.State = "Succeeded"
		currentStageStatus.CompletionTime = &metav1.Time{Time: time.Now()}
		pipelineRun.Status.StageStatuses[currentStageIndex].State = "Succeeded"
		currentStageStatus.Message = ""
		if err := r.Status().Update(ctx, pipelineRun); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("Stage completed successfully", "stage", currentStage.Name)
		//logger.Info("currentStageIndex", "StageStatuses", pipelineRun.Status.StageStatuses[currentStageIndex+1])
		// 当前阶段完成，立即触发下一次协调来处理下一个阶段
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	// 阶段仍在运行中，等待
	logger.Info("Stage still running", "stage", currentStage.Name)
	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

// 为任务创建 Job
func (r *PipelineRunReconciler) createJobForTask(ctx context.Context, pipelineRun *devopsv1.PipelineRun, stage *devopsv1.Stage, task *devopsv1.Task) error {
	logger := log.FromContext(ctx)
	jobName := generateJobName(pipelineRun.Name, stage.Name, task.Name)

	// 检查 Job 是否已存在
	existingJob := &batchv1.Job{}
	if err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: pipelineRun.Namespace}, existingJob); err == nil {
		// Job 已存在，跳过创建
		logger.Info("Job already exists", "job", jobName)
		return nil
	} else if !errors.IsNotFound(err) {
		logger.Error(err, "Failed to check if job exists", "job", jobName)
		return err
	}

	// 设置环境变量
	var envVars []corev1.EnvVar
	for _, env := range task.Env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	// 添加流水线参数
	for k, v := range pipelineRun.Spec.Params {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "PARAM_" + k,
			Value: v,
		})
	}

	// 100秒自动清理job
	ttlSecondsAfterFinished := int32(100)

	//生成内部的volume
	volumes := make([]corev1.Volume, 0)
	volumes = append(volumes, internalVolume...)
	//生成公共的volumeMount
	commonVolumeMount := make([]corev1.VolumeMount, 0)
	commonVolumeMount = append(commonVolumeMount, internalVolumeMount...)

	//解析workspace产生的volume和volumeMount
	wsVs, wsVms := createWorkspaceVolumesAndMount(pipelineRun.Spec.Workspaces)

	volumes = append(volumes, wsVs...)
	commonVolumeMount = append(commonVolumeMount, wsVms...)

	newCmd := ""
	if task.Script != "" {
		configMapName := fmt.Sprintf("pipelinerun-%s-step-%s-script", stage.Name, task.Name)
		err := r.Get(ctx, client.ObjectKey{Namespace: pipelineRun.Namespace, Name: configMapName}, &corev1.ConfigMap{})
		if errors.IsNotFound(err) {
			// ConfigMap 不存在，可以创建它
			logger.Info("ConfigMap 不存在，开始创建", "ConfigMap", configMapName)

			newConfigMap := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: pipelineRun.Namespace,
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(
							pipelineRun,
							pipelineRun.GroupVersionKind(),
						),
					},
				},
				Data: map[string]string{
					"script": task.Script,
				},
			}
			//不存在则创建
			if err := r.Create(ctx, newConfigMap); err != nil {
				logger.Error(err, "Failed to create ConfigMap", "configMapName", configMapName)
				return err
			}

		} else if err != nil {
			logger.Error(err, "Failed to create ConfigMap", "configMapName", configMapName)
			return err
		}

		var modNum *int32
		var modNumValue int32 = 0777
		modNum = &modNumValue
		cm := &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: configMapName,
			},
			DefaultMode: modNum,
		}

		tmp := corev1.Volume{
			Name:         configMapName,
			VolumeSource: corev1.VolumeSource{ConfigMap: cm}}
		volumes = append(volumes, tmp)

		scriptMountPath := fmt.Sprintf("%s/%s", ScriptPrefix, configMapName)
		scVolumeMount := corev1.VolumeMount{
			Name:      configMapName,
			MountPath: scriptMountPath,
		}
		commonVolumeMount = append(commonVolumeMount, scVolumeMount)
		newCmd = fmt.Sprintf("ls -l %s/script; cat %s/script; bash %s/script", scriptMountPath, scriptMountPath, scriptMountPath)
	}

	// 创建 Job 对象
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: pipelineRun.Namespace,
			Labels: map[string]string{
				"pipelinerun": pipelineRun.Name,
				"pipeline":    pipelineRun.Spec.PipelineRef,
				"stage":       stage.Name,
				"task":        task.Name,
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:         "task",
							Image:        task.Image,
							Command:      []string{"bash", "-c", newCmd},
							Env:          envVars,
							VolumeMounts: commonVolumeMount,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	// 设置 OwnerReference
	if err := ctrl.SetControllerReference(pipelineRun, job, r.Scheme); err != nil {
		logger.Error(err, "Failed to set controller reference", "job", jobName)
		return err
	}

	// 创建 Job
	if err := r.Create(ctx, job); err != nil {
		logger.Error(err, "Failed to create Job", "job", jobName)
		return err
	}

	logger.Info("Created Job for task", "pipelineRun", pipelineRun.Name, "stage", stage.Name, "task", task.Name, "job", jobName)
	return nil
}

// 更新 PipelineRun 状态
func (r *PipelineRunReconciler) updateStatus(ctx context.Context, pipelineRun *devopsv1.PipelineRun, status, message string) error {
	pipelineRun.Status.OverallStatus = status
	if status == "Failed" || status == "Succeeded" {
		pipelineRun.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	}
	return r.Status().Update(ctx, pipelineRun)
}

// 生成 Job 名称
func generateJobName(pipelineRun, stage, task string) string {
	return fmt.Sprintf("%s-%s-%s", pipelineRun, stage, task)
}

func isJobCompleted(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobFailed(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func createWorkspaceVolumesAndMount(wb []devopsv1.WorkspaceBinding) (vs []corev1.Volume, vms []corev1.VolumeMount) {
	for _, w := range wb {
		w := w
		vName := fmt.Sprintf("%s-%s", InternalWorkspaceName, w.Name)
		tmpv := corev1.Volume{
			Name: vName,
		}
		tmpvm := corev1.VolumeMount{
			Name:      vName,
			MountPath: fmt.Sprintf("%s-%s", WorkerSpacePrefix, w.Name),
			SubPath:   w.SubPath,
		}

		switch {
		case w.PersistentVolumeClaim != nil:
			pvc := *w.PersistentVolumeClaim
			tmpv.VolumeSource = corev1.VolumeSource{PersistentVolumeClaim: &pvc}
		case w.EmptyDir != nil:
			ed := *w.EmptyDir
			tmpv.VolumeSource = corev1.VolumeSource{EmptyDir: &ed}
		case w.ConfigMap != nil:
			cm := *w.ConfigMap
			tmpv.VolumeSource = corev1.VolumeSource{ConfigMap: &cm}
		case w.Secret != nil:
			s := *w.Secret
			tmpv.VolumeSource = corev1.VolumeSource{Secret: &s}
		}
		vs = append(vs, tmpv)
		vms = append(vms, tmpvm)

	}
	return vs, vms
}

// SetupWithManager 设置控制器管理器
func (r *PipelineRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsv1.PipelineRun{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
