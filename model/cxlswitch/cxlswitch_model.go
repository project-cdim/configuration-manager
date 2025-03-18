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
        
package cxlswitch_model

import (
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

// CXLSwitch is a CXLSwitch structure.
type CXLSwitch struct {
	Properties map[string]any
	Resources  resource_model.ResourceList
}

// NewCXLSwitch is Constructor of CXLSwitch structure.
// Returns an object with all elements having empty values.
func NewCXLSwitch() CXLSwitch {
	return CXLSwitch{
		Properties: map[string]any{},
		Resources:  resource_model.ResourceList{},
	}
}

// Validate reports whether the CXLSwitch instance is valid.
//
// This method checks the validity of the CXLSwitch instance by verifying that the "id" property exists
// and is a non-empty string. It is a crucial step to ensure that each CXLSwitch instance has a unique identifier
// before proceeding with operations that require a valid CXLSwitch.
//
// Returns:
//
//	bool: True if the CXLSwitch has a valid "id" property, false otherwise.
func (c *CXLSwitch) Validate() bool {
	id, ok := c.Properties["id"].(string)
	if !ok || len(id) <= 0 {
		return false
	}

	return true
}

// ToObject creates and returns a map with elements of id, resources, and any optional elements.
//
// This method first validates the CXLSwitch instance. If the instance is not valid, it returns nil.
// Upon successful validation, it proceeds to construct a map (`res`) initialized with the CXLSwitch's properties.
// It then adds the resources, converted to their object form, under the "resources" key in the map.
// The resulting map, which now includes the CXLSwitch's properties and its resources, is returned.
//
// Returns:
//
//	map[string]any: A map representation of the CXLSwitch, including its properties and resources, or nil if the CXLSwitch is invalid.
func (c *CXLSwitch) ToObject() map[string]any {
	if !c.Validate() {
		return nil
	}

	res := c.Properties
	res["resources"] = c.Resources.ToObject()

	return res
}

// ToObject4Chassis creates and returns a map with elements of id, resources, and any optional elements,
// specifically formatted for ChassisObject consumption.
//
// This method first validates the CXLSwitch instance. If the instance is not valid, it returns nil.
// Upon successful validation, it constructs a map (`res`) initialized with the CXLSwitch's properties.
// Unlike ToObject, this method does not include the resources in the returned map, making it specifically
// tailored for ChassisObject consumption, where resources might be handled differently.
//
// Returns:
//
//	map[string]any: A map representation of the CXLSwitch, formatted for ChassisObject, or nil if the CXLSwitch is invalid.
func (c *CXLSwitch) ToObject4Chassis() map[string]any {
	if !c.Validate() {
		return nil
	}

	res := c.Properties

	return res
}
