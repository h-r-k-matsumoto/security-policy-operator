/*

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
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cloudarmorv1beta1 "github.com/h-r-k-matsumoto/security-policy-operator/api/v1beta1"
)

// SecurityPolicyReconciler reconciles a SecurityPolicy object
type SecurityPolicyReconciler struct {
	client.Client
	Log logr.Logger
}

// Reconcile logic.
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies/status,verbs=get;update;patch
func (r *SecurityPolicyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("securitypolicy", req.NamespacedName)

	// Fetch the Gcs instance
	instance := &cloudarmorv1beta1.SecurityPolicy{}

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	myFinalizerName := "securitypolicy.finalizer.cloudarmor.matsumo.dev"

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("delete object.")
		if containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
			// our finalizer is present, so lets handle our external dependency
			if err := r.deleteExternalDependency(instance); err != nil {
				return reconcile.Result{}, err
			}
			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, myFinalizerName)
		if err := r.Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	api := SecurityPolicyAPI{Log: r.Log}
	gceCurrentInstance, err := api.Get(ctx, instance.Spec.Name)
	if gceCurrentInstance == nil {
		log.Info("Create Security Policy")
		if err := api.Create(ctx, &instance.Spec); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		log.Info("Apply Security Policy")
		if err := api.Apply(ctx, &instance.Spec, gceCurrentInstance); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager is reconcile control.
func (r *SecurityPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudarmorv1beta1.SecurityPolicy{}).
		Complete(r)
}

//  delete dependency bucket.
func (r *SecurityPolicyReconciler) deleteExternalDependency(instance *cloudarmorv1beta1.SecurityPolicy) error {
	ctx := context.Background()
	api := SecurityPolicyAPI{Log: r.Log}
	err := api.Delete(ctx, instance.Spec.Name)
	return err
}

// Helper functions to check string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// Helper functions to remove string from slice.
func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
