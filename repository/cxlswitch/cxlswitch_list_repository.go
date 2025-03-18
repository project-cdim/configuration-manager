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
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

// Cypher query to retrieve a list of CXL Switches
// Retrieves the Vertex for the CXL Switch and any connected resources and related Vertices
const getCXLSwitchList string = `
	MATCH (vcx: CXLswitch)
	OPTIONAL MATCH (vcx)-[ecn: Connect]->(vrs)
	OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
	OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
	OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
	OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
	WITH vcx, vrs, van, vrsg, vnd, endt
	RETURN
		vcx,
		CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END, 
		CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
		COLLECT(vrsg.id),
		COLLECT(vnd.id),
		CASE WHEN endt IS NULL THEN true ELSE false END
`
const getCXLSwitchListColumnCount = 6
const (
	getCXLSwitchListIndexCXLSwitch = iota
	getCXLSwitchListIndexResource
	getCXLSwitchListIndexAnnotation
	getCXLSwitchListIndexResourceGroupIDs
	getCXLSwitchListIndexNodeIDs
	getCXLSwitchListIndexNotDetected
)

// CXLSwitchListRepository is a repository structure for getting CXL Switch lists.
type CXLSwitchListRepository struct{}

// NewCXLSwitchListRepository creates and returns a new CXLSwitchListRepository object.
// This function initializes a CXLSwitchListRepository with default values.
func NewCXLSwitchListRepository() CXLSwitchListRepository {
	return CXLSwitchListRepository{}
}

// FindList returns a list of CXL Switches from the database based on the provided CmDb and CmFilter.
// It executes a Cypher query to fetch a list of CXL switches and processes the results into a slice of maps.
// Each map in the slice represents a CXL switch and its associated resources.
// The function returns a slice of maps representing the CXL switches and their resources, or an error if the operation fails.
func (nlr *CXLSwitchListRepository) FindList(cmdb database.CmDb, filter filter.CmFilter) ([]map[string]any, error) {
	common.Log.Debug(getCXLSwitchList)
	cypherCursor, err := cmdb.CmDbExecCypher(getCXLSwitchListColumnCount, getCXLSwitchList)
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
		return compareByCXLSwitchList(records, i, j)
	})

	cxlSwitchList := cxlswitch_model.NewCXLSwitchList()
	cxlSwitch := cxlswitch_model.NewCXLSwitch()
	resources := resource_model.NewResourceList()
	preCXLSwitchID := ""
	for _, row := range records {
		cxlSwitchWork := row[getCXLSwitchListIndexCXLSwitch].(*age.Vertex).Props()
		cxlSwitchID := cxlSwitchWork["id"].(string)

		if len(preCXLSwitchID) > 0 && preCXLSwitchID != cxlSwitchID {
			cxlSwitch.Resources = resources
			if filter.FilterByCondition(cxlSwitch.ToObject()) {
				cxlSwitchList.CXLSwitches = append(cxlSwitchList.CXLSwitches, cxlSwitch)
			}
			cxlSwitch = cxlswitch_model.NewCXLSwitch()
			resources = resource_model.NewResourceList()
		}

		cxlSwitch.Properties = cxlSwitchWork

		resource := resource_repository.ComposeResource(
			row[getCXLSwitchListIndexResource].(*age.Vertex),
			row[getCXLSwitchListIndexAnnotation].(*age.Vertex),
			row[getCXLSwitchListIndexResourceGroupIDs].(*age.SimpleEntity),
			row[getCXLSwitchListIndexNodeIDs].(*age.SimpleEntity),
			row[getCXLSwitchListIndexNotDetected].(*age.SimpleEntity).AsBool(),
			true,
		)
		if resource.Validate() {
			resources.Resources = append(resources.Resources, resource)
		}

		preCXLSwitchID = cxlSwitchID
	}

	cxlSwitch.Resources = resources
	if filter.FilterByCondition(cxlSwitch.ToObject()) {
		cxlSwitchList.CXLSwitches = append(cxlSwitchList.CXLSwitches, cxlSwitch)
	}

	return cxlSwitchList.ToObject(), nil
}

// compareByCXLSwitchList is a wrapper function to sort the contents of a Cypher query execution result.
// It leverages the compareByCXLSwitch function to determine the order of records based on CXL switch and resource information.
// This function is used to order the records by CXL switch ID and associated resources, ensuring that the list is sorted correctly.
func compareByCXLSwitchList(records [][]age.Entity, i, j int) bool {
	return compareByCXLSwitch(records, getCXLSwitchListIndexCXLSwitch, getCXLSwitchListIndexResource, i, j)
}
