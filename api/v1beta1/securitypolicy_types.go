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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LabelSelectors struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SecurityPolicyRule defines rules
type SecurityPolicyRule struct {
	Action            string           `json:"action"`
	Description       string           `json:"description"`
	Priority          int64            `json:"priority"`
	SrcIpRanges       []string         `json:"srcIpRanges,omitempty"`
	NodePoolSelectors []LabelSelectors `json:"nodePoolSelectors,omitempty"`
}

// SecurityPolicySpec defines the desired state of SecurityPolicy
type SecurityPolicySpec struct {
	// +kubebuilder:validation:MinLength=1
	Name          string               `json:"name,omitempty"`
	Description   string               `json:"description,omitempty"`
	DefaultAction string               `json:"defaultAction,omitempty"`
	Rules         []SecurityPolicyRule `json:"rules,omitempty"`
}

// SecurityPolicyStatus defines the observed state of SecurityPolicy
type SecurityPolicyStatus struct {
	// +kubebuilder:validation:MinLength=1
	Name          string               `json:"name,omitempty"`
	Description   string               `json:"description,omitempty"`
	DefaultAction string               `json:"defaultAction,omitempty"`
	Rules         []SecurityPolicyRule `json:"rules,omitempty"`
	Condition     string               `json:"condition,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityPolicy is the Schema for the securitypolicies API
type SecurityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityPolicySpec   `json:"spec,omitempty"`
	Status SecurityPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityPolicyList contains a list of SecurityPolicy
type SecurityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityPolicy{}, &SecurityPolicyList{})
}
