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

package kubernetes

import (
	"strings"

	"github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/constant"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NginxIngressClusterRolePrefix          = "cluster-role"
	NginxIngressWebhookClusterRolePrefix   = "cluster-role-admission"
	NginxIngressWebhookConfigurationPrefix = "validatingwebhook"
)

func CreateNginxIngressClusterRole(client client.Client) error {
	if err := CreateClusterRoleFromYaml(client, constant.NginxIngressControllerClusterRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	if err := CreateClusterRoleFromYaml(client, constant.NginxIngressControllerWebhookClusterRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func DeleteNginxIngressClusterRole(client client.Client) error {
	if err := DeleteClusterRoleFromYaml(client, constant.NginxIngressControllerClusterRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	if err := DeleteClusterRoleFromYaml(client, constant.NginxIngressControllerWebhookClusterRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func CreateNginxIngressControllerStaticResource(client client.Client, ns string) error {

	// 1. Create the ServiceAccount
	if err := CreateServiceAccountFromYaml(client, ns,
		constant.NginxIngressControllerServiceAccount); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 2. Create the Configmap
	if err := CreateConfigMapFromYaml(client, ns,
		constant.NginxIngressControllerConfigMap); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 3. Create the ClusterRoleBinding
	name := strings.Join([]string{NginxIngressClusterRolePrefix, ns}, "-")
	if err := CreateClusterRoleBindingFromYaml(client, name, ns,
		constant.NginxIngressControllerClusterRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 4. Create the Role
	if err := CreateRoleFromYaml(client, ns,
		constant.NginxIngressControllerRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 5. Create the RoleBinding
	if err := CreateRoleBindingFromYaml(client, ns,
		constant.NginxIngressControllerRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 6. Create the Service
	if err := CreateServiceFromYaml(client, ns,
		constant.NginxIngressControllerService); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	return nil
}

func DeleteNginxIngressControllerStaticResource(client client.Client, ns string) error {
	// 1. Delete the ServiceAccount
	if err := DeleteServiceAccountFromYaml(client, ns,
		constant.NginxIngressControllerServiceAccount); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 2. Delete the Configmap
	if err := DeleteConfigMapFromYaml(client, ns,
		constant.NginxIngressControllerConfigMap); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 3. Delete the ClusterRoleBinding
	name := strings.Join([]string{NginxIngressClusterRolePrefix, ns}, "-")
	if err := DeleteClusterRoleBindingFromYaml(client, name, ns,
		constant.NginxIngressControllerClusterRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 4. Delete the Role
	if err := DeleteRoleFromYaml(client, ns,
		constant.NginxIngressControllerRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 5. Delete the RoleBinding
	if err := DeleteRoleBindingFromYaml(client, ns,
		constant.NginxIngressControllerRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 6. Delete the Service
	if err := DeleteServiceFromYaml(client, ns,
		constant.NginxIngressControllerService); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	return nil
}

func CreateNginxIngressControllerDeploymment(client client.Client, ns, poolname string, replicas int32) error {
	if err := CreateDeployFromYaml(client, ns,
		constant.NginxIngressControllerDeployment,
		&replicas,
		map[string]string{
			"nodepool_name": poolname}); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func DeleteNginxIngressControllerDeploymment(client client.Client, ns, poolname string) error {
	if err := DeleteDeployFromYaml(client, ns,
		constant.NginxIngressControllerDeployment,
		map[string]string{
			"nodepool_name": poolname}); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func ScaleNginxIngressControllerDeploymment(client client.Client, ns, poolname string, replicas int32) error {
	if err := UpdateDeployFromYaml(client, ns,
		constant.NginxIngressControllerDeployment,
		&replicas,
		map[string]string{
			"nodepool_name": poolname}); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func CreateNginxIngressControllerWebhookDeploymment(client client.Client, ns, poolname string, replicas int32) error {
	if err := CreateDeployFromYaml(client, ns,
		constant.NginxIngressControllerWebhookDeployment,
		&replicas,
		map[string]string{
			"nodepool_name": poolname}); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func DeleteNginxIngressControllerWebhookDeploymment(client client.Client, ns, poolname string) error {
	if err := DeleteDeployFromYaml(client, ns,
		constant.NginxIngressControllerWebhookDeployment,
		map[string]string{
			"nodepool_name": poolname}); err != nil {
		klog.Errorf("%v", err)
		return err
	}
	return nil
}

func CreateNginxIngressWebhookStaticResource(client client.Client, ns string) error {
	// 1. Create the ServiceAccount
	if err := CreateServiceAccountFromYaml(client, ns,
		constant.NginxIngressControllerWebhookServiceAccount); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 2. Create the ClusterRoleBinding
	name := strings.Join([]string{NginxIngressWebhookClusterRolePrefix, ns}, "-")
	if err := CreateClusterRoleBindingFromYaml(client, name, ns,
		constant.NginxIngressControllerWebhookClusterRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 3. Create the Role
	if err := CreateRoleFromYaml(client, ns,
		constant.NginxIngressControllerWebhookRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 4. Create the RoleBinding
	if err := CreateRoleBindingFromYaml(client, ns,
		constant.NginxIngressControllerWebhookRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 5. Create the Service
	if err := CreateServiceFromYaml(client, ns,
		constant.NginxIngressControllerWebhookService); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 6. Create the ValidatingWebhookConfiguration
	name = strings.Join([]string{NginxIngressWebhookConfigurationPrefix, ns}, "-")
	if err := CreateValidatingWebhookConfigurationFromYaml(client, name, ns,
		constant.NginxIngressControllerValidatingWebhook); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 7. Create the Job
	if err := CreateJobFromYaml(client, ns,
		constant.NginxIngressControllerWebhookJob); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 8. Create the Job Patch
	if err := CreateJobFromYaml(client, ns,
		constant.NginxIngressControllerWebhookJobPatch); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	return nil
}

func DeleteNginxIngressWebhookStaticResource(client client.Client, ns string) error {
	// 1. Delete the ServiceAccount
	if err := DeleteServiceAccountFromYaml(client, ns,
		constant.NginxIngressControllerWebhookServiceAccount); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 2. Delete the ClusterRoleBinding
	name := strings.Join([]string{NginxIngressWebhookClusterRolePrefix, ns}, "-")
	if err := DeleteClusterRoleBindingFromYaml(client, name, ns,
		constant.NginxIngressControllerWebhookClusterRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 3. Delete the Role
	if err := DeleteRoleFromYaml(client, ns,
		constant.NginxIngressControllerWebhookRole); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 4. Delete the RoleBinding
	if err := DeleteRoleBindingFromYaml(client, ns,
		constant.NginxIngressControllerWebhookRoleBinding); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 5. Delete the Service
	if err := DeleteServiceFromYaml(client, ns,
		constant.NginxIngressControllerWebhookService); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 6. Delete the ValidatingWebhookConfiguration
	name = strings.Join([]string{NginxIngressWebhookConfigurationPrefix, ns}, "-")
	if err := DeleteValidatingWebhookConfigurationFromYaml(client, name, ns,
		constant.NginxIngressControllerValidatingWebhook); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 7. Delete the Job
	if err := DeleteJobFromYaml(client, ns,
		constant.NginxIngressControllerWebhookJob); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	// 8. Delete the Job Patch
	if err := DeleteJobFromYaml(client, ns,
		constant.NginxIngressControllerWebhookJobPatch); err != nil {
		klog.Errorf("%v", err)
		return err
	}

	return nil
}
