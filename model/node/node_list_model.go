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
        
package node_model

import (
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
)

// NodeList is a list of nodes.
type NodeList struct {
	Nodes []Node
}

// NewNodeList is the constructor for the NodeList structure.
//
// This function initializes a NodeList struct with an empty slice of Node.
// It is useful for creating a NodeList ready to be populated with Node instances.
//
// Returns:
//
//	NodeList: A new instance of NodeList with an empty slice of Node.
func NewNodeList() NodeList {
	return NodeList{
		Nodes: []Node{},
	}
}

// ToObject creates and returns a map array with elements of
// id, resources, and any optional elements from a node list.
//
// This function iterates over each Node in the NodeList, validates it,
// and then converts it into a map object. Each map object represents a Node
// and includes keys for id, resources, and any other optional elements defined
// within the Node. These map objects are then collected into a slice and returned.
// This is useful for converting a list of Node objects into a more generic data
// structure that can be easily manipulated or serialized.
//
// Returns:
//
//	[]map[string]any: A slice of map objects, each representing a validated Node.
func (nl *NodeList) ToObject() []map[string]any {
	res := []map[string]any{}
	for _, node := range nl.Nodes {
		if node.Validate() {
			res = append(res, node.ToObject())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. node(%v)", node))
		}
	}

	return res
}
