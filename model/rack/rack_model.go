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
        
package rack_model

import (
	chassis_model "github.com/project-cdim/configuration-manager/model/chassis"
)

// Rack is a Rack structure.
type Rack struct {
	Properties map[string]any
	Chassis    chassis_model.ChassisList
}

// NewRack is the constructor for the Rack structure.
//
// This function initializes a Rack struct with all elements having empty values.
// It is useful for creating a Rack instance ready to be populated with properties and chassis.
//
// Returns:
//
//	Rack: A new instance of Rack with empty properties and an empty ChassisList.
func NewRack() Rack {
	return Rack{
		Properties: map[string]any{},
		Chassis:    chassis_model.ChassisList{},
	}
}

// Validate reports whether the receiver Rack is valid.
//
// This method checks the validity of the Rack instance by verifying that the "id" property exists
// and is a non-empty string. It is a crucial step to ensure that each Rack instance has a unique identifier
// before proceeding with operations that require a valid Rack.
//
// Returns:
//
//	bool: True if the Rack has a valid "id" property, false otherwise.
func (r *Rack) Validate() bool {
	id, ok := r.Properties["id"].(string)
	if !ok || len(id) <= 0 {
		return false
	}

	return true
}

// ToObject creates and returns a map with elements of id, resources, and any optional elements.
//
// This method first validates the Rack instance. If the instance is not valid, it returns nil.
// Upon successful validation, it proceeds to construct a map (`res`) initialized with the Rack's properties.
// It then adds the chassis, converted to their object form, under the "chassis" key in the map.
// The resulting map, which now includes the Rack's properties and its chassis, is returned.
//
// Returns:
//
//	map[string]any: A map representation of the Rack, including its properties and chassis, or nil if the Rack is invalid.
func (r *Rack) ToObject() map[string]any {
	if !r.Validate() {
		return nil
	}

	res := r.Properties

	res["chassis"] = r.Chassis.ToObject()

	return res
}
