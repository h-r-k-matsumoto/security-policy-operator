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

	cloudarmorv1beta1 "github.com/h-r-k-matsumoto/security-policy-operator/api/v1beta1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"

	compute "google.golang.org/api/compute/v1"
)

// SecurityPolicyAPI is Google Compute SecurityPolicy API structure.
type SecurityPolicyAPI struct {
}

// Get returns search results by name
func (api *SecurityPolicyAPI) Get(ctx context.Context, name string) (*compute.SecurityPolicy, error) {
	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return nil, err
	}

	req := service.List(credentials.ProjectID)
	list, err := req.Filter(fmt.Sprintf("name=%s", name)).MaxResults(1).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 404 {
				//already deleted.
				return nil, nil
			}
		}
		return nil, err
	}
	if len(list.Items) > 0 {
		return list.Items[0], nil
	}
	return nil, nil
}

// Create calls Security Policy Insert API
func (api *SecurityPolicyAPI) Create(ctx context.Context, spec *cloudarmorv1beta1.SecurityPolicySpec) (*compute.SecurityPolicy, error) {

	service, credentials, err := newSecurityPoliciesService(ctx)
	if err != nil {
		return nil, err
	}

	rb := customResourceToSecurityPolicy(spec)
	resp, err := service.Insert(credentials.ProjectID, rb).Context(ctx).Do()
	fmt.Printf("%#v\n", resp)
	fmt.Printf("%#v\n", err)
	if err != nil {
		return nil, nil
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)
	return nil, nil
}

//
func (api *SecurityPolicyAPI) Apply(ctx context.Context, spec *cloudarmorv1beta1.SecurityPolicySpec, policy *compute.SecurityPolicy) (*compute.SecurityPolicy, error) {
	//diff

	return nil, nil
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

// customResourceToSecurityPolicy convert cloudarmorv1beta1.SecurityPolicySpec to compute.SecurityPolicy.
func customResourceToSecurityPolicy(spec *cloudarmorv1beta1.SecurityPolicySpec) *compute.SecurityPolicy {
	rules := make([]*compute.SecurityPolicyRule, len(spec.Rules))
	for i, _ := range rules {
		rules[i] = &compute.SecurityPolicyRule{
			Action:      spec.Rules[i].Action,
			Description: spec.Rules[i].Description,
			Priority:    spec.Rules[i].Priority,
			Match: &compute.SecurityPolicyRuleMatcher{
				VersionedExpr: "SRC_IPS_V1",
				Config: &compute.SecurityPolicyRuleMatcherConfig{
					SrcIpRanges: spec.Rules[i].SrcIpRanges,
				},
			},
		}
	}
	rules = append(rules, &compute.SecurityPolicyRule{
		Action:      spec.DefaultAction,
		Description: "This is default action",
		Priority:    2147483647,
		Match: &compute.SecurityPolicyRuleMatcher{
			VersionedExpr: "SRC_IPS_V1",
			Config: &compute.SecurityPolicyRuleMatcherConfig{
				SrcIpRanges: []string{"*"},
			},
		},
	})
	rb := &compute.SecurityPolicy{
		Name:        spec.Name,
		Description: spec.Description,
		Rules:       rules,
	}
	return rb
}

// applyToSecurityPolicy return bool that  is true, if a change has occurred. And return SecurityPolicy to apply.
func applyToSecurityPolicy(spec *cloudarmorv1beta1.SecurityPolicySpec) (bool, *compute.SecurityPolicy) {
	return false, nil
}
