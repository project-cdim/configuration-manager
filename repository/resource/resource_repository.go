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
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"

	resource_model "github.com/project-cdim/configuration-manager/model/resource"

	"github.com/apache/age/drivers/golang/age"
)

// Cypher query fragment to retrieve a specific resource
const queryResource_match_return string = `
MATCH (vrs:%s{deviceID: '%s'})
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(:NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`

// Due to the relationships of the data registered in the DB, both UNION and UNION ALL return the same data. Therefore, considering the search speed efficiency, UNION ALL is used.
const queryResource_unionall string = `
UNION ALL`

// getQueryResource generates a SQL query string by combining individual resource type queries.
// It iterates through the resourceTypeList, creating a query for each resource type based on
// queryResource_match_return, and then joins these queries together using queryResource_unionall.
// The resulting string represents a union of queries for all resource types.
func getQueryResource() string {
	items := []string{}
	for range ResourceTypeList {
		items = append(items, queryResource_match_return)
	}
	return strings.Join(items, queryResource_unionall)
}

// getQueryResourceParam generates a slice of any containing the parameters
// needed for querying resources. It iterates through the resourceTypeList,
// appending each resource type and the deviceID to the slice. This slice is
// then used as arguments in a database query.
func getQueryResourceParam(deviceID string) []any {
	items := []any{}
	for _, resourceType := range ResourceTypeList {
		items = append(items, resourceType)
		items = append(items, deviceID)
	}
	return items
}

const getResourceColumnCount = 5
const (
	getResourceIndexResource = iota
	getResourceIndexAnnotation
	getResourceIndexResourceGroupIDs
	getResourceIndexNodeIDs
	getResourceIndexNotDetected
)

// ResourceListRepository is a repository structure for getting a specific resource.
type ResourceRepository struct {
	DeviceID string
}

// NewResourceRepository creates and returns a ResourceRepository object that holds the argument deviceId.
func NewResourceRepository(deviceID string) ResourceRepository {
	return ResourceRepository{
		DeviceID: deviceID,
	}
}

// Find returns a resource. It constructs a query using the deviceID, executes it, and processes the results.
// If a matching resource is found, it is returned; otherwise, an error is returned.
func (rr *ResourceRepository) Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error) {
	query := getQueryResource()
	queryParam := getQueryResourceParam(rr.DeviceID)
	common.Log.Debug(fmt.Sprintf("query: %s", query))
	cypherCursor, err := cmdb.CmDbExecCypher(getResourceColumnCount, query, queryParam...)
	if err != nil {
		return nil, err
	}
	defer cypherCursor.Close()

	resource := resource_model.NewResource()
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}

		// Assemble Resource information from a single record of the search results
		resourceWork := ComposeResource(
			row[getResourceIndexResource].(*age.Vertex),
			row[getResourceIndexAnnotation].(*age.Vertex),
			row[getResourceIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getResourceIndexNodeIDs].(*age.SimpleEntity),
			row[getResourceIndexNotDetected].(*age.SimpleEntity).AsBool(),
			true,
		)
		if filter.FilterByCondition(resourceWork.ToObject()) {
			resource = resourceWork
		}
		break
	}

	return resource.ToObject(), nil
}
