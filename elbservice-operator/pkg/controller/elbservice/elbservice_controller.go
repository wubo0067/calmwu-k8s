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
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// https://github.com/kubernetes/kubernetes/issues/84430

var log = logf.Log.WithName("controller_elbservice")

const (
	finalizerName = "finalizer.elbservice.calmwu"
)

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
	elbServiceInst := &k8sv1alpha1.ELBService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, elbServiceInst)
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

	reqLogger.Info("<<<ELBService>>>", "name:", elbServiceInst.Name, "resourceVersion", elbServiceInst.ResourceVersion)

	switch elbServiceInst.Status.Phase {
	case k8sv1alpha1.ELBServiceNone:
		// 更新状态
		reqLogger.Info("-----k8sv1alpha1.ELBServiceNone-----")
		if err := r.makeCreatingPhase(elbServiceInst); err != nil {
			// 重试
			return reconcile.Result{}, err
		}
	case k8sv1alpha1.ELBServiceCreating, k8sv1alpha1.ELBServiceFailed:
		// 开始创建
		reqLogger.Info("-----k8sv1alpha1.ELBServiceCreating/ELBServiceFailed-----")
		if err := r.createService(elbServiceInst); err != nil {
			return reconcile.Result{}, err
		}
	case k8sv1alpha1.ELBServiceActive:
		reqLogger.Info("-----k8sv1alpha1.ELBServiceActive-----")
		// 判断是否被删除
		if !elbServiceInst.ObjectMeta.DeletionTimestamp.IsZero() {
			// 更新状态为ELBServiceTerminating
			if err := r.makeTerminatingPhase(elbServiceInst); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}

		if err := r.updateUsage(elbServiceInst); err != nil {
			return reconcile.Result{}, err
		}
	case k8sv1alpha1.ELBServiceTerminating:
		reqLogger.Info("-----k8sv1alpha1.ELBServiceTerminating-----")
		// 解绑pod和elb之间的关系，删除finalizer标记
		if err := r.deleteService(elbServiceInst); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func getELBServiceStatus(pods []corev1.Pod, podAddrSet *hashset.Set, logger logr.Logger) *k8sv1alpha1.ELBServiceStatus {
	ipCount := podAddrSet.Size()
	podCount := len(pods)

	status := &k8sv1alpha1.ELBServiceStatus{
		PodCount: int32(ipCount),
		PodInfos: make([]k8sv1alpha1.ELBPodInfo, ipCount),
	}

	for index, pos := 0, 0; pos < ipCount && index < podCount; index++ {
		pod := &pods[index]
		logger.Info("---Pod---", "Pod.Name", pod.Name, "Pod.Phase", pod.Status.Phase, "Pod.IP", pod.Status.PodIP)
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

func (r *ReconcileELBService) updateUsage(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("update Usage")

	epSvcEndPoints := &corev1.Endpoints{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetEPServiceName(elbServiceInst),
		Namespace: elbServiceInst.Namespace,
	}, epSvcEndPoints)
	if err != nil {

	} else {
		// 获得subnet信息
		subnetPodAddrSet := func() *hashset.Set {
			addrSet := hashset.New()
			for i := range epSvcEndPoints.Subsets {
				epSubset := &epSvcEndPoints.Subsets[i]
				for j := range epSubset.Addresses {
					addrSet.Add(epSubset.Addresses[j].IP)
					logger.Info("***addrSet***", "IP", epSubset.Addresses[j].IP)
				}
			}
			return addrSet
		}()

		// 根据label查询对应的pod，需要在pod变化是有回调
		elbPodList := &corev1.PodList{}
		listOpts := []client.ListOption{
			client.InNamespace(elbServiceInst.GetNamespace()),
			client.MatchingLabels(elbServiceInst.Spec.Selector),
		}
		err = r.client.List(context.TODO(), elbPodList, listOpts...)
		if err != nil {
			logger.Error(err, "Failed to list pods.", "ELBService.Namespace", elbServiceInst.Namespace, "ELBService.Name", elbServiceInst.Name,
				"ELBService.Selector", elbServiceInst.Spec.Selector)
			return err
		}
		elbServiceStatus := getELBServiceStatus(elbPodList.Items, subnetPodAddrSet, logger)

		// 判断当前状态和计算出来是否相同
		if !cmp.Equal(elbServiceStatus.PodInfos, elbServiceInst.Status.PodInfos) {
			logger.Info("Update ELBService.Status.PodInfos", "Calculation ELBService.PodInfos", elbServiceStatus.PodInfos,
				"Current ELBService.PodInfos", elbServiceInst.Status.PodInfos)

			// 这里要和elb进行绑定
			/*
						Status:
				  			Podcount:  2
				  			Podinfos:
				    			Name:   elbservice-pod-2
				    			Podip:  10.244.62.163
				    			Name:   elbservice-pod-1
				    			Podip:  10.244.62.175
			*/
			elbServiceInst.Status.PodCount = elbServiceStatus.PodCount
			elbServiceInst.Status.PodInfos = elbServiceStatus.PodInfos
			if err := r.client.Status().Update(context.TODO(), elbServiceInst); err != nil {
				if !k8serrors.IsConflict(err) {
					logger.Error(err, "Update ELBService.Status.PodInfos failed.")
					return err
				}

			} else {
				logger.Info("Update ELBService.Status.PodInfos successed.")
			}
		}
	}

	return nil
}

func (r *ReconcileELBService) makeCreatingPhase(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("update", "OldPhase", elbServiceInst.Status.Phase, " NewPhase", "ELBServiceCreating")
	elbServiceInst.Status.PodInfos = make([]k8sv1alpha1.ELBPodInfo, 0)
	return r.updateServiceStatus(elbServiceInst, k8sv1alpha1.ELBServiceCreating, nil)
}

func (r *ReconcileELBService) makeTerminatingPhase(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("update", "OldPhase", elbServiceInst.Status.Phase, " NewPhase", "ELBServiceTerminating")
	return r.updateServiceStatus(elbServiceInst, k8sv1alpha1.ELBServiceTerminating, nil)
}

func (r *ReconcileELBService) makeFailedPhase(elbServiceInst *k8sv1alpha1.ELBService, err error) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("update", "OldPhase", elbServiceInst.Status.Phase, " NewPhase", "ELBServiceFailed")
	return r.updateServiceStatus(elbServiceInst, k8sv1alpha1.ELBServiceFailed, err)
}

func (r *ReconcileELBService) makeActivePhase(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("update", "OldPhase", elbServiceInst.Status.Phase, " NewPhase", "ELBServiceActive")
	return r.updateServiceStatus(elbServiceInst, k8sv1alpha1.ELBServiceActive, nil)
}

func (r *ReconcileELBService) createService(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("createService")

	// 判断vipservice是否存在，不存在就创建
	vipService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetVIPServiceName(elbServiceInst),
		Namespace: elbServiceInst.Namespace,
	}, vipService)
	if err != nil && errors.IsNotFound(err) {
		vipService = resources.NewVIPServiceForCR(elbServiceInst)
		// 设置owner，这里需要cr的scheme
		controllerutil.SetControllerReference(elbServiceInst, vipService, r.scheme)
		err = r.client.Create(context.TODO(), vipService)
		if err != nil {
			if !k8serrors.IsAlreadyExists(err) {
				logger.Error(err, "Failed to create new VIPService", "VIPService.Namespace", vipService.Namespace, "VIPService.Name", vipService.Name)
				// 需要重试
				r.makeFailedPhase(elbServiceInst, err)
				return err
			}
		}
		logger.Info("Creating a new VIPService successed.", "VIPService.Name", vipService.Name)
		// Service created successfully - return and requeue
		//return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		r.makeFailedPhase(elbServiceInst, err)
		logger.Error(err, "Failed to get VIPService.", "VIPService.Name", vipService.Name)
		return err
	} else {
		logger.Info("VIPService already exists", "VIPService.Namespace", vipService.Namespace, "VIPService.Name", vipService.Name)
	}

	// 判断vipendpoints是否存在，不存在就创建
	vipSvcEndpoints := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetVIPServiceName(elbServiceInst),
		Namespace: elbServiceInst.Namespace,
	}, vipSvcEndpoints)
	if err != nil && errors.IsNotFound(err) {
		vipSvcEndpoints = resources.NewVIPEndpointForCR(elbServiceInst)
		// 设置owner
		controllerutil.SetControllerReference(elbServiceInst, vipSvcEndpoints, r.scheme)
		err = r.client.Create(context.TODO(), vipSvcEndpoints)
		if err != nil {
			if !k8serrors.IsAlreadyExists(err) {
				logger.Error(err, "Failed to create new VIPEndpoints", "VIPEndpoints.Namespace", vipSvcEndpoints.Namespace, "VIPEndpoints.Name", vipSvcEndpoints.Name)
				// 需要重试
				r.makeFailedPhase(elbServiceInst, err)
				return err
			}
		}
		logger.Info("Creating a new VIPSvcEndpoints successed.", "VIPSvcEndpoints.Name", vipSvcEndpoints.Name)
		// 创建一个pod用于测试域名解析
		// dnsTestPod := resources.NewPodForCR(instance)
		// controllerutil.SetControllerReference(instance, dnsTestPod, r.scheme)
		// r.client.Create(context.TODO(), dnsTestPod)

		// Endpoints created successfully - return and requeue
		// return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		r.makeFailedPhase(elbServiceInst, err)
		logger.Error(err, "Failed to get VIPSvcEndpoints.", "VIPSvcEndpoints.Name", vipSvcEndpoints.Name)
		return err
	} else {
		logger.Info("VIPSvcEndpoints already exists", "VIPSvcEndpoints.Namespaec", vipSvcEndpoints.Namespace, "VIPSvcEndpoints.Name", vipSvcEndpoints.Name)
	}

	// 判断epservice是否存在，不存在就创建
	epService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetEPServiceName(elbServiceInst),
		Namespace: elbServiceInst.Namespace,
	}, epService)
	if err != nil && errors.IsNotFound(err) {
		epService = resources.NewEPServiceForCR(elbServiceInst)
		// 设置owner
		controllerutil.SetControllerReference(elbServiceInst, epService, r.scheme)

		err = r.client.Create(context.TODO(), epService)
		if err != nil {
			if !k8serrors.IsAlreadyExists(err) {
				logger.Error(err, "Failed to create new EPService", "EPService.Namespace", epService.Namespace, "EPService.Name", epService.Name)
				r.makeFailedPhase(elbServiceInst, err)
				return err
			}
		}
		logger.Info("Creating a new EPService successed.", "EPService.Name", epService.Name)
		// Service created successfully - return and requeue
		//return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		r.makeFailedPhase(elbServiceInst, err)
		logger.Error(err, "Failed to get EPService.", "EPService.Name", epService.Name)
		return err
	} else {
		logger.Info("EPService already exists", "EPService.Namespace", epService.Namespace, "EPService.Name", epService.Name)
	}

	// 尝试去获取epSvcEndPoints
	epSvcEndPoints := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resources.GetEPServiceName(elbServiceInst),
		Namespace: elbServiceInst.Namespace,
	}, epSvcEndPoints)
	if err != nil && errors.IsNotFound(err) {
		// 找不到继续
		logger.Error(err, "Not found EPSvcEndpoints so requeue.", "EPSvcEndpoints.Name", epSvcEndPoints.Name)
		return err
	} else if err == nil {
		// 修改ownerreference
		// 修改owner，这里只需要做一次
		err = controllerutil.SetControllerReference(elbServiceInst, epSvcEndPoints, r.scheme)
		if err != nil {
			logger.Error(err, "Set EPSvcEndpoints OwnerReferences failed.", "EPSvcEndpoints.Name", epSvcEndPoints.Name)
			r.makeFailedPhase(elbServiceInst, err)
			return err
		} else {
			// 更新
			err = r.client.Update(context.TODO(), epSvcEndPoints)
			if err != nil {
				r.makeFailedPhase(elbServiceInst, err)
				logger.Error(err, "Update EPSvcEndpoints failed.", "EPSvcEndpoints.Name", epSvcEndPoints.Name)
				return err
			}
		}
		logger.Info("found EPSvcEndpoints and update successed.", "EPSvcEndpoints.Name", epSvcEndPoints.Name)
	}

	if elbServiceInst.Status.Phase == k8sv1alpha1.ELBServiceCreating {
		addFinalizer(&elbServiceInst.ObjectMeta, finalizerName)
		elbServiceInst.Status.Phase = k8sv1alpha1.ELBServiceActive
		elbServiceInst.Status.LastUpdateTime = metav1.NewTime(time.Now())
		// 这个要放前面，不然会报错。
		r.makeActivePhase(elbServiceInst)

		if err := r.client.Update(context.TODO(), elbServiceInst); err != nil {
			logger.Error(err, "update ELBService add Finializer failed.")
		} else {
			logger.Info("update ELBService object add Finializer successed.")
		}
	} else {
		r.makeActivePhase(elbServiceInst)
	}

	return nil
}

func (r *ReconcileELBService) deleteService(elbServiceInst *k8sv1alpha1.ELBService) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)
	logger.Info("delete service")

	// 和真正的elb解绑

	// 删除finializer标志
	if removeFinalizer(&elbServiceInst.ObjectMeta, finalizerName) {
		if err := r.client.Update(context.TODO(), elbServiceInst); err != nil {
			logger.Error(err, "update ELBService object remove Finializer failed.")
		}
		logger.Info("update ELBService object remove Finializer successed.")
	}

	// 当前面的finializer被删除后，对象已经不存在了，所以这里会报错
	// if err := r.makeTerminatingPhase(elbServiceInst); err != nil {
	// 	logger.Error(err, "make terminating phase failed.")
	// 	return err
	// }
	logger.Info("delete service successed.")
	return nil
}

func (r *ReconcileELBService) updateServiceStatus(elbServiceInst *k8sv1alpha1.ELBService, phase k8sv1alpha1.ELBServicePhase, reason error) error {
	logger := log.WithValues("Request.Namespace", elbServiceInst.Namespace, "Request.Name", elbServiceInst.Name)

	elbServiceInst.Status.Reason = ""
	if reason != nil {
		elbServiceInst.Status.Reason = reason.Error()
	}

	elbServiceInst.Status.Phase = phase
	elbServiceInst.Status.LastUpdateTime = metav1.NewTime(time.Now())
	// 修改状态，我想知道更新状态是否会导致调和被调用
	err := r.client.Status().Update(context.TODO(), elbServiceInst)
	if err != nil {
		logger.Error(err, "Failed to update ELBService status.")
		return err
	}
	logger.Info("Update ELBService status successed")
	return nil
}
