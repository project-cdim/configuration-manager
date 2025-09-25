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

package cxlswitch_repository

import (
	"fmt"
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

// Cypher query to get a specific CXLSwitch
// Retrieves the Vertex of the CXLSwitch and the resources and related Vertices associated with the CXLSwitch
const getCXLSwitch string = `
	MATCH (vcx:CXLswitch {id: '%s'})
	OPTIONAL MATCH (vcx)-[:Connect]->(vrs)
	OPTIONAL MATCH (vrs)-[:Have]->(van)
	OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
	OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
	OPTIONAL MATCH (vrs)-[endt:NotDetected]->(:NotDetectedDevice)
	WITH vcx, vrs, van, vrsg, vnd, endt
	RETURN
		vcx,
		CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END,
		CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
		COLLECT(vrsg.id),
		COLLECT(vnd.id),
		CASE WHEN endt IS NULL THEN true ELSE false END
`
const getCXLSwitchColumnCount = 6
const (
	getCXLSwitchIndexCXLSwitch = iota
	getCXLSwitchIndexResource
	getCXLSwitchIndexAnnotation
	getCXLSwitchIndexResourceGroupIDs
	getCXLSwitchIndexNodeIDs
	getCXLSwitchIndexNotDetected
)

// CXLSwitchRepository is a repository structure for getting a specific cxlSwitch.
type CXLSwitchRepository struct {
	CXLSwitchID string
}

// NewCXLSwitchRepository creates and returns a CXLSwitchRepository object that holds the argument cxlSwitchID.
// This constructor function initializes a new instance of CXLSwitchRepository with the specified CXLSwitchID.
func NewCXLSwitchRepository(cxlSwitchID string) CXLSwitchRepository {
	return CXLSwitchRepository{
		CXLSwitchID: cxlSwitchID,
	}
}

// Find retrieves a CXL switch from the database based on the provided CmDb and CmFilter.
// It executes a Cypher query to fetch CXL switch information and processes the results into a structured map.
// The function returns a map representing the CXL switch and its resources, or an error if the operation fails.
func (nr *CXLSwitchRepository) Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error) {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", getCXLSwitch, nr.CXLSwitchID))
	cypherCursor, err := cmdb.CmDbExecCypher(getCXLSwitchColumnCount, getCXLSwitch, nr.CXLSwitchID)
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
		return compareByCXLSwitchSingle(records, i, j)
	})

	cxlSwitchWork := cxlswitch_model.NewCXLSwitch()
	for _, row := range records {
		// Retrieve the properties of the CXLSwitch and store them in the model
		cxlSwitchWork.Properties = row[getCXLSwitchIndexCXLSwitch].(*age.Vertex).Props()

		// From one record of the search result, assemble the Resource information associated with the CXLSwitch, and add it to the list of Resources
		// *If the assembled Resource information is empty, it is determined that there are no Resources associated with the CXLSwitch, and it is not added to the list of Resources
		resource := resource_repository.ComposeResource(
			row[getCXLSwitchIndexResource].(*age.Vertex),
			row[getCXLSwitchIndexAnnotation].(*age.Vertex),
			row[getCXLSwitchIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getCXLSwitchIndexNodeIDs].(*age.SimpleEntity),
			row[getCXLSwitchIndexNotDetected].(*age.SimpleEntity).AsBool(),
			false,
		)
		if resource.Validate() {
			cxlSwitchWork.Resources.Resources = append(cxlSwitchWork.Resources.Resources, resource)
		}
	}

	cxlSwitch := cxlswitch_model.NewCXLSwitch()
	if filter.FilterByCondition(cxlSwitchWork.ToObject()) {
		cxlSwitch = cxlSwitchWork
	}

	return cxlSwitch.ToObject(), nil
}

// compareByCXLSwitchSingle is a wrapper function to sort the contents of a Cypher query execution result.
// It leverages the compareByCXLSwitch function to determine the order of records based on CXL switch and resource information.
func compareByCXLSwitchSingle(records [][]age.Entity, i, j int) bool {
	return compareByCXLSwitch(records, getCXLSwitchIndexCXLSwitch, getCXLSwitchIndexResource, i, j)
}
