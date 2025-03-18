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
        
package rack_repository

import (
	"strings"

	"github.com/apache/age/drivers/golang/age"
)

// compareByRack sorts the contents of a Cypher query execution result based on the following criteria:
// - First sort key: Rack > Chassis' unitPosition (numeric, -1 if absent, ascending order)
// - Second sort key: Rack > Chassis > Device's type (string, empty string if absent, ascending order)
// - Third sort key: Rack > Chassis > Device's deviceID (string, empty string if absent, ascending order)
// - Fourth sort key: Rack > Chassis > Device's id (string, empty string if absent, ascending order)
// This function is used to order records by chassis and device information within a rack.
func compareByRack(records [][]age.Entity, chassisIdx, deviceIdx, i, j int) bool {
	row1 := records[i]
	row2 := records[j]

	chassisProp1 := row1[chassisIdx].(*age.Vertex).Props()
	chassisProp2 := row2[chassisIdx].(*age.Vertex).Props()

	// Obtain the first sort key for each and compare if they are not equal
	position1, ok := chassisProp1["unitPosition"].(int64)
	if !ok {
		position1 = int64(-1)
	}
	position2, ok := chassisProp2["unitPosition"].(int64)
	if !ok {
		position2 = int64(-1)
	}

	if position1 != position2 {
		return position1 < position2
	}

	devProp1 := row1[deviceIdx].(*age.Vertex).Props()
	devProp2 := row2[deviceIdx].(*age.Vertex).Props()

	// Obtain the second sort key for each and compare if they are not equal
	type1, ok := devProp1["type"].(string)
	if !ok {
		type1 = ""
	}
	type2, ok := devProp2["type"].(string)
	if !ok {
		type2 = ""
	}

	if type1 != type2 {
		return strings.Compare(type1, type2) < 0
	}

	// Obtain the third sort key for each and compare if they are not equal
	deviceID1, ok := devProp1["deviceID"].(string)
	if !ok {
		deviceID1 = ""
	}
	deviceID2, ok := devProp2["deviceID"].(string)
	if !ok {
		deviceID2 = ""
	}

	if deviceID1 != deviceID2 {
		return strings.Compare(deviceID1, deviceID2) < 0
	}

	// Obtain the fourth sort key for each and compare
	id1, ok := devProp1["id"].(string)
	if !ok {
		id1 = ""
	}
	id2, ok := devProp2["id"].(string)
	if !ok {
		id2 = ""
	}

	return strings.Compare(id1, id2) < 0
}
