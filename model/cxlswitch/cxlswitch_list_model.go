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
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
)

// CXLSwitchList is a list of CXLSwitches.
type CXLSwitchList struct {
	CXLSwitches []CXLSwitch
}

// NewCXLSwitchList is a constructor for the CXLSwitchList structure.
//
// This function initializes a CXLSwitchList struct with an empty slice of CXLSwitch.
// It is useful for creating a CXLSwitchList ready to be populated with CXLSwitch instances.
//
// Returns:
//
//	CXLSwitchList: A new instance of CXLSwitchList with an empty slice of CXLSwitch.
func NewCXLSwitchList() CXLSwitchList {
	return CXLSwitchList{
		CXLSwitches: []CXLSwitch{},
	}
}

// ToObject creates and returns a map array with elements of
// id, resources, and any optional elements from a node list.
//
// This function iterates over each CXLSwitch in the CXLSwitchList, validates it,
// and then converts it into a map object. Each map object represents a CXLSwitch
// and includes keys for id, resources, and any other optional elements defined
// within the CXLSwitch. These map objects are then collected into a slice and returned.
// This is useful for converting a list of CXLSwitch objects into a more generic data
// structure that can be easily manipulated or serialized.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated CXLSwitch.
func (nl *CXLSwitchList) ToObject() []map[string]any {
	res := []map[string]any{}
	for _, cxlswitch := range nl.CXLSwitches {
		if cxlswitch.Validate() {
			res = append(res, cxlswitch.ToObject())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. CXLSwitch(%v)", cxlswitch))
		}
	}

	return res
}

// ToObject4Chassis creates for ChassisObject and returns a map array with elements of
// id, resources, and any optional elements from a node list.
//
// This function iterates over each CXLSwitch in the CXLSwitchList, validates it,
// and then converts it into a map object specifically formatted for ChassisObject consumption.
// Each map object represents a CXLSwitch and includes keys for id, resources, and any other
// optional elements defined within the CXLSwitch. These map objects are then collected into
// a slice and returned. This is particularly useful for aggregating CXLSwitch data in a format
// that is compatible with ChassisObject requirements.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated CXLSwitch, formatted for ChassisObject.
func (nl *CXLSwitchList) ToObject4Chassis() []map[string]any {
	res := []map[string]any{}
	for _, cxlswitch := range nl.CXLSwitches {
		if cxlswitch.Validate() {
			res = append(res, cxlswitch.ToObject4Chassis())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. CXLSwitch(%v)", cxlswitch))
		}
	}

	return res
}
