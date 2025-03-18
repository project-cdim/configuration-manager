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
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

// Node is a node structure.
type Node struct {
	Properties map[string]any
	Resources  resource_model.ResourceList
}

// NewNode is the constructor for the Node structure.
//
// This function initializes a Node struct with all elements having empty values.
// It is useful for creating a Node instance ready to be populated with properties and resources.
//
// Returns:
//
//	Node: A new instance of Node with empty properties and resources.
func NewNode() Node {
	return Node{
		Properties: map[string]any{},
		Resources:  resource_model.ResourceList{},
	}
}

// Validate reports whether the receiver Node is valid.
//
// This method checks the validity of the Node instance by verifying that the "id" property exists
// and is a non-empty string. It is a crucial step to ensure that each Node instance has a unique identifier
// before proceeding with operations that require a valid Node.
//
// Returns:
//
//	bool: True if the Node has a valid "id" property, false otherwise.
func (n *Node) Validate() bool {
	id, ok := n.Properties["id"].(string)
	if !ok || len(id) <= 0 {
		return false
	}

	return true
}

// ToObject creates and returns a map with elements of id, resources, and any optional elements.
//
// This method first validates the Node instance. If the instance is not valid, it returns nil.
// Upon successful validation, it proceeds to construct a map (`res`) initialized with the Node's properties.
// It then adds the resources, converted to their object form for Node, under the "resources" key in the map.
// The resulting map, which now includes the Node's properties and its resources, is returned.
//
// Returns:
//
//	map[string]any: A map representation of the Node, including its properties and resources, or nil if the Node is invalid.
func (n *Node) ToObject() map[string]any {
	if !n.Validate() {
		return nil
	}

	res := n.Properties
	res["resources"] = n.Resources.ToObject4Node()

	return res
}
