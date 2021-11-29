/*
Copyright 2020 The OpenYurt Authors.

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
	"bytes"
	"context"
	"fmt"
	"text/template"

	"k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateClusterRoleFromYaml creates the ClusterRole from the yaml template.
func CreateClusterRoleFromYaml(client client.Client, crTmpl string) error {
	obj, err := YamlToObject([]byte(crTmpl))
	if err != nil {
		return err
	}
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return fmt.Errorf("fail to assert clusterrole")
	}
	err = client.Create(context.Background(), cr)
	if err != nil {
		return fmt.Errorf("fail to create the clusterrole/%s: %v", cr.Name, err)
	}
	klog.V(4).Infof("clusterrole/%s is created", cr.Name)
	return nil
}

// DeleteClusterRoleFromYaml deletes the ClusterRole from the yaml template.
func DeleteClusterRoleFromYaml(client client.Client, crTmpl string) error {
	obj, err := YamlToObject([]byte(crTmpl))
	if err != nil {
		return err
	}
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return fmt.Errorf("fail to assert clusterrole")
	}
	err = client.Delete(context.Background(), cr)
	if err != nil {
		return fmt.Errorf("fail to delete the clusterrole/%s: %v", cr.Name, err)
	}
	klog.V(4).Infof("clusterrole/%s is deleted", cr.Name)
	return nil
}

// CreateServiceAccountFromYaml creates the ServiceAccount from the yaml template.
func CreateServiceAccountFromYaml(client client.Client, ns, saTmpl string) error {
	obj, err := YamlToObject([]byte(saTmpl))
	if err != nil {
		return err
	}
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return fmt.Errorf("fail to assert serviceaccount")
	}
	sa.Namespace = ns
	err = client.Create(context.Background(), sa)
	if err != nil {
		return fmt.Errorf("fail to create the serviceaccount/%s: %v", sa.Name, err)
	}
	klog.V(4).Infof("serviceaccount/%s is created", sa.Name)
	return nil
}

// DeleteServiceAccountFromYaml deletes the ServiceAccount from the yaml template.
func DeleteServiceAccountFromYaml(client client.Client, ns, saTmpl string) error {
	obj, err := YamlToObject([]byte(saTmpl))
	if err != nil {
		return err
	}
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return fmt.Errorf("fail to assert serviceaccount")
	}
	sa.Namespace = ns
	err = client.Delete(context.Background(), sa)
	if err != nil {
		return fmt.Errorf("fail to delete the serviceaccount/%s: %v", sa.Name, err)
	}
	klog.V(4).Infof("serviceaccount/%s is deleted", sa.Name)
	return nil
}

// CreateClusterRoleBindingFromYaml creates the ClusterRoleBinding from the yaml template.
func CreateClusterRoleBindingFromYaml(client client.Client, name, subject_ns, crbTmpl string) error {
	obj, err := YamlToObject([]byte(crbTmpl))
	if err != nil {
		return err
	}
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert clusterrolebinding")
	}
	crb.Name = name
	crb.Subjects[0].Namespace = subject_ns
	err = client.Create(context.Background(), crb)
	if err != nil {
		return fmt.Errorf("fail to create the clusterrolebinding/%s: %v", crb.Name, err)
	}
	klog.V(4).Infof("clusterrolebinding/%s is created", crb.Name)
	return nil
}

// DeleteClusterRoleBindingFromYaml deletes the ClusterRoleBinding from the yaml template.
func DeleteClusterRoleBindingFromYaml(client client.Client, name, subject_ns, crbTmpl string) error {
	obj, err := YamlToObject([]byte(crbTmpl))
	if err != nil {
		return err
	}
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert clusterrolebinding")
	}
	crb.Name = name
	crb.Subjects[0].Namespace = subject_ns
	err = client.Delete(context.Background(), crb)
	if err != nil {
		return fmt.Errorf("fail to delete the clusterrolebinding/%s: %v", crb.Name, err)
	}
	klog.V(4).Infof("clusterrolebinding/%s is deleted", crb.Name)
	return nil
}

// CreateConfigMapFromYaml creates the ConfigMap from the yaml template.
func CreateConfigMapFromYaml(client client.Client, ns, cmTmpl string) error {
	obj, err := YamlToObject([]byte(cmTmpl))
	if err != nil {
		return err
	}
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return fmt.Errorf("fail to assert configmap")
	}
	cm.Namespace = ns
	err = client.Create(context.Background(), cm)
	if err != nil {
		return fmt.Errorf("fail to create the configmap/%s: %v", cm.Name, err)
	}
	klog.V(4).Infof("configmap/%s is created", cm.Name)
	return nil
}

// DeleteConfigMapFromYaml deletes the ConfigMap from the yaml template.
func DeleteConfigMapFromYaml(client client.Client, ns, cmTmpl string) error {
	obj, err := YamlToObject([]byte(cmTmpl))
	if err != nil {
		return err
	}
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return fmt.Errorf("fail to assert configmap")
	}
	cm.Namespace = ns
	err = client.Delete(context.Background(), cm)
	if err != nil {
		return fmt.Errorf("fail to delete the configmap/%s: %v", cm.Name, err)
	}
	klog.V(4).Infof("configmap/%s is deleted", cm.Name)
	return nil
}

// CreateDeployFromYaml creates the Deployment from the yaml template.
func CreateDeployFromYaml(client client.Client, ns, dplyTmpl string, replicas *int32, ctx interface{}) error {
	dp, err := SubsituteTemplate(dplyTmpl, ctx)
	if err != nil {
		return err
	}
	dpObj, err := YamlToObject([]byte(dp))
	if err != nil {
		return err
	}
	dply, ok := dpObj.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("fail to assert deployment")
	}
	dply.Namespace = ns
	dply.Spec.Replicas = replicas
	err = client.Create(context.Background(), dply)
	if err != nil {
		return fmt.Errorf("fail to create the deployment/%s: %v", dply.Name, err)
	}
	klog.V(4).Infof("deployment/%s is created", dply.Name)
	return nil
}

// DeleteDeployFromYaml delete the Deployment from the yaml template.
func DeleteDeployFromYaml(client client.Client, ns, dplyTmpl string, ctx interface{}) error {
	dp, err := SubsituteTemplate(dplyTmpl, ctx)
	if err != nil {
		return err
	}
	dpObj, err := YamlToObject([]byte(dp))
	if err != nil {
		return err
	}
	dply, ok := dpObj.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("fail to assert deployment")
	}
	dply.Namespace = ns
	err = client.Delete(context.Background(), dply)
	if err != nil {
		return fmt.Errorf("fail to delete the deployment/%s: %v", dply.Name, err)
	}
	klog.V(4).Infof("deployment/%s is deleted", dply.Name)
	return nil
}

// UpdateDeployFromYaml updates the Deployment from the yaml template.
func UpdateDeployFromYaml(client client.Client, ns, dplyTmpl string, replicas *int32, ctx interface{}) error {
	dp, err := SubsituteTemplate(dplyTmpl, ctx)
	if err != nil {
		return err
	}
	dpObj, err := YamlToObject([]byte(dp))
	if err != nil {
		return err
	}
	dply, ok := dpObj.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("fail to assert deployment")
	}
	dply.Namespace = ns
	dply.Spec.Replicas = replicas
	err = client.Update(context.Background(), dply)
	if err != nil {
		return fmt.Errorf("fail to update the deployment/%s: %v", dply.Name, err)
	}
	klog.V(4).Infof("deployment/%s is updated", dply.Name)
	return nil
}

// CreateServiceFromYaml creates the Service from the yaml template.
func CreateServiceFromYaml(client client.Client, ns, svcTmpl string) error {
	obj, err := YamlToObject([]byte(svcTmpl))
	if err != nil {
		return err
	}
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return fmt.Errorf("fail to assert service")
	}
	svc.Namespace = ns
	err = client.Create(context.Background(), svc)
	if err != nil {
		return fmt.Errorf("fail to create the service/%s: %v", svc.Name, err)
	}
	klog.V(4).Infof("service/%s is created", svc.Name)
	return nil
}

// DeleteServiceFromYaml deletes the Service from the yaml template.
func DeleteServiceFromYaml(client client.Client, ns, svcTmpl string) error {
	obj, err := YamlToObject([]byte(svcTmpl))
	if err != nil {
		return err
	}
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return fmt.Errorf("fail to assert service")
	}
	svc.Namespace = ns
	err = client.Delete(context.Background(), svc)
	if err != nil {
		return fmt.Errorf("fail to delete the service/%s: %v", svc.Name, err)
	}
	klog.V(4).Infof("service/%s is deleted", svc.Name)
	return nil
}

// CreateRoleFromYaml creates the Role from the yaml template.
func CreateRoleFromYaml(client client.Client, ns, rTmpl string) error {
	obj, err := YamlToObject([]byte(rTmpl))
	if err != nil {
		return err
	}
	r, ok := obj.(*rbacv1.Role)
	if !ok {
		return fmt.Errorf("fail to assert role")
	}
	r.Namespace = ns
	err = client.Create(context.Background(), r)
	if err != nil {
		return fmt.Errorf("fail to create the role/%s: %v", r.Name, err)
	}
	klog.V(4).Infof("role/%s is created", r.Name)
	return nil
}

// DeleteRoleFromYaml deletes the Role from the yaml template.
func DeleteRoleFromYaml(client client.Client, ns, rTmpl string) error {
	obj, err := YamlToObject([]byte(rTmpl))
	if err != nil {
		return err
	}
	r, ok := obj.(*rbacv1.Role)
	if !ok {
		return fmt.Errorf("fail to assert role")
	}
	r.Namespace = ns
	err = client.Delete(context.Background(), r)
	if err != nil {
		return fmt.Errorf("fail to delete the role/%s: %v", r.Name, err)
	}
	klog.V(4).Infof("role/%s is deleted", r.Name)
	return nil
}

// CreateRoleBindingFromYaml creates the RoleBinding from the yaml template.
func CreateRoleBindingFromYaml(client client.Client, ns, rbTmpl string) error {
	obj, err := YamlToObject([]byte(rbTmpl))
	if err != nil {
		return err
	}
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert rolebinding")
	}
	rb.Namespace = ns
	rb.Subjects[0].Namespace = ns
	err = client.Create(context.Background(), rb)
	if err != nil {
		return fmt.Errorf("fail to create the rolebinding/%s: %v", rb.Name, err)
	}
	klog.V(4).Infof("rolebinding/%s is created", rb.Name)
	return nil
}

// DeleteRoleBindingFromYaml delete the RoleBinding from the yaml template.
func DeleteRoleBindingFromYaml(client client.Client, ns, rbTmpl string) error {
	obj, err := YamlToObject([]byte(rbTmpl))
	if err != nil {
		return err
	}
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert rolebinding")
	}
	rb.Namespace = ns
	rb.Subjects[0].Namespace = ns
	err = client.Delete(context.Background(), rb)
	if err != nil {
		return fmt.Errorf("fail to delete the rolebinding/%s: %v", rb.Name, err)
	}
	klog.V(4).Infof("rolebinding/%s is deleted", rb.Name)
	return nil
}

// CreateValidatingWebhookConfigurationFromYaml creates the validatingwebhookconfiguration from the yaml template.
func CreateValidatingWebhookConfigurationFromYaml(client client.Client, name, svc_ns, vwcTmpl string) error {
	obj, err := YamlToObject([]byte(vwcTmpl))
	if err != nil {
		return err
	}
	vwc, ok := obj.(*v1beta1.ValidatingWebhookConfiguration)
	if !ok {
		return fmt.Errorf("fail to assert validatingwebhookconfiguration")
	}
	vwc.Name = name
	vwc.Webhooks[0].ClientConfig.Service.Namespace = svc_ns
	err = client.Create(context.Background(), vwc)
	if err != nil {
		return fmt.Errorf("fail to create the validatingwebhookconfiguration/%s: %v", vwc.Name, err)
	}
	klog.V(4).Infof("validatingwebhookconfiguration/%s is created", vwc.Name)
	return nil
}

// DeleteValidatingWebhookConfigurationFromYaml delete the validatingwebhookconfiguration from the yaml template.
func DeleteValidatingWebhookConfigurationFromYaml(client client.Client, name, svc_ns, vwcTmpl string) error {
	obj, err := YamlToObject([]byte(vwcTmpl))
	if err != nil {
		return err
	}
	vwc, ok := obj.(*v1beta1.ValidatingWebhookConfiguration)
	if !ok {
		return fmt.Errorf("fail to assert validatingwebhookconfiguration")
	}
	vwc.Name = name
	vwc.Webhooks[0].ClientConfig.Service.Namespace = svc_ns
	err = client.Delete(context.Background(), vwc)
	if err != nil {
		return fmt.Errorf("fail to delete the validatingwebhookconfiguration/%s: %s", vwc.Name, err)
	}
	klog.V(4).Infof("validatingwebhookconfiguration/%s is deleted", vwc.Name)
	return nil
}

// CreateJobFromYaml creates the Job from the yaml template.
func CreateJobFromYaml(client client.Client, ns, jobTmpl string) error {
	obj, err := YamlToObject([]byte(jobTmpl))
	if err != nil {
		return err
	}
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return fmt.Errorf("fail to assert job")
	}
	job.Namespace = ns
	err = client.Create(context.Background(), job)
	if err != nil {
		return fmt.Errorf("fail to create the job/%s: %v", job.Name, err)
	}
	klog.V(4).Infof("job/%s is created", job.Name)
	return nil
}

// DeleteJobFromYaml deletes the Job from the yaml template.
func DeleteJobFromYaml(client client.Client, ns, jobTmpl string) error {
	obj, err := YamlToObject([]byte(jobTmpl))
	if err != nil {
		return err
	}
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return fmt.Errorf("fail to assert job")
	}
	job.Namespace = ns
	err = client.Delete(context.Background(), job)
	if err != nil {
		return fmt.Errorf("fail to delete the job/%s: %v", job.Name, err)
	}
	klog.V(4).Infof("job/%s is deleted", job.Name)
	return nil
}

// CreateJobPatchFromYaml creates the Job patch from the yaml template.
func CreateJobPatchFromYaml(client client.Client, ns, jobTmpl string, ctx interface{}) error {
	jp, err := SubsituteTemplate(jobTmpl, ctx)
	if err != nil {
		return err
	}
	jpObj, err := YamlToObject([]byte(jp))
	if err != nil {
		return err
	}
	job, ok := jpObj.(*batchv1.Job)
	if !ok {
		return fmt.Errorf("fail to assert job")
	}
	job.Namespace = ns
	err = client.Create(context.Background(), job)
	if err != nil {
		return fmt.Errorf("fail to create the job patch/%s: %v", job.Name, err)
	}
	klog.V(4).Infof("job patch/%s is created", job.Name)
	return nil
}

// YamlToObject deserializes object in yaml format to a runtime.Object
func YamlToObject(yamlContent []byte) (k8sruntime.Object, error) {
	decode := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer().Decode
	obj, _, err := decode(yamlContent, nil, nil)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// SubsituteTemplate fills out the template based on the context
func SubsituteTemplate(tmpl string, context interface{}) (string, error) {
	t, tmplPrsErr := template.New("test").Option("missingkey=zero").Parse(tmpl)
	if tmplPrsErr != nil {
		return "", tmplPrsErr
	}
	writer := bytes.NewBuffer([]byte{})
	if err := t.Execute(writer, context); nil != err {
		return "", err
	}
	return writer.String(), nil
}
