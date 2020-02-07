// Package elbservice for implement controller
package elbservice

import (
	"context"

	k8sv1alpha1 "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
	"calmwu.org/elbservice-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_elbservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ELBService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileELBService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("elbservice-controller", mgr, controller.Options{
		Reconciler:              r,
		MaxConcurrentReconciles: 1, //启动一个worker
	})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ELBService
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.ELBService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// 监控二级资源 vipservice
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.ELBService{},
	})
	if err != nil {
		return err
	}

	// 监控二级资源 vipendpoints
	err = c.Watch(&source.Kind{Type: &corev1.Endpoints{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.ELBService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileELBService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileELBService{}

// ReconcileELBService reconciles a ELBService object
type ReconcileELBService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ELBService object and makes changes based on the state read
// and what is in the ELBService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileELBService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ELBService")

	// Fetch the ELBService instance
	instance := &k8sv1alpha1.ELBService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("ELBService resource not found, Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get ELBService.")
		return reconcile.Result{}, err
	}

	// 判断vipservice是否存在，不存在就创建
	vipService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetVIPServiceName(instance),
		Namespace: request.Namespace,
	}, vipService)
	if err != nil && errors.IsNotFound(err) {
		vipService = resources.NewVIPServiceForCR(instance)
		// 设置owner
		controllerutil.SetControllerReference(instance, vipService, r.scheme)
		reqLogger.Info("Creating a new VIPService.", "VIPService.Namespaec", vipService.Namespace, "VIPService.Name", vipService.Name)
		err = r.client.Create(context.TODO(), vipService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new VIPService", "VIPService.Namespaec", vipService.Namespace, "VIPService.Name", vipService.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VIPService.")
		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Skip reconcile: VIPService already exists", "VIPService.Namespaec", vipService.Namespace, "VIPService.Name", vipService.Name)
	}

	// 判断vipendpoints是否存在，不存在就创建
	vipEndpoints := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetVIPServiceName(instance),
		Namespace: request.Namespace,
	}, vipEndpoints)
	if err != nil && errors.IsNotFound(err) {
		vipEndpoints = resources.NewVIPEndpointForCR(instance)
		// 设置owner
		controllerutil.SetControllerReference(instance, vipEndpoints, r.scheme)
		reqLogger.Info("Creating a new VIPEndpoints.", "VIPEndpoints.Namespaec", vipEndpoints.Namespace, "VIPEndpoints.Name", vipEndpoints.Name)
		err = r.client.Create(context.TODO(), vipEndpoints)
		if err != nil {
			reqLogger.Error(err, "Failed to create new VIPEndpoints", "VIPEndpoints.Namespaec", vipEndpoints.Namespace, "VIPEndpoints.Name", vipEndpoints.Name)
			return reconcile.Result{}, err
		}

		// 创建一个pod用于测试域名解析
		dnsTestPod := resources.NewPodForCR(instance)
		controllerutil.SetControllerReference(instance, dnsTestPod, r.scheme)
		r.client.Create(context.TODO(), dnsTestPod)
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VIPEndpoints.")
		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Skip reconcile: VIPEndpoints already exists", "VIPEndpoints.Namespaec", vipEndpoints.Namespace, "VIPEndpoints.Name", vipEndpoints.Name)
	}

	return reconcile.Result{}, nil
}
