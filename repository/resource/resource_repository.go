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
MATCH (vrs: %s{deviceID: '%s'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`

// Due to the relationships of the data registered in the DB, both UNION and UNION ALL return the same data. Therefore, considering the search speed efficiency, UNION ALL is used.
const queryResource_unionall string = `
UNION ALL`

// getQueryResource constructs a Cypher query to retrieve a specific resource by deviceID.
// It iterates through a list of resource types, creating a part of the query for each,
// and then joins these parts with a UNION ALL clause to form the final query.
func getQueryResource(deviceID string) string {
	items := []string{}
	for _, resourceType := range resourceTypeList {
		items = append(items, fmt.Sprintf(queryResource_match_return, resourceType, deviceID))
	}
	return strings.Join(items, queryResource_unionall)
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
	query := getQueryResource(rr.DeviceID)
	common.Log.Debug(query)
	cypherCursor, err := cmdb.CmDbExecCypher(getResourceColumnCount, query)
	if err != nil {
		return nil, err
	}

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
	cypherCursor.Close()

	return resource.ToObject(), nil
}
