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
        
package node_repository

import (
	"strings"

	"github.com/apache/age/drivers/golang/age"
)

// compareByNode sorts the contents of a Cypher query execution result based on the following criteria:
// - First sort key: Node ID (string, empty string if absent, ascending order)
// - Second sort key: Device ID of the resources linked to the node (string, empty string if absent, ascending order)
// This function is used to order records by node and associated resources.
func compareByNode(records [][]age.Entity, nodeIdx, resourceIdx, i, j int) bool {
	row1 := records[i]
	row2 := records[j]

	nodeProp1 := row1[nodeIdx].(*age.Vertex).Props()
	nodeProp2 := row2[nodeIdx].(*age.Vertex).Props()

	// Obtain the first sort key for each and compare if they are not equal
	nodeID1, ok := nodeProp1["id"].(string)
	if !ok {
		nodeID1 = ""
	}
	nodeID2, ok := nodeProp2["id"].(string)
	if !ok {
		nodeID2 = ""
	}

	if nodeID1 != nodeID2 {
		return strings.Compare(nodeID1, nodeID2) < 0
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
