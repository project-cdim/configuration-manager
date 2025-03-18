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
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
)

// ChassisList is a list of Chassis.
type ChassisList struct {
	Chassis []Chassis
}

// NewChassisList creates and returns a new instance of ChassisList.
//
// This function initializes a ChassisList struct with an empty slice of Chassis.
// It is useful for creating a ChassisList ready to be populated with Chassis instances.
//
// Returns:
//
//	ChassisList: A new instance of ChassisList with an empty slice of Chassis.
func NewChassisList() ChassisList {
	return ChassisList{
		Chassis: []Chassis{},
	}
}

// ToObject creates and returns a map array with elements of
// id, resources, and any optional elements from a node list.
//
// This method iterates over each Chassis in the ChassisList, validates it,
// and then converts it into a map object. Each map object represents a Chassis
// and includes keys for id, resources, and any other optional elements defined
// within the Chassis. These map objects are then collected into a slice and returned.
// This is useful for converting a list of Chassis objects into a more generic data
// structure that can be easily manipulated or serialized.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated Chassis.
func (cl *ChassisList) ToObject() []map[string]any {
	res := []map[string]any{}
	for _, chassis := range cl.Chassis {
		if chassis.Validate() {
			res = append(res, chassis.ToObject())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. chassis(%v)", chassis))
		}
	}

	return res
}
