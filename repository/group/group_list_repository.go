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
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	group_model "github.com/project-cdim/configuration-manager/model/group"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

const getGroupList string = `
MATCH (vrsg: ResourceGroups)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN 
	vrsg, 
	CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`

const getGroupListColumnCount = 6
const (
	getGroupListIndexGroup = iota
	getGroupListIndexResource
	getGroupListIndexAnnotation
	getGroupListIndexResourceGroupIDs
	getGroupListIndexNodeIDs
	getGroupListIndexNotDetected
)

// GroupListRepository is a repository that manages a list of groups.
// The WithResources field indicates whether the repository should include
// resources associated with the groups.
type GroupListRepository struct {
	WithResources bool
}

// NewGroupListRepository creates a new instance of GroupListRepository.
// The withResources parameter specifies whether the repository should include resources.
//
// Parameters:
//
//	withResources - a boolean indicating if resources should be included.
//
// Returns:
//
//	A new instance of GroupListRepository.
func NewGroupListRepository(withResources bool) GroupListRepository {
	return GroupListRepository{
		WithResources: withResources,
	}
}

// FindList retrieves a list of groups from the database based on the provided filter.
// It executes a Cypher query to fetch the group data and processes the results into a structured format.
//
// Parameters:
//   - cmdb: An instance of the CmDb database connection.
//   - filter: A CmFilter instance used to filter the groups.
//
// Returns:
//   - A slice of maps containing group data if successful.
//   - An error if there is an issue executing the query or processing the results.
//
// The function performs the following steps:
//  1. Executes a Cypher query to fetch group data.
//  2. Iterates through the query results and processes each row.
//  3. Sorts the records based on a custom comparison function.
//  4. Constructs group and resource objects from the processed data.
//  5. Applies the provided filter to the groups.
//  6. Returns the filtered groups, optionally including resources based on the repository configuration.
func (glr *GroupListRepository) FindList(cmdb database.CmDb, filter filter.CmFilter) ([]map[string]any, error) {
	common.Log.Debug(getGroupList)
	cypherCursor, err := cmdb.CmDbExecCypher(getGroupListColumnCount, getGroupList)
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
		return compareByGroupList(records, i, j)
	})

	groups := group_model.NewGroupList()
	group := group_model.NewGroup()
	resources := resource_model.NewResourceList()
	preGroupID := ""
	for _, row := range records {
		groupProps := row[getGroupListIndexGroup].(*age.Vertex).Props()
		groupID := groupProps["id"].(string)

		if len(preGroupID) > 0 && preGroupID != groupID {
			group.Resources = resources
			if filter.FilterByCondition(group.ToObject()) {
				groups.Groups = append(groups.Groups, group)
			}
			group = group_model.NewGroup()
			resources = resource_model.NewResourceList()
		}

		group.Id = groupID
		group.Properties = groupProps
		group.CreatedAt = groupProps["createdAt"].(string)
		group.UpdatedAt = groupProps["updatedAt"].(string)

		resource := resource_repository.ComposeResource(
			row[getGroupListIndexResource].(*age.Vertex),
			row[getGroupListIndexAnnotation].(*age.Vertex),
			row[getGroupListIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getGroupListIndexNodeIDs].(*age.SimpleEntity),
			row[getGroupListIndexNotDetected].(*age.SimpleEntity).AsBool(),
			true,
		)

		resources.Resources = append(resources.Resources, resource)
		preGroupID = groupID
	}

	group.Resources = resources
	if filter.FilterByCondition(group.ToObject()) {
		groups.Groups = append(groups.Groups, group)
	}

	if !glr.WithResources {
		return groups.ToObject(), nil
	}
	return groups.ToObjectWithResources(), nil
}

// compareByGroupList compares two records from a list of age.Entity slices based on their group and resource indices.
// It uses the compareByGroup function with specific index retrieval functions for group and resource.
// Parameters:
// - records: A 2D slice of age.Entity representing the records to be compared.
// - i: The index of the first record to compare.
// - j: The index of the second record to compare.
// Returns:
// - A boolean indicating whether the record at index i should sort before the record at index j.
func compareByGroupList(records [][]age.Entity, i, j int) bool {
	return compareByGroup(records, getGroupListIndexGroup, getGroupListIndexResource, i, j)
}
