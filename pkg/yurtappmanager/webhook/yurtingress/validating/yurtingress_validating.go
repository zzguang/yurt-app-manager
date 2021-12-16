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
	"k8s.io/apimachinery/pkg/util/validation/field"

	appsv1alpha1 "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/apis/apps/v1alpha1"
)

// validateYurtIngressSpec validates the yurt ingress spec.
func validateYurtIngressSpec(spec *appsv1alpha1.YurtIngressSpec) field.ErrorList {
	return nil
}

// validateYurtIngressSpecUpdate tests if required fields in the YurtIngress spec are set.
func validateYurtIngressSpecUpdate(spec, oldSpec *appsv1alpha1.YurtIngressSpec) field.ErrorList {
	if allErrs := validateYurtIngressSpec(spec); allErrs != nil {
		return allErrs
	}

	return nil
}