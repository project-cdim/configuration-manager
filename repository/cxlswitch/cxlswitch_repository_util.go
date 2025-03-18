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
        
package cxlswitch_repository

import (
	"strings"

	"github.com/apache/age/drivers/golang/age"
)

// compareByCXLSwitch sorts the contents of a Cypher query execution result based on the following criteria:
// - First sort key: CXLSwitch ID (string, empty string if absent, ascending order)
// - Second sort key: Device ID of the resources linked to the CXLSwitch (string, empty string if absent, ascending order)
// This function is used to order records by CXLSwitch and associated resources.
func compareByCXLSwitch(records [][]age.Entity, switchIdx, resourceIdx, i, j int) bool {
	row1 := records[i]
	row2 := records[j]

	switchProp1 := row1[switchIdx].(*age.Vertex).Props()
	switchProp2 := row2[switchIdx].(*age.Vertex).Props()

	// Obtain the first sort key for each and compare if they are not equal
	switchID1, ok := switchProp1["id"].(string)
	if !ok {
		switchID1 = ""
	}
	switchID2, ok := switchProp2["id"].(string)
	if !ok {
		switchID2 = ""
	}

	if switchID1 != switchID2 {
		return strings.Compare(switchID1, switchID2) < 0
	}

	resProp1 := row1[resourceIdx].(*age.Vertex).Props()
	resProp2 := row2[resourceIdx].(*age.Vertex).Props()

	// Obtain the second sort key for each and compare
	deviceID1, ok := resProp1["deviceID"].(string)
	if !ok {
		deviceID1 = ""
	}
	deviceID2, ok := resProp2["deviceID"].(string)
	if !ok {
		deviceID2 = ""
	}

	return strings.Compare(deviceID1, deviceID2) < 0
}
