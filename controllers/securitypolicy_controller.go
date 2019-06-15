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
	"fmt"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("delete object.")
		// The object is being deleted
		// add finializer logic.
		return ctrl.Result{}, nil
	}

	api := SecurityPolicyAPI{}
	gcp_instance, err := api.Get(ctx, instance.Spec.Name)
	log.Info("====================")
	log.Info(fmt.Sprintf("%v", gcp_instance))
	log.Info("====================")
	if gcp_instance == nil {
		_, err := api.Create(ctx, &instance.Spec)
		if err != nil {
			log.Error(err, "error")
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
