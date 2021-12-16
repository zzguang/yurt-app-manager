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

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutil "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/webhook/util"
)

// YurtIngressCreateUpdateHandler handles YurtIngress
type YurtIngressCreateUpdateHandler struct {
	Client client.Client

	// Decoder decodes objects
	Decoder *admission.Decoder
}

var _ webhookutil.Handler = &YurtIngressCreateUpdateHandler{}

func (h *YurtIngressCreateUpdateHandler) SetOptions(options webhookutil.Options) {
	return
}

// Handle handles admission requests.
func (h *YurtIngressCreateUpdateHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	//np_ing := appsv1alpha1.YurtIngress{}

	switch req.AdmissionRequest.Operation {
	case admissionv1.Create:
		/*		klog.V(4).Info("capture the nodepool creation request")
				err := h.Decoder.Decode(req, &np)
				if err != nil {
					return admission.Errored(http.StatusBadRequest, err)
				}
				if allErrs := validateNodePoolSpec(&np.Spec); len(allErrs) > 0 {
					return admission.Errored(http.StatusUnprocessableEntity,
						allErrs.ToAggregate())
				}*/
	case admissionv1.Update:
		/*		klog.V(4).Info("capture the nodepool update request")
				err := h.Decoder.Decode(req, &np)
				if err != nil {
					return admission.Errored(http.StatusBadRequest, err)
				}
				onp := appsv1alpha1.NodePool{}
				err = h.Decoder.DecodeRaw(req.OldObject, &onp)
				if err != nil {
					return admission.Errored(http.StatusBadRequest, err)
				}

				if allErrs := validateNodePoolSpecUpdate(&np.Spec, &onp.Spec); len(allErrs) > 0 {
					return admission.Errored(http.StatusUnprocessableEntity,
						allErrs.ToAggregate())
				}*/
	case admissionv1.Delete:
		/*		klog.V(4).Info("capture the nodepool deletion request")
				err := h.Decoder.DecodeRaw(req.OldObject, &np)
				if err != nil {
					return admission.Errored(http.StatusBadRequest, err)
				}
				if allErrs := validateNodePoolDeletion(h.Client, &np); len(allErrs) > 0 {
					return admission.Errored(http.StatusUnprocessableEntity,
						allErrs.ToAggregate())
				}*/
	}

	return admission.ValidationResponse(true, "")
}

var _ admission.DecoderInjector = &YurtIngressCreateUpdateHandler{}

// InjectDecoder injects the decoder into the YurtIngressCreateUpdateHandler
func (h *YurtIngressCreateUpdateHandler) InjectDecoder(d *admission.Decoder) error {
	h.Decoder = d
	return nil
}

var _ inject.Client = &YurtIngressCreateUpdateHandler{}

// InjectClient injects the client into the PodCreateHandler
func (h *YurtIngressCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}
