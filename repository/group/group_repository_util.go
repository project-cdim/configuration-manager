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
        
package group_repository

import (
	"strings"

	"github.com/project-cdim/configuration-manager/common"

	"github.com/apache/age/drivers/golang/age"
)

// compareByGroup compares two records based on their group and resource properties.
// It first compares the group IDs of the records. If the group IDs are different,
// it prioritizes the default group ID defined in common.DefaultGroupId. If neither
// group ID is the default, it compares the group IDs lexicographically.
// If the group IDs are the same, it compares the device IDs of the resources
// lexicographically.
//
// Parameters:
// - records: A 2D slice of age.Entity representing the records to be compared.
// - groupIdx: The index of the group property in the records.
// - resourceIdx: The index of the resource property in the records.
// - i: The index of the first record to compare.
// - j: The index of the second record to compare.
//
// Returns:
// - A boolean value indicating whether the first record should be ordered before the second record.
func compareByGroup(records [][]age.Entity, groupIdx, resourceIdx, i, j int) bool {
	row1 := records[i]
	row2 := records[j]

	groupProp1 := row1[groupIdx].(*age.Vertex).Props()
	groupProp2 := row2[groupIdx].(*age.Vertex).Props()

	groupID1, ok := groupProp1["id"].(string)
	if !ok {
		groupID1 = ""
	}
	groupID2, ok := groupProp2["id"].(string)
	if !ok {
		groupID2 = ""
	}

	if groupID1 != groupID2 {
		if groupID1 == common.DefaultGroupId {
			return true
		} else if groupID2 == common.DefaultGroupId {
			return false
		}
		return strings.Compare(groupID1, groupID2) < 0
	}

	resProp1 := row1[resourceIdx].(*age.Vertex).Props()
	resProp2 := row2[resourceIdx].(*age.Vertex).Props()

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
