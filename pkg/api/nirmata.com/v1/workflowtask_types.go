/*
Copyright 2021.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	StatePending = "PENDING"
	StateExecuting = "EXECUTING"
	StateCompleted = "COMPLETED"
	StateFailed = "FAILED"
)

// WorkflowTask is the Schema for the workflowtasks API
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type WorkflowTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkflowTaskSpec   `json:"spec,omitempty"`
	Status WorkflowTaskStatus `json:"status,omitempty"`
}

// WorkflowTaskSpec defines the desired state of WorkflowTask
type WorkflowTaskSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type string `json:"type,omitempty"`
}

// WorkflowTaskStatus defines the observed state of WorkflowTask
type WorkflowTaskStatus struct {
	State string `json:"state,omitempty"`
	StartTimeUTC string `json:"startTimeUTC,omitempty"`
	CompletionTimeUTC string `json:"completionTimeUTC,omitempty"`
	Executor string `json:"executor,omitempty"`
	Error string `json:"error,omitempty"`
}

// WorkflowTaskList contains a list of WorkflowTask
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type WorkflowTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkflowTask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkflowTask{}, &WorkflowTaskList{})
}
