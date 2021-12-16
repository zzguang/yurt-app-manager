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

package yurtingress

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1alpha1 "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/apis/apps/v1alpha1"
	"github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/util/gate"
	yurtapputil "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/util/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	controllerName = "yurtingress-controller"
)

// YurtIngressReconciler reconciles a YurtIngress object
type YurtIngressReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Add creates a new YurtIngress Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and start it when the Manager is started.
func Add(mgr manager.Manager, ctx context.Context) error {
	if !gate.ResourceEnabled(&appsv1alpha1.YurtIngress{}) {
		return nil
	}
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
//func newReconciler(mgr manager.Manager, createSingletonPoolIngress bool) reconcile.Reconciler {
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &YurtIngressReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		recorder: mgr.GetEventRecorderFor(controllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to YurtIngress
	err = c.Watch(&source.Kind{Type: &appsv1alpha1.YurtIngress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appsv1alpha1.YurtIngress{},
	})
	if err != nil {
		return err
	}
	return nil
}

//+kubebuilder:rbac:groups=apps.openyurt.io,resources=yurtingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.openyurt.io,resources=yurtingresses/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *YurtIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	if req.Name != appsv1alpha1.SingletonYurtIngressInstanceName {
		return ctrl.Result{}, nil
	}
	// Fetch the YurtIngress instance
	instance := &appsv1alpha1.YurtIngress{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Add finalizer if not exist
	if !controllerutil.ContainsFinalizer(instance, appsv1alpha1.YurtIngressFinalizer) {
		klog.V(4).Infof("Now add finalizers %s", appsv1alpha1.YurtIngressFinalizer)
		controllerutil.AddFinalizer(instance, appsv1alpha1.YurtIngressFinalizer)
		if err := r.Update(context.TODO(), instance); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Handle ingress controller resources cleanup
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.cleanupIngressResources(instance)
	}
	desiredNodePoolList := instance.Spec.Pools
	currentNodePoolList := instance.DeepCopy().Status.Pools
	statusNeedUpdate := false
	addedPools, removedPools, unchangedPools := getPools(desiredNodePoolList, currentNodePoolList)
	if addedPools != nil {
		klog.V(4).Infof("added pool list is %s", addedPools)
		statusNeedUpdate = true
		if currentNodePoolList == nil {
			yurtapputil.CreateNginxIngressCommonResource(r.Client)
		}
		for _, pool := range addedPools {
			yurtapputil.CreateNginxIngressSpecificResource(r.Client, pool, instance.Spec.Replicas)
		}
	}
	if removedPools != nil {
		klog.V(4).Infof("removed pool list is %s", removedPools)
		statusNeedUpdate = true
		for _, pool := range removedPools {
			yurtapputil.DeleteNginxIngressSpecificResource(r.Client, pool)
		}
		if desiredNodePoolList == nil {
			yurtapputil.DeleteNginxIngressCommonResource(r.Client)
		}
	}
	if unchangedPools != nil {
		klog.V(4).Infof("unchanged pool list is %s", unchangedPools)
		desiredControllerReplicas := instance.Spec.Replicas
		currentControllerReplicas := instance.Status.Replicas
		isControllerReplicasChanged := desiredControllerReplicas != currentControllerReplicas
		if isControllerReplicasChanged {
			statusNeedUpdate = true
			klog.V(4).Infof("Per-Pool ingress controller replicas is changed!")
			for _, pool := range unchangedPools {
				yurtapputil.ScaleNginxIngressControllerDeploymment(r.Client, pool, desiredControllerReplicas)
			}
		}
	}
	if statusNeedUpdate {
		r.updateStatus(instance)
	}
	return ctrl.Result{}, nil
}

func getPools(desired, current []string) (added, removed, unchanged []string) {
	swap := false
	for i := 0; i < 2; i++ {
		for _, s1 := range desired {
			found := false
			for _, s2 := range current {
				if s1 == s2 {
					found = true
					if !swap {
						unchanged = append(unchanged, s1)
					}
					break
				}
			}
			if !found {
				if !swap {
					added = append(added, s1)
				} else {
					removed = append(removed, s1)
				}
			}
		}
		if i == 0 {
			swap = true
			desired, current = current, desired
		}
	}
	return added, removed, unchanged
}

func (r *YurtIngressReconciler) updateStatus(ying *appsv1alpha1.YurtIngress) (*appsv1alpha1.YurtIngress, error) {
	ying.Status.Pools = ying.Spec.Pools
	ying.Status.Replicas = ying.Spec.Replicas
	ying.Status.Version = appsv1alpha1.NginxIngressControllerVersion

	var updateErr error
	for i, obj := 0, ying; ; i++ {
		updateErr = r.Client.Status().Update(context.TODO(), obj)
		if updateErr == nil {
			klog.V(4).Infof("%s status is updated!", obj.Name)
			return obj, nil
		}
		if i >= 5 {
			break
		}
	}
	klog.Errorf("fail to update YurtIngress %s status: %s", ying.Name, updateErr)
	return nil, updateErr
}

func (r *YurtIngressReconciler) cleanupIngressResources(instance *appsv1alpha1.YurtIngress) (ctrl.Result, error) {
	pools := instance.Spec.Pools
	if pools != nil {
		for _, pool := range pools {
			yurtapputil.DeleteNginxIngressSpecificResource(r.Client, pool)
		}
		yurtapputil.DeleteNginxIngressCommonResource(r.Client)
		r.updateStatus(instance)
	}
	if controllerutil.ContainsFinalizer(instance, appsv1alpha1.YurtIngressFinalizer) {
		controllerutil.RemoveFinalizer(instance, appsv1alpha1.YurtIngressFinalizer)
		if err := r.Update(context.TODO(), instance); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}
