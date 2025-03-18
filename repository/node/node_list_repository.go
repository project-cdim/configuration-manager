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
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	node_model "github.com/project-cdim/configuration-manager/model/node"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

// getNodeList is cypher query to get a list of nodes.
const getNodeList string = `
	MATCH (vnd: Node)
	OPTIONAL MATCH (vnd)-[ecm:Compose]->(vrs)
	OPTIONAL MATCH (vrs)-[ehv:Have]->(van)
	OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
	OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
	WITH vnd, vrs, van, vrsg, endt 
	RETURN
		vnd, 
		CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END, 
		CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
		COLLECT(vrsg.id), 
		CASE WHEN endt IS NULL THEN true ELSE false END 
`
const getNodeListColumnCount = 5
const (
	getNodeListIndexNode = iota
	getNodeListIndexResource
	getNodeListIndexAnnotation
	getNodeListIndexResourceGroupIDs
	getNodeListIndexNotDetected
)

// NodeListRepository is a repository structure for getting node lists.
type NodeListRepository struct{}

// NewNodeListRepository creates and returns a NodeListRepository object.
// This constructor function initializes a new instance of NodeListRepository with default values.
func NewNodeListRepository() NodeListRepository {
	return NodeListRepository{}
}

// FindList retrieves a list of nodes from the database based on the provided CmDb and CmFilter.
// It executes a Cypher query to fetch node information and processes the results into a structured list.
// Each node in the list includes its properties and associated resources.
// The function returns a slice of maps, each representing a node and its resources, or an error if the operation fails.
func (nlr *NodeListRepository) FindList(cmdb database.CmDb, filter filter.CmFilter) ([]map[string]any, error) {
	query := getNodeList

	common.Log.Debug(query)
	cypherCursor, err := cmdb.CmDbExecCypher(getNodeListColumnCount, query)
	if err != nil {
		return nil, err
	}

	records := [][]age.Entity{}
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}
		records = append(records, row)
	}
	cypherCursor.Close()

	sort.Slice(records, func(i, j int) bool {
		return compareByNodeList(records, i, j)
	})

	nodeList := node_model.NewNodeList()
	node := node_model.NewNode()
	resources := resource_model.NewResourceList()
	preNodeID := ""
	for _, row := range records {
		// Retrieve the property information from the Vertex data (Properties can be obtained in map format)
		nodeWork := row[getNodeListIndexNode].(*age.Vertex).Props()
		nodeID := nodeWork["id"].(string)

		if len(preNodeID) > 0 && preNodeID != nodeID {
			// If the Node information of the previous row read is different, append the accumulated Node information to the nodes array and clear the Node information
			node.Resources = resources
			if filter.FilterByCondition(node.ToObject()) {
				nodeList.Nodes = append(nodeList.Nodes, node)
			}
			node = node_model.NewNode()
			resources = resource_model.NewResourceList()
		}

		node.Properties = nodeWork

		// From one record of the search result, assemble the Resource information associated with the Node, and add it to the list of Resources
		// *If the assembled Resource information is empty, it is determined that there are no Resources associated with the Node, and it is not added to the list of Resources
		resource := resource_repository.ComposeResource(
			row[getNodeListIndexResource].(*age.Vertex),
			row[getNodeListIndexAnnotation].(*age.Vertex),
			row[getNodeListIndexResourceGroupIDs].(*age.SimpleEntity),
			age.NewSimpleEntity([]any{}),
			row[getNodeListIndexNotDetected].(*age.SimpleEntity).AsBool(),
			false,
		)
		if resource.Validate() {
			resources.Resources = append(resources.Resources, resource)
		}

		// Update the nodeID being read
		preNodeID = nodeID
	}

	// Append the last accumulated Node information to the nodes array
	node.Resources = resources
	if filter.FilterByCondition(node.ToObject()) {
		nodeList.Nodes = append(nodeList.Nodes, node)
	}

	return nodeList.ToObject(), nil
}

// compareByNodeList is a wrapper function to sort the contents of a Cypher query execution result.
// It leverages the compareByNode function to determine the order of records based on node and resource information.
func compareByNodeList(records [][]age.Entity, i, j int) bool {
	return compareByNode(records, getNodeListIndexNode, getNodeListIndexResource, i, j)
}
