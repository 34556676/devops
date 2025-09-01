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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PipelineSpec 定义流水线的期望状态
type PipelineSpec struct {
	// +kubebuilder:validation:Required
	Stages     []Stage                `json:"stages,omitempty"`     //步骤
	Workspaces []WorkspaceDeclaration `json:"workspaces,omitempty"` //存储
}

// Stage 定义流水线的一个阶段
type Stage struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:default=false
	Parallel bool `json:"parallel,omitempty"`

	// +kubebuilder:validation:Required
	Tasks []Task `json:"tasks"`
}

// Task 定义阶段中的一个任务
type Task struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Optional
	Script string `json:"script,omitempty"`

	// +kubebuilder:validation:Optional
	Env []EnvVar `json:"env,omitempty"`
}

// EnvVar 定义环境变量
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// workspace代表存储配置
type WorkspaceDeclaration struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MountPath   string `json:"mountPath,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty"` //是否只读
	Optional    bool   `json:"optional,omitempty"` //是否必填，相当于pipelineRun实例化时，是否必须要填写
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineSpec defines the desired state of Pipeline.

// PipelineStatus 定义流水线的观察状态
type PipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API.
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline.
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
