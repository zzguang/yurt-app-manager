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

package validating

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/apis/apps/v1alpha1"
)

// validateNodePoolIngressSpec validates the nodepool ingress spec.
func validateNodePoolIngressSpec(c client.Client, spec *appsv1alpha1.NodePoolIngressSpec) field.ErrorList {

	if len(spec.Pools) > 0 {
		var err error
		var errmsg string
		nps := appsv1alpha1.NodePoolList{}
		if err = c.List(context.TODO(), &nps, &client.ListOptions{}); err != nil {
			errmsg = "List nodepool list error!"
			klog.Errorf(errmsg)
			return field.ErrorList([]*field.Error{
				field.Forbidden(field.NewPath("spec").Child("pools"), errmsg)})
		}

		// validate whether the nodepool exist
		var found = false
		if len(nps.Items) > 0 {
			for _, snp := range spec.Pools { //go through the nodepools setting in yaml
				for _, cnp := range nps.Items { //go through the nodepools in cluster
					if snp == cnp.ObjectMeta.Name {
						found = true
						break
					}
				}
				if !found {
					errmsg = snp + " does not exist in the cluster!"
					klog.Errorf(errmsg)
					return field.ErrorList([]*field.Error{
						field.Forbidden(field.NewPath("spec").Child("pools"), errmsg)})
				}
				found = false
			}

		}
	}
	return nil
}

func validateNodePoolIngressSpecUpdate(c client.Client, spec, oldSpec *appsv1alpha1.NodePoolIngressSpec) field.ErrorList {
	return validateNodePoolIngressSpec(c, spec)
}

func validateNodePoolIngressDeletion(c client.Client, spec *appsv1alpha1.NodePoolIngressSpec) field.ErrorList {
	return validateNodePoolIngressSpec(c, spec)
}
