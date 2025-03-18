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
        
package chassis_model

import (
	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

// Chassis is a Chassis structure.
type Chassis struct {
	Properties  map[string]any
	Resources   resource_model.ResourceList
	CXLSwitches cxlswitch_model.CXLSwitchList
}

// NewChassis is a constructor for the Chassis structure.
//
// This function initializes a Chassis struct with all of its elements set to their zero values.
// Specifically, it sets Properties to an empty map, Resources to an empty ResourceList,
// and CXLSwitches to an empty CXLSwitchList. This is useful for creating a new Chassis instance
// that is ready to be populated with data.
//
// Returns:
//
//	Chassis: A new Chassis instance with all elements having empty values.
func NewChassis() Chassis {
	return Chassis{
		Properties:  map[string]any{},
		Resources:   resource_model.ResourceList{},
		CXLSwitches: cxlswitch_model.CXLSwitchList{},
	}
}

// Validate reports whether the receiver is valid.
//
// This method checks the validity of the Chassis instance by verifying that the "id" property exists
// and is a non-empty string. It is a crucial step to ensure that each Chassis instance has a unique identifier
// before proceeding with operations that require a valid Chassis.
//
// Returns:
//
//	bool: True if the Chassis has a valid "id" property, false otherwise.
func (c *Chassis) Validate() bool {
	id, ok := c.Properties["id"].(string)
	if !ok || len(id) <= 0 {
		return false
	}

	return true
}

// ToObject creates and returns a map with elements of id, resources, and any optional elements.
//
// This method first validates the Chassis instance. If the instance is not valid, it returns nil.
// Upon successful validation, it proceeds to construct a map (`res`) initialized with the Chassis's properties.
// It then aggregates resources from both CXLSwitches and Resources, appending them into a single slice.
// This aggregated resources slice is then added to the `res` map under the "resources" key.
// The resulting map, which now includes the Chassis's properties and its aggregated resources, is returned.
//
// Returns:
//
//	map[string]any: A map representation of the Chassis, including its properties and resources, or nil if the Chassis is invalid.
func (c *Chassis) ToObject() map[string]any {
	if !c.Validate() {
		return nil
	}

	res := c.Properties

	resources := c.CXLSwitches.ToObject4Chassis()
	resources = append(resources, c.Resources.ToObject()...)

	res["resources"] = resources

	return res
}
