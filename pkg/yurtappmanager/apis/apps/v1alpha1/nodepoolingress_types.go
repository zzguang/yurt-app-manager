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
	// SingletonNodePoolIngressInstanceName defines the singleton instance name of NodePoolIngress
	SingletonNodePoolIngressInstanceName = "nodepool-ingress"
	// NodePoolNameSpacePrefix defines the prefix of the nodepool namespace for nginx ingress
	NodePoolNameSpacePrefix = "nodepool"
)

// NodePoolIngressSpec defines the desired state of NodePoolIngress
type NodePoolIngressSpec struct {
	// Indicates the number of the ingress controllers to be deployed under all the specified nodepools.
	// +optional
	Replicas int32 `json:"ingress_controller_replicas_per_pool,omitempty"`

	// Indicates all the nodepools on which to enable ingress.
	// +optional
	Pools []string `json:"pools,omitempty"`
}

// NodePoolIngressStatus defines the observed state of NodePoolIngress
type NodePoolIngressStatus struct {
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
	ReadyPoolNum int32 `json:"readyPoolNum"`

	// Total number of unready pools on which ingress is enabling or enable failed.
	// +optional
	UnreadyPoolNum int32 `json:"unreadyPoolNum"`
}

// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=nping
// +kubebuilder:printcolumn:name="ReadyPools",type="integer",JSONPath=".status.readyPoolNum",description="The number of pools on which ingress is enabled"
// +kubebuilder:printcolumn:name="NotReadyPools",type="integer",JSONPath=".status.unreadyPoolNum",description="The number of pools on which ingress is enabling or enable failed"
// +kubebuilder:printcolumn:name="Replicas-Per-Pool",type="integer",JSONPath=".status.Replicas",description="The ingress controller replicas per pool"
// +kubebuilder:printcolumn:name="IngressController-Version",type="string",JSONPath=".status.Verson",description="The ingress controller version"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// NodePoolIngress is the Schema for the nodepoolingresses API
type NodePoolIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodePoolIngressSpec   `json:"spec,omitempty"`
	Status NodePoolIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodePoolIngressList contains a list of NodePoolIngress
type NodePoolIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodePoolIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodePoolIngress{}, &NodePoolIngressList{})
}
