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
	cloudarmorv1beta1 "github.com/h-r-k-matsumoto/security-policy-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// SecurityPolicyNodeReconciler reconciles a SecurityPolicy object
type SecurityPolicyNodeReconciler struct {
	client.Client
	Log logr.Logger
}

// Reconcile logic
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudarmor.matsumo.dev,resources=securitypolicies/status,verbs=get;update;patch
func (r *SecurityPolicyNodeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("securitypolicy(node)", req.NamespacedName)
	log.Info("Reconcile start.")

	instance := &cloudarmorv1beta1.SecurityPolicyList{}

	err := r.List(ctx, instance)
	if err != nil {
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	for _, policy := range instance.Items {
		for _, rule := range policy.Spec.Rules {
			if rule.NodePoolSelectors != nil && len(rule.NodePoolSelectors) > 0 {
				policy.Status.Condition = "node event update"
				if err := r.Update(ctx, &policy); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager is reconcile control.
func (r *SecurityPolicyNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		WithEventFilter(&NodeEventPredicate{}).
		Complete(r)
}
