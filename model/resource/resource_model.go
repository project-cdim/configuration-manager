// Copyright (C) 2025 NEC Corporation.
// 
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
        
package resource_model

import (
	annotation_model "github.com/project-cdim/configuration-manager/model/annotation"
)

// Resource is a resource structure.
type Resource struct {
	Device           map[string]any
	Annotation       annotation_model.Annotation
	ResourceGroupIDs []string
	NodeIDs          []string
	Detected         bool
}

// NewResource is the constructor for the Resource structure.
//
// This function initializes a Resource struct with all elements having empty values.
// It is useful for creating a Resource instance ready to be populated with device information,
// annotations, resource group IDs, node IDs, and detection status.
//
// Returns:
//
//	Resource: A new instance of Resource with empty device information, annotations,
//	          resource group IDs, node IDs, and detection status set to false.
func NewResource() Resource {
	return Resource{
		Device:           map[string]any{},
		Annotation:       annotation_model.NewAnnotation(),
		ResourceGroupIDs: []string{},
		NodeIDs:          []string{},
		Detected:         false,
	}
}

// Validate reports whether the receiver is valid.
//
// This method checks the validity of the Resource instance by verifying that at least one
// of the Device, Annotation.Properties, ResourceGroupIDs, or NodeIDs is not empty.
// It ensures that the Resource instance has enough information to be considered valid.
//
// Returns:
//
//	bool: True if the Resource instance is valid, false otherwise.
func (r *Resource) Validate() bool {
	if len(r.Device) == 0 && len(r.Annotation.Properties) == 0 && len(r.ResourceGroupIDs) == 0 && len(r.NodeIDs) == 0 {
		return false
	}
	return true
}

// ToObject creates and returns a map with elements
// of device, annotation, resourceGroupIDs, nodeIds, and detected.
//
// This method first validates the Resource instance. If the instance is not valid, it returns nil.
// Upon successful validation, it constructs a map (`res`) initialized with the Resource's device information,
// annotations (formatted specifically for Resource), resource group IDs, node IDs, and detection status.
// The resulting map is returned and includes all necessary information about the Resource.
//
// Returns:
//
//	map[string]any: A map representation of the Resource, including its device information,
//	                annotations, resource group IDs, node IDs, and detection status, or nil if the Resource is invalid.
func (r *Resource) ToObject() map[string]any {
	if !r.Validate() {
		return nil
	}

	return map[string]any{
		"device":           r.Device,
		"annotation":       r.Annotation.ToObject(),
		"resourceGroupIDs": r.ResourceGroupIDs,
		"nodeIDs":          r.NodeIDs,
		"detected":         r.Detected,
	}
}

// ToObject4Node creates for nodeObject and returns a map with elements
// of device, annotation, resourceGroupIDs, and detected.
//
// This method is similar to ToObject but tailored for Node objects. It first validates the Resource instance.
// If the instance is not valid, it returns nil. Upon successful validation, it constructs a map (`res`) initialized
// with the Resource's device information, annotations (formatted specifically for Resource), resource group IDs,
// and detection status. The resulting map is tailored for consumption by Node objects.
//
// Returns:
//
//	map[string]any: A map representation of the Resource, formatted specifically for Node consumption,
//	                including its device information, annotations, resource group IDs, and detection status,
//	                or nil if the Resource is invalid.
func (r *Resource) ToObject4Node() map[string]any {
	if !r.Validate() {
		return nil
	}

	return map[string]any{
		"device":           r.Device,
		"annotation":       r.Annotation.ToObject(),
		"resourceGroupIDs": r.ResourceGroupIDs,
		"detected":         r.Detected,
	}
}

// ToObject4Unused converts the Resource instance into a map representation suitable for unusing operations.
// It first validates the Resource instance, and if validation fails, it returns nil.
// If validation succeeds, it returns a map containing the device, annotation, and resourceGroupIDs fields.
//
// Returns:
//
//	A map[string]any representing the Resource instance if validation succeeds, otherwise nil.
func (r *Resource) ToObject4Unused() map[string]any {
	if !r.Validate() {
		return nil
	}

	return map[string]any{
		"device":           r.Device,
		"annotation":       r.Annotation.ToObject(),
		"resourceGroupIDs": r.ResourceGroupIDs,
	}
}
