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

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/model"
)

// UpdateGroupRepository is a repository that handles the update operations for groups.
// It provides methods to update group information in the data store.
type UpdateGroupRepository struct{}

// NewUpdateGroupRepository creates a new instance of UpdateGroupRepository.
// It returns an empty UpdateGroupRepository struct.
func NewUpdateGroupRepository() UpdateGroupRepository {
	return UpdateGroupRepository{}
}

// Set updates a group in the database using the provided CmDb and CmModelMapper.
// It converts the model to an object, generates a Cypher property map, and executes
// a Cypher query to merge the resource group.
//
// Parameters:
//
//	cmdb - The database connection object.
//	model  - The model mapper to convert the group model to an object.
//
// Returns:
//
//	A map representing the updated group object and an error if any occurred during the process.
func (ugr *UpdateGroupRepository) Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error) {
	groupObject := model.ToObject()
	id := groupObject["id"]

	property, err := common.Map2CypherProperty(groupObject)
	if err != nil {
		return nil, err
	}

	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", mergeResourceGroup, id, property))
	_, err = cmdb.CmDbExecCypher(mergeResourceGroupColumnCount, mergeResourceGroup, id, property)
	if err != nil {
		return nil, err
	}

	return groupObject, nil
}
