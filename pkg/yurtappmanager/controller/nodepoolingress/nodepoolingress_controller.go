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

package nodepoolingress

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1alpha1 "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/apis/apps/v1alpha1"
	"github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/util/gate"
	yurtapputil "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/util/kubernetes"
)

var concurrentReconciles = 3

const (
	controllerName = "nodepoolingress-controller"

	/*eventTypeIngressEnabling     = "IngressEnabling"
	eventTypeIngressEnabled      = "IngressEnabled"
	eventTypeIngressEnableFailed = "IngressEnableFailed"*/
)

// NodePoolIngressReconciler reconciles a NodePoolIngress object
type NodePoolIngressReconciler struct {
	client.Client
	Scheme                     *runtime.Scheme
	recorder                   record.EventRecorder
	createSingletonPoolIngress bool
}

// Add creates a new NodePoolIngress Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager, ctx context.Context) error {
	if !gate.ResourceEnabled(&appsv1alpha1.NodePoolIngress{}) {
		return nil
	}

	/*	inf := ctx.Value(constant.ContextKeyCreateSingletonPoolIngress)
		cdp, ok := inf.(bool)
		if !ok {
			return errors.New("fail to assert interface to bool for command line option createSingletonPoolIngress")
		}*/

	return add(mgr, newReconciler(mgr, true))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, createSingletonPoolIngress bool) reconcile.Reconciler {
	return &NodePoolIngressReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),

		recorder:                   mgr.GetEventRecorderFor(controllerName),
		createSingletonPoolIngress: createSingletonPoolIngress,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: concurrentReconciles})
	if err != nil {
		return err
	}

	npir, ok := r.(*NodePoolIngressReconciler)
	if !ok {
		return errors.New("fail to assert interface to NodePoolIngressReconciler")
	}
	// Watch for changes to NodePoolIngress
	err = c.Watch(&source.Kind{Type: &appsv1alpha1.NodePoolIngress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if npir.createSingletonPoolIngress {
		go createSingletonPoolIngress(mgr.GetClient())
	}
	return nil
}

//createSingletonPoolIngress creates the singleton NodePoolIngress CR, it will try 5 times if it fails
func createSingletonPoolIngress(client client.Client) {
	name := appsv1alpha1.SingletonNodePoolIngressInstanceName
	namespace := appsv1alpha1.SingletonNodePoolIngressNameSpace
	replicas := appsv1alpha1.DefaultIngressControllerReplicasPerPool
	for i := 0; i < 5; i++ {
		np_ing := appsv1alpha1.NodePoolIngress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: appsv1alpha1.NodePoolIngressSpec{
				Replicas: replicas,
			},
		}
		err := client.Create(context.TODO(), &np_ing)
		if err == nil {
			klog.V(4).Infof("the singleton nodepoolingress(%s) is created", name)
			return
		}
		if apierrors.IsAlreadyExists(err) {
			klog.V(4).Infof("the singleton nodepoolingress(%s) already exist", name)
			return
		}
		klog.Errorf("fail to create the singleton nodepoolingress(%s): %s, try again!", name, err)
		time.Sleep(2 * time.Second)
	}
	klog.V(4).Info("fail to create the singleton nodepoolingress after trying for 5 times")
}

//+kubebuilder:rbac:groups=apps.openyurt.io,resources=nodepoolingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.openyurt.io,resources=nodepoolingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.openyurt.io,resources=nodepoolingresses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NodePoolIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NodePoolIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	klog.V(4).Infof("Reconcile NodePoolIngress %s/%s", req.Namespace, req.Name)
	if req.Name != appsv1alpha1.SingletonNodePoolIngressInstanceName {
		return ctrl.Result{}, nil
	}

	// Fetch the NodePoolIngress instance
	instance := &appsv1alpha1.NodePoolIngress{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	desiredNodePoolList := instance.Spec.Pools
	currentNodePoolList := instance.DeepCopy().Status.Pools

	statusNeedUpdate := false

	addedPools, removedPools, unchangedPools := getPools(desiredNodePoolList, currentNodePoolList)

	if addedPools != nil {
		klog.V(4).Infof("added pool list is %s", addedPools)
		statusNeedUpdate = true
		if currentNodePoolList == nil {
			yurtapputil.CreateNginxIngressClusterRole(r.Client)
		}
		for _, pool := range addedPools {
			DeployNginxIngressController(r.Client, pool, instance.Spec.Replicas)
		}
	}
	if removedPools != nil {
		klog.V(4).Infof("removed pool list is %s", removedPools)
		statusNeedUpdate = true
		for _, pool := range removedPools {
			DeleteNginxIngressController(r.Client, pool)
		}
		if desiredNodePoolList == nil {
			yurtapputil.DeleteNginxIngressClusterRole(r.Client)
		}
	}
	if unchangedPools != nil {
		klog.V(4).Infof("unchanged pool list is %s", unchangedPools)
		desiredControllerReplicas := instance.Spec.Replicas
		currentControllerReplicas := instance.Status.Replicas
		isControllerReplicasChanged := desiredControllerReplicas != currentControllerReplicas
		if isControllerReplicasChanged {
			statusNeedUpdate = true
			klog.V(4).Infof("Per-pool ingress controller replicas is changed!")
			for _, pool := range unchangedPools {
				ns := GetPoolNameSpace(pool)
				yurtapputil.ScaleNginxIngressControllerDeploymment(r.Client, ns, pool, desiredControllerReplicas)
			}
		}
	}
	r.updateStatus(instance, statusNeedUpdate)

	return ctrl.Result{}, err
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

func (r *NodePoolIngressReconciler) updateStatus(nping *appsv1alpha1.NodePoolIngress,
	update bool) (*appsv1alpha1.NodePoolIngress, error) {
	if !update {
		return nping, nil
	}

	nping.Status.Pools = nping.Spec.Pools
	nping.Status.Replicas = nping.Spec.Replicas

	var updateErr error
	for i, obj := 0, nping; ; i++ {
		updateErr = r.Client.Status().Update(context.TODO(), obj)
		if updateErr == nil {
			klog.V(4).Infof("%s status is updated!", obj.Name)
			return obj, nil
		}
		if i >= 5 {
			break
		}
	}

	klog.Errorf("fail to update NodePoolIngress %s/%s status: %s", nping.Namespace, nping.Name, updateErr)
	return nil, updateErr
}

func DeployNginxIngressController(client client.Client, poolname string, desiredreplicas int32) {
	ns, err := CreateNodePoolNamespace(client, poolname)
	if err != nil {
		klog.V(4).Infof("namespace for %s is created failed: %v", poolname, err)
		return
	}
	yurtapputil.CreateNginxIngressControllerStaticResource(client, ns)
	yurtapputil.CreateNginxIngressControllerDeploymment(client, ns, poolname, desiredreplicas)
	yurtapputil.CreateNginxIngressWebhookStaticResource(client, ns)
	yurtapputil.CreateNginxIngressControllerWebhookDeploymment(client, ns, poolname, 1)
}

func DeleteNginxIngressController(client client.Client, poolname string) {
	ns := GetPoolNameSpace(poolname)
	yurtapputil.DeleteNginxIngressControllerStaticResource(client, ns)
	yurtapputil.DeleteNginxIngressControllerDeploymment(client, ns, poolname)
	yurtapputil.DeleteNginxIngressWebhookStaticResource(client, ns)
	yurtapputil.DeleteNginxIngressControllerWebhookDeploymment(client, ns, poolname)
	DeleteNodePoolNamespace(client, poolname)
}

func GetPoolNameSpace(poolname string) string {
	return strings.Join([]string{appsv1alpha1.NodePoolNameSpacePrefix, poolname}, "-")
}

func CreateNodePoolNamespace(client client.Client, poolname string) (string, error) {
	name := GetPoolNameSpace(poolname)
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := client.Create(context.TODO(), &namespace); err != nil {
		return "", fmt.Errorf("fail to create the namespace /%s: %v", namespace.Name, err)
	}
	klog.V(4).Infof("nodepool namespace %s is created", name)
	return name, nil
}

func DeleteNodePoolNamespace(client client.Client, poolname string) error {
	name := GetPoolNameSpace(poolname)
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := client.Delete(context.TODO(), &namespace)
	if err != nil {
		return fmt.Errorf("fail to delete the namespace /%s: %v", namespace.Name, err)
	}
	klog.V(4).Infof("namespace/%s is deleted", namespace.Name)

	return nil
}
