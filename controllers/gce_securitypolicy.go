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
	context "context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"

	cloudarmorv1beta1 "github.com/h-r-k-matsumoto/security-policy-operator/api/v1beta1"
	compute "google.golang.org/api/compute/v1"
)

// SecurityPolicyAPI is Google Compute SecurityPolicy API structure.
type SecurityPolicyAPI struct {
	Log logr.Logger
}

// Get returns search results by id
func (api *SecurityPolicyAPI) Get(ctx context.Context, name string) (*compute.SecurityPolicy, error) {
	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return nil, err
	}
	policy, err := service.Get(credentials.ProjectID, name).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 404 {
				//already deleted.
				return nil, nil
			}
		}
		return nil, err
	}
	return policy, nil
}

// Create calls Security Policy Insert API
func (api *SecurityPolicyAPI) Create(ctx context.Context, spec *cloudarmorv1beta1.SecurityPolicyStatus) error {
	log := api.Log.WithValues("gcp_securitypolicy", spec.Name)

	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return err
	}

	log.Info("Insert SecurityPolicy")
	rb := customResourceToSecurityPolicy(spec)
	req := service.Insert(credentials.ProjectID, rb).Context(ctx)
	if _, err := req.Do(); err != nil {
		return nil
	}
	return nil
}

//
func (api *SecurityPolicyAPI) Apply(ctx context.Context, spec *cloudarmorv1beta1.SecurityPolicyStatus, current *compute.SecurityPolicy) error {
	log := api.Log.WithValues("gcp_securitypolicy", spec.Name)

	update := customResourceToSecurityPolicy(spec)
	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return err
	}

	// generate priority map.
	currentPriorityMap := make(map[int64]*compute.SecurityPolicyRule, len(current.Rules))
	for _, rule := range current.Rules {
		currentPriorityMap[rule.Priority] = rule
	}

	updatePriorityMap := make(map[int64]*compute.SecurityPolicyRule, len(update.Rules))
	for _, rule := range update.Rules {
		updatePriorityMap[rule.Priority] = rule
	}

	//Rules check
	for priority, updateRule := range updatePriorityMap {
		if currentRule, ok := currentPriorityMap[priority]; ok {
			if updateRule.Action != currentRule.Action || updateRule.Description != currentRule.Description || !reflect.DeepEqual(updateRule.Match.Config.SrcIpRanges, currentRule.Match.Config.SrcIpRanges) {
				log.Info(fmt.Sprintf("Patch SecurityPolicy Rule [ priority=%d ]", updateRule.Priority))
				req := service.PatchRule(credentials.ProjectID, update.Name, updateRule).Context(ctx).Priority(updateRule.Priority)
				if _, err := req.Do(); err != nil {
					return err
				}
			}
		} else {
			log.Info(fmt.Sprintf("Add SecurityPolicy Rule [ priority=%d ]", updateRule.Priority))
			req := service.AddRule(credentials.ProjectID, update.Name, updateRule).Context(ctx)
			if _, err := req.Do(); err != nil {
				return err
			}
		}
	}
	for priority, currentRule := range currentPriorityMap {
		if _, ok := updatePriorityMap[priority]; !ok {
			log.Info(fmt.Sprintf("Remove SecurityPolicy Rule [ priority=%d ]", currentRule.Priority))
			req := service.RemoveRule(credentials.ProjectID, update.Name).Context(ctx).Priority(currentRule.Priority)
			if _, err := req.Do(); err != nil {
				return err
			}
		}
	}

	if update.Name != current.Name || update.Description != current.Description {
		log.Info("Patch SecurityPolicy")
		update.Fingerprint = current.Fingerprint
		update.Id = current.Id
		update.Rules = nil
		req := service.Patch(credentials.ProjectID, update.Name, update).Context(ctx)
		if _, err := req.Do(); err != nil {
			return err
		}
	}
	return nil
}

// Delete is delete security policy.
func (api *SecurityPolicyAPI) Delete(ctx context.Context, name string) error {
	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return err
	}

	_, err = service.Delete(credentials.ProjectID, name).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 404 {
				//already deleted.
				return nil
			}
		}
		return err
	}
	return nil
}

// newSecurityPoliciesService returns SecurityPoliciesService and default credential.
func newSecurityPoliciesService(ctx context.Context) (*compute.SecurityPoliciesService, *google.Credentials, error) {
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, nil, err
	}
	computeService, err := compute.New(c)
	if err != nil {
		return nil, nil, err
	}
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		return nil, nil, err
	}

	return computeService.SecurityPolicies, credentials, nil
}

// customResourceToSecurityPolicyRule convert cloudarmorv1beta1.SecurityPolicyRule to compute.SecurityPolicyRule
func customResourceToSecurityPolicyRule(rule *cloudarmorv1beta1.SecurityPolicyRule) *compute.SecurityPolicyRule {
	return &compute.SecurityPolicyRule{
		Action:      rule.Action,
		Description: rule.Description,
		Priority:    rule.Priority,
		Match: &compute.SecurityPolicyRuleMatcher{
			VersionedExpr: "SRC_IPS_V1",
			Config: &compute.SecurityPolicyRuleMatcherConfig{
				SrcIpRanges: rule.SrcIpRanges,
			},
		},
	}
}

// defaultSecurityPolicyRule generates default security policy rule.
func defaultSecurityPolicyRule(spec *cloudarmorv1beta1.SecurityPolicyStatus) *compute.SecurityPolicyRule {
	return &compute.SecurityPolicyRule{
		Action:      spec.DefaultAction,
		Description: "This is default action",
		Priority:    2147483647,
		Match: &compute.SecurityPolicyRuleMatcher{
			VersionedExpr: "SRC_IPS_V1",
			Config: &compute.SecurityPolicyRuleMatcherConfig{
				SrcIpRanges: []string{"*"},
			},
		},
	}
}

// customResourceToSecurityPolicy convert cloudarmorv1beta1.SecurityPolicyStatus to compute.SecurityPolicy.
func customResourceToSecurityPolicy(spec *cloudarmorv1beta1.SecurityPolicyStatus) *compute.SecurityPolicy {
	rules := make([]*compute.SecurityPolicyRule, len(spec.Rules))
	for i, _ := range rules {
		rules[i] = customResourceToSecurityPolicyRule(&spec.Rules[i])
	}
	rules = append(rules, defaultSecurityPolicyRule(spec))
	rb := &compute.SecurityPolicy{
		Name:        spec.Name,
		Description: spec.Description,
		Rules:       rules,
	}
	return rb
}
