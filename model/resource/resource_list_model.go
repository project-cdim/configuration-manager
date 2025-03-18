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
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
)

// ResourceList is a list of resources.
type ResourceList struct {
	Resources []Resource
}

// NewResourceList is the constructor for the ResourceList structure.
//
// This function initializes a ResourceList struct with an empty slice of Resource.
// It is useful for creating a ResourceList ready to be populated with Resource instances.
//
// Returns:
//
//	ResourceList: A new instance of ResourceList with an empty slice of Resource.
func NewResourceList() ResourceList {
	return ResourceList{
		Resources: []Resource{},
	}
}

// ToObject creates and returns a slice of maps with elements
// of device, annotation, and resourceGroupIDs from a resource list.
//
// This method iterates over each Resource in the ResourceList, validates it,
// and then converts it into a map object. Each map object represents a Resource
// and includes keys for device, annotation, and resourceGroupIDs. These map objects
// are then collected into a slice and returned. This is useful for converting a list
// of Resource objects into a more generic data structure that can be easily manipulated
// or serialized.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated Resource.
func (rl *ResourceList) ToObject() []map[string]any {
	res := []map[string]any{}
	for _, resource := range rl.Resources {
		if resource.Validate() {
			res = append(res, resource.ToObject())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. resource(%v)", resource))
		}
	}

	return res
}

// ToObject4Node creates for nodeObject and returns a slice of maps with elements
// of device, annotation, and resourceGroupIDs from a resource list, specifically formatted for Node consumption.
//
// This method is similar to ToObject but tailored for Node objects. It iterates over each Resource in the ResourceList,
// validates it, and then converts it into a map object formatted specifically for Node consumption. Each map object
// includes keys for device, annotation, and resourceGroupIDs. These map objects are then collected into a slice and returned.
// This method is useful when a Node-specific representation of the ResourceList is required.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated Resource, formatted for Node consumption.
func (rl *ResourceList) ToObject4Node() []map[string]any {
	res := []map[string]any{}
	for _, resource := range rl.Resources {
		if resource.Validate() {
			res = append(res, resource.ToObject4Node())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. resource(%v)", resource))
		}
	}

	return res
}

// ToObject4Unused converts the ResourceList instance into a slice of map representations suitable for unusing operations.
// It iterates over each Resource in the ResourceList, validates it, and if validation succeeds,
// it converts the Resource to a map using the ToObject4Unused method of the Resource struct.
// The resulting maps are collected into a slice.
//
// Returns:
//
//	A slice of map[string]any, where each map represents a validated Resource instance.
func (rl *ResourceList) ToObject4Unused() []map[string]any {
	res := []map[string]any{}
	for _, resource := range rl.Resources {
		if resource.Validate() {
			res = append(res, resource.ToObject4Unused())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. resource(%v)", resource))
		}
	}

	return res
}
