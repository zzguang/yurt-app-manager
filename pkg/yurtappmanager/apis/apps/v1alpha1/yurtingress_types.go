/*
Copyright 2021 The OpenYurt Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Define the default nodepool ingress related values
const (
	// DefaultIngressControllerReplicasPerPool defines the default ingress controller replicas per pool
	DefaultIngressControllerReplicasPerPool int32 = 1
	// NginxIngressControllerVersion defines the nginx ingress controller version
	NginxIngressControllerVersion = "0.48.1"
	// SingletonYurtIngressInstanceName defines the singleton instance name of YurtIngress
	SingletonYurtIngressInstanceName = "yurtingress-singleton"
	// YurtIngressFinalizer is used to cleanup ingress resources when singleton YurtIngress CR is deleted
	YurtIngressFinalizer = "ingress.operator.openyurt.io"
)

// YurtIngressSpec defines the desired state of YurtIngress
type YurtIngressSpec struct {
	// Indicates the number of the ingress controllers to be deployed under all the specified nodepools.
	// +optional
	Replicas int32 `json:"ingress_controller_replicas_per_pool,omitempty"`

	// Indicates all the nodepools on which to enable ingress.
	// +optional
	Pools []string `json:"pools,omitempty"`
}

// YurtIngressStatus defines the observed state of YurtIngress
type YurtIngressStatus struct {
	// Indicates the number of the ingress controllers deployed under all the specified nodepools.
	// +optional
	Replicas int32 `json:"ingress_controller_replicas_per_pool,omitempty"`

	// Indicates all the nodepools on which to enable ingress.
	// +optional
	Pools []string `json:"pools,omitempty"`

	// Indicates the nginx ingress controller version deployed under all the specified nodepools.
	// +optional
	Version string `json:"nginx_ingress_controller_version,omitempty"`

	// Total number of ready pools on which ingress is enabled.
	// +optional
	ReadyNum int32 `json:"readyNum"`

	// Total number of unready pools on which ingress is enabling or enable failed.
	// +optional
	UnreadyNum int32 `json:"unreadyNum"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=yurtingresses,shortName=ying,categories=all
// +kubebuilder:printcolumn:name="Nginx-Ingress-Version",type="string",JSONPath=".status.nginx_ingress_controller_version",description="The nginx ingress controller version"
// +kubebuilder:printcolumn:name="Replicas-Per-Pool",type="integer",JSONPath=".status.ingress_controller_replicas_per_pool",description="The nginx ingress controller replicas per pool"
// +kubebuilder:printcolumn:name="ReadyNum",type="integer",JSONPath=".status.readyNum",description="The number of pools on which ingress is enabled"
// +kubebuilder:printcolumn:name="NotReadyNum",type="integer",JSONPath=".status.unreadyNum",description="The number of pools on which ingress is enabling or enable failed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +genclient:nonNamespaced

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// YurtIngress is the Schema for the yurtingresses API
type YurtIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   YurtIngressSpec   `json:"spec,omitempty"`
	Status YurtIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// YurtIngressList contains a list of YurtIngress
type YurtIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []YurtIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&YurtIngress{}, &YurtIngressList{})
}
