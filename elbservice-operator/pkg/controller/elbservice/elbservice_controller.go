// Package elbservice for implement controller
package elbservice

import (
	"context"
	"time"

	k8sv1alpha1 "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
	"calmwu.org/elbservice-operator/pkg/resources"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
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
	// 加入mgr中
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
		MaxConcurrentReconciles: 2, //启动一个worker
	})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ELBService
	// 监控主资源
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.ELBService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// 监控依赖主资源的dependent资源，这些dependent的owner为ELBService
	// 监控二级资源 headlessservice，这些资源owner必须是ELBService
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.ELBService{},
	})
	if err != nil {
		return err
	}

	// 监控二级资源 endpoints
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
		reqLogger.Info("Creating a new VIPService.", "VIPService.Namespace", vipService.Namespace, "VIPService.Name", vipService.Name)
		err = r.client.Create(context.TODO(), vipService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new VIPService", "VIPService.Namespace", vipService.Namespace, "VIPService.Name", vipService.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VIPService.")
		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Skip reconcile: VIPService already exists", "VIPService.Namespace", vipService.Namespace, "VIPService.Name", vipService.Name)
	}

	// 判断vipendpoints是否存在，不存在就创建
	vipSvcEndpoints := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetVIPServiceName(instance),
		Namespace: request.Namespace,
	}, vipSvcEndpoints)
	if err != nil && errors.IsNotFound(err) {
		vipSvcEndpoints = resources.NewVIPEndpointForCR(instance)
		// 设置owner
		controllerutil.SetControllerReference(instance, vipSvcEndpoints, r.scheme)
		reqLogger.Info("Creating a new VIPEndpoints.", "VIPEndpoints.Namespace", vipSvcEndpoints.Namespace, "VIPEndpoints.Name", vipSvcEndpoints.Name)
		err = r.client.Create(context.TODO(), vipSvcEndpoints)
		if err != nil {
			reqLogger.Error(err, "Failed to create new VIPEndpoints", "VIPEndpoints.Namespace", vipSvcEndpoints.Namespace, "VIPEndpoints.Name", vipSvcEndpoints.Name)
			return reconcile.Result{}, err
		}

		// 创建一个pod用于测试域名解析
		// dnsTestPod := resources.NewPodForCR(instance)
		// controllerutil.SetControllerReference(instance, dnsTestPod, r.scheme)
		// r.client.Create(context.TODO(), dnsTestPod)

		// Endpoints created successfully - return and requeue
		return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VIPEndpoints.")
		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Skip reconcile: VIPEndpoints already exists", "VIPEndpoints.Namespaec", vipSvcEndpoints.Namespace, "VIPEndpoints.Name", vipSvcEndpoints.Name)
	}

	// 判断epservice是否存在，不存在就创建
	epService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetEPServiceName(instance),
		Namespace: request.Namespace,
	}, epService)
	if err != nil && errors.IsNotFound(err) {
		epService = resources.NewEPServiceForCR(instance)
		// 设置owner
		controllerutil.SetControllerReference(instance, epService, r.scheme)
		reqLogger.Info("Creating a new EPService.", "EPService.Namespace", epService.Namespace, "EPService.Name", epService.Name)
		err = r.client.Create(context.TODO(), epService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new EPService", "EPService.Namespace", epService.Namespace, "EPService.Name", epService.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get EPService.")
		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Skip reconcile: EPService already exists", "EPService.Namespace", epService.Namespace, "EPService.Name", epService.Name)
	}

	epSvcEndPoints := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetEPServiceName(instance),
		Namespace: request.Namespace,
	}, epSvcEndPoints)
	if err != nil && errors.IsNotFound(err) {
		// 找不到继续
		return reconcile.Result{Requeue: true, RequeueAfter: 3 * time.Second}, nil
	} else if err == nil {
		// 修改owner，这里只需要做一次
		reqLogger.Info("----Get EPSvcEndpoints----", "EPSvcEndpoints.Namespace", epSvcEndPoints.Namespace, "EPSvcEndpoints.Name", epSvcEndPoints.Name)
		err = controllerutil.SetControllerReference(instance, epSvcEndPoints, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Set EPSvcEndpoints OwnerReferences failed.")
		} else {
			// 更新
			err = r.client.Update(context.TODO(), epSvcEndPoints)
			if err != nil {
				reqLogger.Error(err, "Update EPSvcEndpoints failed.")
			}
		}

		// 获得subnet信息
		subnetPodAddrSet := func() *hashset.Set {
			addrSet := hashset.New()
			for i := range epSvcEndPoints.Subsets {
				epSubset := &epSvcEndPoints.Subsets[i]
				for j := range epSubset.Addresses {
					addrSet.Add(epSubset.Addresses[j].IP)
					reqLogger.Info("***addrSet***", "IP", epSubset.Addresses[j].IP)
				}
			}
			return addrSet
		}()

		// 根据label查询对应的pod，需要在pod变化是有回调
		elbPodList := &corev1.PodList{}
		listOpts := []client.ListOption{
			client.InNamespace(instance.GetNamespace()),
			client.MatchingLabels(instance.Spec.Selector),
		}
		err = r.client.List(context.TODO(), elbPodList, listOpts...)
		if err != nil {
			reqLogger.Error(err, "Failed to list pods.", "ELBService.Namespace", instance.Namespace, "ELBService.Name", instance.Name,
				"ELBService.Selector", instance.Spec.Selector)
			return reconcile.Result{}, err
		}
		elbServiceStatus := getELBServiceStatus(elbPodList.Items, subnetPodAddrSet, reqLogger)

		// 判断当前状态和计算出来是否相同
		if !cmp.Equal(*elbServiceStatus, instance.Status) {
			reqLogger.Info("Update ELBService status", "ELBService.Namespace", instance.Namespace, "ELBService.Name", instance.Name,
				"Calculation ELBService.Status", elbServiceStatus, "Current ELBService.instance.Status", instance.Status)
			// 更新状态
			/*
						Status:
				  			Podcount:  2
				  			Podinfos:
				    			Name:   elbservice-pod-2
				    			Podip:  10.244.62.163
				    			Name:   elbservice-pod-1
				    			Podip:  10.244.62.175
			*/
			instance.Status = *elbServiceStatus
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to Update ELBService status.")
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

func getELBServiceStatus(pods []corev1.Pod, podAddrSet *hashset.Set, reqLogger logr.Logger) *k8sv1alpha1.ELBServiceStatus {
	ipCount := podAddrSet.Size()
	podCount := len(pods)

	status := &k8sv1alpha1.ELBServiceStatus{
		PodCount: int32(ipCount),
		PodInfos: make([]k8sv1alpha1.ELBPodInfo, ipCount),
	}

	for index, pos := 0, 0; pos < ipCount && index < podCount; index++ {
		pod := &pods[index]
		reqLogger.Info("---Pod---", "Pod.Name", pod.Name, "Pod.Phase", pod.Status.Phase, "Pod.IP", pod.Status.PodIP)
		if pod.Status.HostIP == "" {
			continue
		}

		if podAddrSet.Contains(pod.Status.PodIP) {
			status.PodInfos[pos] = k8sv1alpha1.ELBPodInfo{
				Name:   pod.Name,
				PodIP:  pod.Status.PodIP,
				Status: pod.Status.Phase,
			}
			pos++
		}
	}

	return status
}
