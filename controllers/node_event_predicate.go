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
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type NodeEventPredicate struct {
}

// Create returns true if the Create event should be processed
func (p *NodeEventPredicate) Create(event.CreateEvent) bool {
	return true
}

// Delete returns true if the Delete event should be processed
func (p *NodeEventPredicate) Delete(event.DeleteEvent) bool {
	return true
}

// Update returns true if the Update event should be processed
func (p *NodeEventPredicate) Update(event.UpdateEvent) bool {
	return false
}

// Generic returns true if the Generic event should be processed
func (p *NodeEventPredicate) Generic(event.GenericEvent) bool {
	return false
}
