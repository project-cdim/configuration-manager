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

package resource_repository

import (
	"fmt"
	"sort"
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	resource_filter "github.com/project-cdim/configuration-manager/filter/resource"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"

	"github.com/apache/age/drivers/golang/age"
)

// ResourceTypeList is a list of resource types.
// Note: The type has been changed from [...]string to []any to allow expanding
// this list as variadic arguments when needed.
var ResourceTypeList = []any{
	"CPU",
	"Accelerator",
	"DSP",
	"FPGA",
	"GPU",
	"UnknownProcessor",
	"Memory",
	"Storage",
	"NetworkInterface",
	"GraphicController",
	"VirtualMedia",
}

// Cypher query fragment to retrieve a specific resource
const queryResourceList_match_return string = `
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`

// Due to the relationships of the data registered in the DB, both UNION and UNION ALL return the same data. Therefore, considering the search speed efficiency, UNION ALL is used.
const queryResourceList_unionall string = `
UNION ALL`

// getQueryResourceList generates a SQL query string by combining multiple resource type queries.
// It iterates over the resourceTypeList, appending a predefined query pattern for each resource type.
// The resulting queries are then joined together using a UNION ALL clause to form a single query.
func getQueryResourceList() string {
	items := []string{}
	for range ResourceTypeList {
		items = append(items, queryResourceList_match_return)
	}
	return strings.Join(items, queryResourceList_unionall)
}

const getResourceListColumnCount = 5
const (
	getResourceListIndexResource = iota
	getResourceListIndexAnnotation
	getResourceListIndexResourceGroupIDs
	getResourceListIndexNodeIDs
	getResourceListIndexNotDetected
)

// ResourceListRepository is a repository structure for getting resource lists.
type ResourceListRepository struct {
	Detail bool
}

// NewResourceListRepository creates and returns a ResourceListRepository object that holds the arguments detail.
// This constructor function initializes a ResourceListRepository with a detail flag indicating whether to retrieve detailed resource information.
func NewResourceListRepository(detail bool) ResourceListRepository {
	return ResourceListRepository{
		Detail: detail,
	}
}

// FindList retrieves a list of resources that match the given filter conditions.
// It constructs a query to fetch resources, executes it, and processes the results.
// Each resource is composed into a map[string]any format, and if it passes the filter conditions, it's added to the result list.
// The function returns a slice of map[string]any representing the resources, or an error if the operation fails.
func (rlr *ResourceListRepository) FindList(cmdb database.CmDb, filter filter.CmFilter) ([]map[string]any, error) {
	query := getQueryResourceList()
	common.Log.Debug(fmt.Sprintf("query: %s, params: %v", query, ResourceTypeList))
	cypherCursor, err := cmdb.CmDbExecCypher(getResourceListColumnCount, query, ResourceTypeList...)
	if err != nil {
		return nil, err
	}
	defer cypherCursor.Close()

	resourceList := resource_model.NewResourceList()
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}

		// Assemble Resource information from a single record of the search results
		resource := ComposeResource(
			row[getResourceListIndexResource].(*age.Vertex),
			row[getResourceListIndexAnnotation].(*age.Vertex),
			row[getResourceListIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getResourceListIndexNodeIDs].(*age.SimpleEntity),
			row[getResourceListIndexNotDetected].(*age.SimpleEntity).AsBool(),
			rlr.Detail,
		)
		if filter.FilterByCondition(resource.ToObject()) {
			// Append a single record of search results to the variable resources (information of search results)
			resourceList.Resources = append(resourceList.Resources, resource)
		}
	}

	switch filter.(type) {
	case resource_filter.ResourceAvailableFilter:
		// sort by resourceType and deviceID
		sortResourceList(resourceList.Resources)
		return resourceList.ToObject(), nil
	case resource_filter.ResourceUnusedFilter:
		// sort by resourceType and deviceID
		sortResourceList(resourceList.Resources)
		return resourceList.ToObject4Unused(), nil
	default:
		// sort by deviceID
		sort.Slice(resourceList.Resources, func(i, j int) bool {
			return strings.Compare(resourceList.Resources[j].Device["deviceID"].(string), resourceList.Resources[i].Device["deviceID"].(string)) > 0
		})
		return resourceList.ToObject(), nil
	}
}
