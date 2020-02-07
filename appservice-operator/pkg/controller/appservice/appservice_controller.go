// Package appservice for implement
package appservice

import (
	"context"
	"encoding/json"
	"time"

	appservicecontrollerv1 "calmwu.org/appservice-operator/pkg/apis/appservicecontroller/v1"
	v1 "calmwu.org/appservice-operator/pkg/apis/appservicecontroller/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const appServiceFinalizer = "finalizer.appservice"

var log = logf.Log.WithName("controller_appservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AppService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAppService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("appservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AppService
	err = c.Watch(&source.Kind{Type: &appservicecontrollerv1.AppService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AppService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appservicecontrollerv1.AppService{},
	})
	if err != nil {
		return err
	}

	// 添加需要监听的资源类型
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appservicecontrollerv1.AppService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAppService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppService{}

// ReconcileAppService reconciles a AppService object
type ReconcileAppService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AppService object and makes changes based on the state read
// and what is in the AppService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("calmwu-log", "request:", request)

	// Fetch the AppService instance
	instance := &appservicecontrollerv1.AppService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("---------------NotFound-------------")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		// 表明crd对象被删除
		now := time.Now()
		nowName := now.Format("2006-01-02 15:04:05 -0700")
		deletionTimestampName := instance.DeletionTimestamp.Format("2006-01-02 15:04:05 -0700")
		reqLogger.Info("**calmwu-log**", "deletionTimestampName:", deletionTimestampName, "now:", nowName)

		// 如果包含appServiceFinalizer信息，就开始做手工清理工作
		if contains(instance.GetFinalizers(), appServiceFinalizer) {
			reqLogger.Info("**calmwu-log** Successfully finalized AppService")
		}

		// 同时清理Finalizers所有内容
		instance.SetFinalizers([]string{})
		// 更新
		r.client.Update(context.TODO(), instance)

		return reconcile.Result{}, err
	}

	// 给cr加上finalizer
	if !contains(instance.GetFinalizers(), appServiceFinalizer) {
		reqLogger.Info("**calmwu-log** Adding Finalizer for the AppService")
		instance.SetFinalizers(append(instance.GetFinalizers(), appServiceFinalizer))
		// 这个不可修改
		//var sec int64 = 60
		//instance.SetDeletionGracePeriodSeconds(&sec)
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update AppService with finalizer")
			return reconcile.Result{}, err
		}
	}

	//instanceInfo := fmt.Sprintf("%#v", instance)
	//reqLogger.Info("calmwu-log", "AppService instance info:", instanceInfo)

	//controllerutil.SetControllerReference
	deploy := &appsv1.Deployment{}
	// deployment名字和appservice相同
	if err := r.client.Get(context.TODO(), request.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		// 如果deployment不存在，创建
		deploy = newDeploy(instance)
		if err := r.client.Create(context.TODO(), deploy); err != nil {
			reqLogger.Error(err, "**calmwu-log** Create deployment failed.")
			return reconcile.Result{}, err
		} else {
			reqLogger.Info("calmwu-log ---------- create deployment")
		}

		// 创建service
		service := newService(instance)
		if err = r.client.Create(context.TODO(), service); err != nil {
			reqLogger.Error(err, "**calmwu-log** Create service failed.")
			return reconcile.Result{}, err
		}

		// 将当前的信息写入到appservice的annotation中
		reqLogger.Info("**calmwu-log**", "instance.spec", instance.Spec)
		data, _ := json.Marshal(instance.Spec)
		if instance.Annotations != nil {
			instance.Annotations["spec"] = string(data)
		} else {
			instance.Annotations = map[string]string{
				"spec": string(data),
			}
		}

		// 然后修改instance，将spec加入到instance中
		if err := r.client.Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "calmwu-log", "update instance space failed.", instance.Spec)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	}

	//reqLogger.Info("**calmwu-log**", "annotation.spec", instance.Annotations["spec"])

	// 判断instance状态是否变化，带过来的annotaion.Spec来判断
	annotationSpecStr := instance.Annotations["spec"]
	specStr, _ := json.Marshal(instance.Spec)

	// 判断是否需要更新，instance中spec和annotation中带的内容不对应，需要更新
	if annotationSpecStr != string(specStr) {
		//diffStr := cmp.Diff(instance.Spec, oldSpec)

		reqLogger.Info("**calmwu-log** update")
		reqLogger.Info("**calmwu-log**", "annotationSpecStr", annotationSpecStr)
		reqLogger.Info("**calmwu-log**", "specStr", string(specStr))

		instance.Annotations["spec"] = string(specStr)
		if err := r.client.Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "-------calmwu-log-----------", "update instance spec failed.", instance.Spec)
			return reconcile.Result{}, nil
		}

		// 将这个spec内容设置到deployment中
		newDP := newDeploy(instance)
		newDP.GetNamespace()
		oldDP := &appsv1.Deployment{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldDP); err != nil {
			reqLogger.Error(err, "Get deployment:", request.NamespacedName, " failed.")
			return reconcile.Result{}, err
		}
		oldDP.Spec = newDP.Spec
		// 更新
		if err := r.client.Update(context.TODO(), oldDP); err != nil {
			reqLogger.Error(err, "Update deployment:", request.NamespacedName, " failed.")
			return reconcile.Result{}, err
		}

		newSvc := newService(instance)
		oldSvc := &corev1.Service{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldSvc); err != nil {
			return reconcile.Result{}, err
		}
		oldSvc.Spec = newSvc.Spec
		if err := r.client.Update(context.TODO(), oldSvc); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	reqLogger.Info("---------------Reconcile exit-------------")

	return reconcile.Result{}, nil
}

// 构建一个新的deployment
func newDeploy(cr *appservicecontrollerv1.AppService) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}

	selector := &metav1.LabelSelector{MatchLabels: labels}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,

			// 设置deployment的归属资源
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   appservicecontrollerv1.SchemeGroupVersion.Group,
					Version: appservicecontrollerv1.SchemeGroupVersion.Version,
					Kind:    "AppService",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.Size,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: newContainers(cr),
				},
			},
			Selector: selector,
		},
	}
}

func newContainers(app *v1.AppService) []corev1.Container {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range app.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	return []corev1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.Image,
			Resources:       app.Spec.Resources,
			Ports:           containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             app.Spec.Envs,
		},
	}
}

func newService(app *v1.AppService) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "AppService",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: app.Spec.Ports,
			Selector: map[string]string{
				"app": app.Name,
			},
		},
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
