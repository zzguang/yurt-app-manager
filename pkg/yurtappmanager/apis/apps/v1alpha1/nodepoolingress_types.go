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
	// SingletonNodePoolIngressInstanceName defines the singleton instance name of NodePoolIngress
	SingletonNodePoolIngressInstanceName = "nodepool-ingress"
	// SingletonNodePoolIngressNameSpace defines the namespace of the singleton NodePoolIngress instance
	SingletonNodePoolIngressNameSpace = "kube-system"
	// NodePoolNameSpacePrefix defines the prefix of the nodepool namespace for nginx ingress
	NodePoolNameSpacePrefix = "nodepool"
)

// NodePoolIngressConditionType indicates valid conditions type of a NodePoolIngress
type NodePoolIngressConditionType string

const (
	// IngressEnabling means ingress is being enabled on some nodepools.
	IngressEnabling NodePoolIngressConditionType = "IngressEnabling"
	// IngressEnabled means ingress is enabled successfully on some nodepools.
	IngressEnabled NodePoolIngressConditionType = "IngressEnabled"
	// IngressEnableFailed means ingress is enabled failed on some nodepools.
	IngressEnableFailed NodePoolIngressConditionType = "IngressEnableFailed"
)

// Topology defines the spread detail of each nodepool under NodePoolIngress
// A NodePoolIngress CRD manages ingress feature support for multi nodepools
// Each of the nodepools under the NodePoolIngress is described in Topology.
/*
type IngressTopology struct {
	// Indicates details of each nodepool which will be managed by NodePoolIngress.
	Pools []IngressPool `json:"pools,omitempty"`
}

// IngressPool defines the detail of a nodepool
type IngressPool struct {
	// Indicates the nodepool name.
	Name string `json:"name,omitempty"`

	// Indicates the number of the ingress controllers to be created under this nodepool.
	//Replicas *int32 `json:"ingress_controller_replicas,omitempty"`

	// Indicates the ingress controller external IP to be exposed.
	//ExternalIP string `json:"ingress_controller_externalIP,omitempty"`
}*/

// NodePoolStatus defines the observed state of the related nodepoolingress
type NodePoolIngressCondition struct {
	// Indicates the nodepoolingress state.
	Condition NodePoolIngressConditionType `json:"state,omitempty"`
	// Indicates the related nodepool names
	Pools []string `json:"pools,omitempty"`
}

// NodePoolIngressSpec defines the desired state of NodePoolIngress
type NodePoolIngressSpec struct {
	// Indicates the number of the ingress controllers to be deployed under all the specified nodepools.
	Replicas int32 `json:"ingress_controller_replicas_per_pool,omitempty"`

	// Indicates the nginx ingress controller version to be deployed under all the specified nodepools.
	//Version string `json:"nginx_ingress_controller_version,omitempty"`

	// Indicates all the nodepools on which to enable ingress.
	//Topology IngressTopology `json:"topology,omitempty"`
	Pools []string `json:"pools,omitempty"`
}

// NodePoolIngressStatus defines the observed state of NodePoolIngress
type NodePoolIngressStatus struct {
	// Indicates the number of the ingress controllers to be deployed under all the specified nodepools.
	Replicas int32 `json:"ingress_controller_replicas_per_pool,omitempty"`

	// Indicates the nginx ingress controller version to be deployed under all the specified nodepools.
	//Version string `json:"nginx_ingress_controller_version,omitempty"`

	// Indicates all the nodepools on which to enable ingress.
	Pools []string `json:"pools,omitempty"`

	// Indicates all the nodepools ingress status.
	Status []NodePoolIngressCondition `json:"pools_status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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
