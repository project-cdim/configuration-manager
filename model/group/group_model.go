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
        
package group_model

import (
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/model"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

// Group represents a group entity with its properties and associated resources.
//
// Fields:
// - Id: A unique identifier for the group.
// - Properties: A map containing various properties of the group.
// - CreatedAt: The timestamp when the group was created.
// - UpdatedAt: The timestamp when the group was last updated.
// - Resources: A list of resources associated with the group.
type Group struct {
	Id         string
	Properties map[string]any
	CreatedAt  string
	UpdatedAt  string
	Resources  resource_model.ResourceList
}

// NewGroup creates and returns a new Group instance with default values.
// The Id is initialized as an empty string, Properties is an empty map,
// CreatedAt and UpdatedAt are empty strings, and Resources is initialized
// using the NewResourceList function from the resource_model package.
func NewGroup() Group {
	return Group{
		Id:         "",
		Properties: map[string]any{},
		CreatedAt:  "",
		UpdatedAt:  "",
		Resources:  resource_model.NewResourceList(),
	}
}

// NewGroupWithCreateTimeStampsNow creates a new Group instance, sets its creation timestamps to the current time,
// and assigns the provided properties to the Group.
//
// Parameters:
//   - properties: A map containing the properties to be assigned to the new Group.
//
// Returns:
//
//	A new Group instance with the specified properties and current creation timestamps.
func NewGroupWithCreateTimeStampsNow(properties map[string]any) Group {
	g := NewGroup()
	g.createTimeStampsNow()
	g.Properties = properties
	return g
}

// NewGroupForUpdate creates a new Group instance for updating purposes.
// It takes the existing group data from the database and a map of properties to be updated.
// The function initializes a new Group, sets its ID and creation timestamp from the database,
// updates the timestamps to the current time, and assigns the provided properties.
//
// Parameters:
//   - groupFromDb: A map containing the existing group data from the database.
//   - properties: A map containing the properties to be updated.
//
// Returns:
//   - Group: A new Group instance with updated properties and timestamps.
func NewGroupForUpdate(groupFromDb map[string]any, properties map[string]any) Group {
	g := NewGroup()
	g.Id = groupFromDb["id"].(string)
	g.CreatedAt = groupFromDb["createdAt"].(string)
	g.updateTimeStampsNow()
	g.Properties = properties
	return g
}

// createTimeStampsNow sets the CreatedAt and UpdatedAt fields of the Group to the current time
// in ISO8601 format.
func (g *Group) createTimeStampsNow() {
	now := model.CurrentTimeISO8601()
	g.CreatedAt = now
	g.UpdatedAt = now
}

// updateTimeStampsNow updates the UpdatedAt field of the Group struct to the current time
// in ISO 8601 format.
func (g *Group) updateTimeStampsNow() {
	now := model.CurrentTimeISO8601()
	g.UpdatedAt = now
}

// Validate checks the properties and timestamps of the Group object.
// It returns true if all validations pass, otherwise it returns false.
// The function performs the following checks:
// 1. Validates the properties of the Group using ValidateProperty function.
// 2. Validates the CreatedAt timestamp to ensure it is in ISO8601 format.
// 3. Validates the UpdatedAt timestamp to ensure it is in ISO8601 format.
// If any of these validations fail, the function returns false and prints an error message.
func (g *Group) Validate() bool {
	if !ValidateProperty(g.Properties) {
		return false
	}

	if !model.ValidateISO8601(g.CreatedAt) {
		common.Log.Warn(fmt.Sprintf("createdAt is not ISO8601. createdAt(%v)", g.CreatedAt))
		return false
	}

	if !model.ValidateISO8601(g.UpdatedAt) {
		common.Log.Warn(fmt.Sprintf("updatedAt is not ISO8601. updatedAt(%v)", g.UpdatedAt))
		return false
	}

	return true
}

// ToObject converts the Group struct to a map representation.
// It returns nil if the Group is not valid.
// The returned map contains the following keys:
// - "id": the ID of the group
// - "name": the name of the group
// - "description": the description of the group
// - "createdAt": the creation timestamp of the group
// - "updatedAt": the last update timestamp of the group
func (g *Group) ToObject() map[string]any {
	if !g.Validate() {
		return nil
	}

	return map[string]any{
		"id":          g.Id,
		"name":        g.Properties["name"].(string),
		"description": g.Properties["description"].(string),
		"createdAt":   g.CreatedAt,
		"updatedAt":   g.UpdatedAt,
	}
}

// ToObjectWithResources converts the Group object into a map representation
// including its resources. It returns nil if the Group object is not valid.
// The returned map contains the following keys:
// - "id": the ID of the group
// - "name": the name of the group
// - "description": the description of the group
// - "createdAt": the creation timestamp of the group
// - "updatedAt": the last update timestamp of the group
// - "resources": the resources associated with the group, converted to an object
func (g *Group) ToObjectWithResources() map[string]any {
	if !g.Validate() {
		return nil
	}

	return map[string]any{
		"id":          g.Id,
		"name":        g.Properties["name"].(string),
		"description": g.Properties["description"].(string),
		"createdAt":   g.CreatedAt,
		"updatedAt":   g.UpdatedAt,
		"resources":   g.Resources.ToObject(),
	}
}
