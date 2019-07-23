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
	"time"

	"github.com/go-logr/logr"
	cloudarmorv1beta1 "github.com/h-r-k-matsumoto/security-policy-operator/api/v1beta1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// SecurityPolicyReconciler reconciles a SecurityPolicy object
type SecurityPolicyReconciler struct {
	client.Client
	Log logr.Logger
}

// Reconcile logic
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=node,verbs=get;list;watch

func (r *SecurityPolicyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("securitypolicy", req.NamespacedName)
	log.Info("Reconcile start.")

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
	instance.Status.Name = instance.Spec.Name
	instance.Status.Description = instance.Spec.Description
	instance.Status.DefaultAction = instance.Spec.DefaultAction
	instance.Status.Rules = instance.Spec.Rules
	if !containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, myFinalizerName)
	}
	nodeCalculator := &NodeCalculator{Log: r.Log, Reconciler: r}
	instance, err = nodeCalculator.Calculate(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	api := SecurityPolicyAPI{Log: r.Log}
	err = retry(
		func() error {
			gceCurrentInstance, err := api.Get(ctx, instance.Status.Name)
			if err != nil {
				return err
			}
			if gceCurrentInstance == nil {
				log.Info("Create Security Policy")
				if err := api.Create(ctx, &instance.Status); err != nil {
					return err
				}
			} else {
				log.Info("Apply Security Policy")
				if err := api.Apply(ctx, &instance.Status, gceCurrentInstance); err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	instance.Status.Condition = "security operator updated."
	if err := r.Update(ctx, instance); err != nil {
		return reconcile.Result{RequeueAfter: 5 * time.Second}, err
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
	err := api.Delete(ctx, instance.Status.Name)
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

func retry(fn func() error) error {
	i := 3
	sleepTime := 1 * time.Second
	var err error = nil
	for i > 0 {
		if err = fn(); err != nil {
			i--
			if i > 0 {
				time.Sleep(sleepTime)
				continue
			}
			return err
		}
		return nil
	}
	return err
}
