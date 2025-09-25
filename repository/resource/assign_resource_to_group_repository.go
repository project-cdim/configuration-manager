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

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/model"
)

const (
	deleteIncludeEdge = `
		MATCH (vrsg)-[ein:Include]->({deviceID: '%s'})
		DELETE ein
`
	deleteIncludeEdgeCount = 0
)

const (
	createIncludeEdge = `
		MATCH (vrs:%s {deviceID: '%s'})
		MATCH (vrsg:ResourceGroups {id: '%s'})
		CREATE (vrsg)-[:Include]->(vrs)
`
	createIncludeEdgeCount = 0
)

// AssignResourceToGroupRepository represents a repository for updating resource groups associated with a device.
// It contains the device ID, the type of the device in the database, and the new resource groups to be assigned.
type AssignResourceToGroupRepository struct {
	DeviceID          string
	DbDeviceType      string
	NewResourceGroups []string
}

// NewAssignResourceToGroupRepository creates a new instance of UpdateGroupOfResourceRepository with the provided device ID,
// database device type, and new resource groups.
//
// Parameters:
//   - deviceID: A string representing the unique identifier of the device.
//   - dbDeviceType: A string representing the type of the device in the database.
//   - newResourceGroups: A slice of strings representing the new resource groups to be associated with the device.
//
// Returns:
//
//	An instance of AssignResourceToGroupRepository initialized with the provided values.
func NewAssignResourceToGroupRepository(deviceID string, dbDeviceType string, newResourceGroups []string) AssignResourceToGroupRepository {
	return AssignResourceToGroupRepository{
		DeviceID:          deviceID,
		DbDeviceType:      dbDeviceType,
		NewResourceGroups: newResourceGroups,
	}
}

// Set updates the resource groups associated with a device in the database.
// It first deletes existing associations and then creates new ones based on the provided model.
//
// Parameters:
// - cmdb: The database connection object.
// - model: The model containing the new resource group associations.
//
// Returns:
// - A map containing the device ID and the new resource group IDs.
// - An error if any database operation fails.
func (argr *AssignResourceToGroupRepository) Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error) {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", deleteIncludeEdge, argr.DeviceID))
	_, err := cmdb.CmDbExecCypher(deleteIncludeEdgeCount, deleteIncludeEdge, argr.DeviceID)
	if err != nil {
		return nil, err
	}

	for _, resourceGroupID := range argr.NewResourceGroups {
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", createIncludeEdge, argr.DbDeviceType, argr.DeviceID, resourceGroupID))
		_, err = cmdb.CmDbExecCypher(createIncludeEdgeCount, createIncludeEdge, argr.DbDeviceType, argr.DeviceID, resourceGroupID)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		"resourceGroupIDs": argr.NewResourceGroups,
	}, nil
}
