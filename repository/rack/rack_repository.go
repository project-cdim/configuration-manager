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

package rack_repository

import (
	"fmt"
	"sort"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	chassis_model "github.com/project-cdim/configuration-manager/model/chassis"
	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	rack_model "github.com/project-cdim/configuration-manager/model/rack"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/apache/age/drivers/golang/age"
)

// getRack is cypher query to retrieve a specific rack.
const getRack string = `
	MATCH (vrc:Rack{id: '%s'})
	OPTIONAL MATCH (vrc)-[:Attach]->(vch)
	OPTIONAL MATCH (vch)-[:Mount]->(vrs)
	OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
	OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
	OPTIONAL MATCH (vrs)-[:Have]->(van:Annotation)
	OPTIONAL MATCH (vrs)-[endt:NotDetected]->(:NotDetectedDevice)
	WITH vrc, vch, vrs, van, vrsg, vnd, endt
	RETURN
		vrc,
		CASE WHEN vch IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vch END,
		CASE WHEN vrs IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE vrs END,
		CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
		COLLECT(vrsg.id),
		COLLECT(vnd.id),
		CASE WHEN endt IS NULL THEN true ELSE false END
`

const getRackColumnCount = 7
const (
	getRackIndexRack = iota
	getRackIndexChassis
	getRackIndexDevice
	getRackIndexAnnotation
	getRackIndexResourceGroupIDs
	getRackIndexNodeIDs
	getRackIndexNotDetected
)

// RackRepository is a repository structure for getting a specific rack.
type RackRepository struct {
	RackID string
	Detail bool
}

// NewRackRepository creates and returns a RackRepository object that holds the argument rackID.
// This constructor function initializes a RackRepository with a rackID and a detail flag indicating whether to retrieve detailed rack information.
func NewRackRepository(rackID string, detail bool) RackRepository {
	return RackRepository{
		RackID: rackID,
		Detail: detail,
	}
}

// Find retrieves a rack based on the provided CmDb and CmFilter.
// It constructs and executes a Cypher query to fetch rack information from the database.
// The function processes the query results, assembling them into a structured map representing the rack and its components.
// If successful, it returns the assembled rack as a map[string]any, or an error if the operation fails.
func (rr *RackRepository) Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error) {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", getRack, rr.RackID))
	cypherCursor, err := cmdb.CmDbExecCypher(getRackColumnCount, getRack, rr.RackID)
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
		return compareByRackSingle(records, i, j)
	})

	rack := rack_model.NewRack()
	chassis := chassis_model.NewChassis()
	resources := resource_model.NewResourceList()
	cxlswitches := cxlswitch_model.NewCXLSwitchList()
	preChassisID := ""
	for _, row := range records {
		if len(rack.Properties) <= 0 {
			rack.Properties = row[getRackIndexRack].(*age.Vertex).Props()
		}

		chassisWork := row[getRackIndexChassis].(*age.Vertex).Props()
		chassisID, ok := chassisWork["id"].(string)
		if !ok {
			continue
		}

		if len(preChassisID) > 0 && preChassisID != chassisID {
			chassis.Resources = resources
			chassis.CXLSwitches = cxlswitches
			if filter.FilterByCondition(chassis.ToObject()) {
				rack.Chassis.Chassis = append(rack.Chassis.Chassis, chassis)
			}
			chassis = chassis_model.NewChassis()
			resources = resource_model.NewResourceList()
			cxlswitches = cxlswitch_model.NewCXLSwitchList()
		}

		chassis.Properties = chassisWork

		deviceVertex := row[getRackIndexDevice].(*age.Vertex)
		if len(deviceVertex.Props()) > 0 {
			if deviceVertex.Label() == "CXLswitch" {
				cxlswitch := cxlswitch_model.NewCXLSwitch()
				cxlswitch.Properties = deviceVertex.Props()
				if cxlswitch.Validate() {
					cxlswitches.CXLSwitches = append(cxlswitches.CXLSwitches, cxlswitch)
				}
			} else {
				resource := resource_repository.ComposeResource(
					deviceVertex,
					row[getRackIndexAnnotation].(*age.Vertex),
					row[getRackIndexResourceGroupIDs].(*age.SimpleEntity),
					row[getRackIndexNodeIDs].(*age.SimpleEntity),
					row[getRackIndexNotDetected].(*age.SimpleEntity).AsBool(),
					rr.Detail,
				)
				if resource.Validate() {
					resources.Resources = append(resources.Resources, resource)
				}
			}
		}

		preChassisID = chassisID
	}

	chassis.Resources = resources
	chassis.CXLSwitches = cxlswitches
	if filter.FilterByCondition(chassis.ToObject()) {
		rack.Chassis.Chassis = append(rack.Chassis.Chassis, chassis)
	}

	return rack.ToObject(), nil
}

// compareByRackSingle is a wrapper function to sort the contents of a Cypher query execution result.
// It leverages the compareByRack function to determine the order of records based on chassis and device information.
func compareByRackSingle(records [][]age.Entity, i, j int) bool {
	return compareByRack(records, getRackIndexChassis, getRackIndexDevice, i, j)
}
