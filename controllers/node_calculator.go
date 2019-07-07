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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeCalculator struct {
	Log        logr.Logger
	Reconciler *SecurityPolicyReconciler
}

func (n *NodeCalculator) Calculate(instance *cloudarmorv1beta1.SecurityPolicy) (*cloudarmorv1beta1.SecurityPolicy, error) {
	for i, rule := range instance.Status.Rules {
		if rule.NodePoolSelectors == nil {
			continue
		}
		labels := map[string]string{}
		for _, selector := range rule.NodePoolSelectors {
			labels[selector.Key] = selector.Value
		}
		addresses, err := n.List(labels)
		if err != nil {
			return nil, err
		}
		instance.Status.Rules[i].SrcIpRanges = addresses
	}
	return instance, nil
}

// List is returned node external ip list.
func (n *NodeCalculator) List(labels map[string]string) ([]string, error) {
	log := n.Log.WithValues("gcp_securitypolicy", "node_handler")
	log.Info("Node Address List")
	ctx := context.Background()
	addresses := []string{}

	nodelist := &corev1.NodeList{}
	filter := client.MatchingLabels(labels)
	err := n.Reconciler.List(ctx, nodelist, filter)
	if err != nil {
		return addresses, err
	}
	for _, node := range nodelist.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeExternalIP {
				addresses = append(addresses, address.Address)
			}
		}
	}
	return addresses, nil
}
