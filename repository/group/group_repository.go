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
	"fmt"
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"

	group_model "github.com/project-cdim/configuration-manager/model/group"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

const getGroup string = `
MATCH (vrsg:ResourceGroups {id: '%s'})
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(:NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN
	vrsg,
	CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`

const getGroupColumnCount = 6
const (
	getGroupIndexGroup = iota
	getGroupIndexResource
	getGroupIndexAnnotation
	getGroupIndexResourceGroupIDs
	getGroupIndexNodeIDs
	getGroupIndexNotDetected
)

// GroupRepository represents a repository for managing groups.
// It contains the GroupID which uniquely identifies the group,
// and a boolean WithResources indicating whether the group includes resources.
type GroupRepository struct {
	GroupID       string
	WithResources bool
}

// NewGroupRepository creates a new instance of GroupRepository with the specified groupID and resource inclusion flag.
// Parameters:
//   - groupID: A string representing the unique identifier of the group.
//   - withResources: A boolean flag indicating whether to include resources in the repository.
//
// Returns:
//
//	A new instance of GroupRepository.
func NewGroupRepository(groupID string, withResources bool) GroupRepository {
	return GroupRepository{
		GroupID:       groupID,
		WithResources: withResources,
	}
}

// Find retrieves a group and its associated resources from the database based on the provided filter.
// It executes a Cypher query to fetch the group and resource data, processes the results, and returns
// the group data in a map format.
//
// Parameters:
//   - cmdb: An instance of the CmDb database connection.
//   - filter: A CmFilter instance to filter the results.
//
// Returns:
//   - A map containing the group data and its resources if the filter conditions are met.
//   - An error if any issues occur during the database query or data processing.
func (gr *GroupRepository) Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error) {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", getGroup, gr.GroupID))
	cypherCursor, err := cmdb.CmDbExecCypher(getGroupColumnCount, getGroup, gr.GroupID)
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
		return compareByGroupSingle(records, i, j)
	})

	group := group_model.NewGroup()
	resources := resource_model.NewResourceList()

	for i, row := range records {
		if i == 0 {
			groupProps := row[getGroupListIndexGroup].(*age.Vertex).Props()
			group.Id = groupProps["id"].(string)
			group.Properties = groupProps
			group.CreatedAt = groupProps["createdAt"].(string)
			group.UpdatedAt = groupProps["updatedAt"].(string)
		}
		resourceWork := resource_repository.ComposeResource(
			row[getGroupIndexResource].(*age.Vertex),
			row[getGroupIndexAnnotation].(*age.Vertex),
			row[getGroupIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getGroupIndexNodeIDs].(*age.SimpleEntity),
			row[getGroupIndexNotDetected].(*age.SimpleEntity).AsBool(),
			true,
		)

		resources.Resources = append(resources.Resources, resourceWork)
	}

	group.Resources = resources

	res := group_model.NewGroup()
	if filter.FilterByCondition(group.ToObjectWithResources()) {
		res = group
	}

	if !gr.WithResources {
		return res.ToObject(), nil
	}
	return res.ToObjectWithResources(), nil
}

// compareByGroupSingle compares two records based on their group and resource indices.
// It uses the compareByGroup function with specific index retrieval functions for groups and resources.
//
// Parameters:
//
//	records - A 2D slice of age.Entity representing the records to be compared.
//	i - The index of the first record to compare.
//	j - The index of the second record to compare.
//
// Returns:
//
//	A boolean value indicating whether the record at index i should be sorted before the record at index j.
func compareByGroupSingle(records [][]age.Entity, i, j int) bool {
	return compareByGroup(records, getGroupIndexGroup, getGroupIndexResource, i, j)
}
