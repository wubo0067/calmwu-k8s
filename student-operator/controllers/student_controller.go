/*
Copyright 2020 calmwu.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	studentcontrollerv1 "calmwu.org/student-operator/api/v1"
)

// StudentReconciler reconciles a Student object
type StudentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=studentcontroller.calmwu.org,resources=students,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=studentcontroller.calmwu.org,resources=students/status,verbs=get;update;patch

func (r *StudentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("student", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *StudentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&studentcontrollerv1.Student{}).
		Complete(r)
}
