// +build !ignore_autogenerated

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

// autogenerated by controller-gen object, do not modify manually

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityPolicy) DeepCopyInto(out *SecurityPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityPolicy.
func (in *SecurityPolicy) DeepCopy() *SecurityPolicy {
	if in == nil {
		return nil
	}
	out := new(SecurityPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecurityPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityPolicyList) DeepCopyInto(out *SecurityPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecurityPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityPolicyList.
func (in *SecurityPolicyList) DeepCopy() *SecurityPolicyList {
	if in == nil {
		return nil
	}
	out := new(SecurityPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecurityPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityPolicySpec) DeepCopyInto(out *SecurityPolicySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityPolicySpec.
func (in *SecurityPolicySpec) DeepCopy() *SecurityPolicySpec {
	if in == nil {
		return nil
	}
	out := new(SecurityPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityPolicyStatus) DeepCopyInto(out *SecurityPolicyStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityPolicyStatus.
func (in *SecurityPolicyStatus) DeepCopy() *SecurityPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(SecurityPolicyStatus)
	in.DeepCopyInto(out)
	return out
}
