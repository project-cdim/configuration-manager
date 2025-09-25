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
	"fmt"
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	node_model "github.com/project-cdim/configuration-manager/model/node"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

// getNode is cypher query to retrieve a specific node.
const getNode string = `
	MATCH (vnd:Node {id: '%s'})
	OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
	OPTIONAL MATCH (vrs)-[:Have]->(van)
	OPTIONAL MATCH (vrs)-[endt:NotDetected]->(:NotDetectedDevice)
	OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
	WITH vnd, vrs, van, vrsg, endt
	RETURN
		vnd,
		CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END,
		CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
		COLLECT(vrsg.id),
		CASE WHEN endt IS NULL THEN true ELSE false END
`
const getNodeColumnCount = 5
const (
	getNodeIndexNode = iota
	getNodeIndexResource
	getNodeIndexAnnotation
	getNodeIndexResourceGroupIDs
	getNodeIndexNotDetected
)

// NodeRepository is a repository structure for getting a specific node.
type NodeRepository struct {
	NodeID string
}

// NewNodeRepository creates and returns a NodeRepository object that holds the argument nodeID.
// This function initializes a NodeRepository with a specific nodeID.
func NewNodeRepository(nodeID string) NodeRepository {
	return NodeRepository{
		NodeID: nodeID,
	}
}

// Find retrieves a node based on the provided CmDb and CmFilter.
// It constructs and executes a Cypher query to fetch node information from the database.
// The function processes the query results, assembling them into a structured map representing the node and its components.
// If successful, it returns the assembled node as a map[string]any, or an error if the operation fails.
func (nr *NodeRepository) Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error) {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", getNode, nr.NodeID))
	cypherCursor, err := cmdb.CmDbExecCypher(getNodeColumnCount, getNode, nr.NodeID)
	if err != nil {
		return nil, err
	}
	defer cypherCursor.Close()

	records := [][]age.Entity{}
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}
		records = append(records, row)
	}

	sort.Slice(records, func(i, j int) bool {
		return compareByNodeSingle(records, i, j)
	})

	nodeWork := node_model.NewNode()
	for _, row := range records {
		// Retrieve the properties of the Node and store them in the model
		nodeWork.Properties = row[getNodeIndexNode].(*age.Vertex).Props()

		// From one record of the search result, assemble the Resource information associated with the Node, and add it to the list of Resources
		// *If the assembled Resource information is empty, it is determined that there are no Resources associated with the Node, and it is not added to the list of Resources
		resource := resource_repository.ComposeResource(
			row[getNodeIndexResource].(*age.Vertex),
			row[getNodeIndexAnnotation].(*age.Vertex),
			row[getNodeIndexResourceGroupIDs].(*age.SimpleEntity),
			age.NewSimpleEntity([]any{}),
			row[getNodeIndexNotDetected].(*age.SimpleEntity).AsBool(),
			true,
		)
		if resource.Validate() {
			nodeWork.Resources.Resources = append(nodeWork.Resources.Resources, resource)
		}
	}

	node := node_model.NewNode()
	if filter.FilterByCondition(nodeWork.ToObject()) {
		node = nodeWork
	}

	return node.ToObject(), nil
}

// compareByNodeSingle is a wrapper function to sort the contents of a Cypher query execution result.
// It leverages the compareByNode function to determine the order of records based on node and resource information.
func compareByNodeSingle(records [][]age.Entity, i, j int) bool {
	return compareByNode(records, getNodeIndexNode, getNodeIndexResource, i, j)
}
