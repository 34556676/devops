/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineRunSpec 定义 PipelineRun 的期望状态
type PipelineRunSpec struct {
	// +kubebuilder:validation:Required
	PipelineRef string `json:"pipelineRef"`

	// +kubebuilder:validation:Optional
	Params map[string]string `json:"params,omitempty"`

	Workspaces []WorkspaceBinding `json:"workspaces,omitempty"`
}

// StageStatus 定义阶段的执行状态
type StageStatus struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
	State string `json:"state"`

	// +kubebuilder:validation:Optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// +kubebuilder:validation:Optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

// PipelineRunStatus 定义 PipelineRun 的观察状态
type PipelineRunStatus struct {
	// +kubebuilder:validation:Optional
	StageStatuses []StageStatus `json:"stageStatuses,omitempty"`

	// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
	OverallStatus string `json:"overallStatus,omitempty"`

	// +kubebuilder:validation:Optional
	CurrentStep string `json:"currentStep,omitempty"`

	// +kubebuilder:validation:Optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// +kubebuilder:validation:Optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CurrentStep",type=string,JSONPath=".status.currentStep",description="currentStep",priority=0
// +kubebuilder:printcolumn:name="OverallStatus",type=string,JSONPath=".status.overallStatus",description="The overall status of the PipelineRun"
// +kubebuilder:printcolumn:name="CompletionTime",type=date,JSONPath=".status.completionTime",description="The completion time of the PipelineRun"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// PipelineRun is the Schema for the pipelineruns API
type PipelineRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineRunSpec   `json:"spec,omitempty"`
	Status PipelineRunStatus `json:"status,omitempty"`
}

type WorkspaceBinding struct {
	Name                  string                                    `json:"name"`
	SubPath               string                                    `json:"subPath,omitempty"`
	VolumeClaimTemplate   *corev1.PersistentVolumeClaim             `json:"volumeClaimTemplate,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
	EmptyDir              *corev1.EmptyDirVolumeSource              `json:"emptyDir,omitempty"`
	ConfigMap             *corev1.ConfigMapVolumeSource             `json:"configMap,omitempty"`
	Secret                *corev1.SecretVolumeSource                `json:"secret,omitempty"`
}

// +kubebuilder:object:root=true
// PipelineRunList contains a list of PipelineRun.
type PipelineRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineRun{}, &PipelineRunList{})
}
