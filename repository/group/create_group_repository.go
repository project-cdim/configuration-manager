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
	"errors"
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/model"

	"github.com/apache/age/drivers/golang/age"
	"github.com/google/uuid"
)

const (
	getResourceGroupCount = `
		MATCH (vrsg: ResourceGroups {id: '%s'})
		return COUNT(vrsg)
`
	getResourceGroupCountColumnCount = 1
)

const (
	mergeResourceGroup = `
		MERGE (vrsg: ResourceGroups {id: '%s'})
		SET vrsg = %s
`
	mergeResourceGroupColumnCount = 0
)

// CreateGroupRepository is a repository for creating groups.
// It provides methods to interact with the data source for group creation operations.
type CreateGroupRepository struct{}

// NewCreateGroupRepository creates a new instance of CreateGroupRepository.
// It returns an empty CreateGroupRepository struct.
func NewCreateGroupRepository() CreateGroupRepository {
	return CreateGroupRepository{}
}

// Set creates a new group in the database using the provided CmDb and CmModelMapper.
// It generates a unique resource group ID, converts the model to an object, and sets the ID.
// The object is then converted to a Cypher property map and used to execute a Cypher query
// to merge the resource group into the database.
//
// Parameters:
//
//	cmdb - The database connection object.
//	model  - The model mapper to convert the model to an object.
//
// Returns:
//
//	A map representing the created group object, or an error if the operation fails.
func (cgr *CreateGroupRepository) Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error) {
	id, err := generateResourceGroupID(cmdb)
	if err != nil {
		return nil, err
	}

	groupObject := model.ToObject()
	groupObject["id"] = id

	property, err := common.Map2CypherProperty(groupObject)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(mergeResourceGroup, id, property)
	common.Log.Debug(query)
	_, err = cmdb.CmDbExecCypher(mergeResourceGroupColumnCount, query)
	if err != nil {
		return nil, err
	}

	return groupObject, nil
}

// generateResourceGroupID generates a unique resource group ID using UUID version 7.
// It attempts to generate a unique ID up to 10 times, checking for duplicates in the database.
// If a unique ID is found, it is returned as a string. If no unique ID is found after 10 attempts,
// an error is returned.
//
// Parameters:
//   - cmdb: An instance of the CmDb database.
//
// Returns:
//   - A unique resource group ID as a string, or an error if a unique ID could not be generated.
func generateResourceGroupID(cmdb database.CmDb) (string, error) {
	for i := 0; i < 10; i++ {
		id, _ := uuid.NewV7()

		query := fmt.Sprintf(getResourceGroupCount, id.String())
		common.Log.Debug(query)
		cypherCursor, err := cmdb.CmDbExecCypher(getResourceGroupCountColumnCount, query)
		if err != nil {
			common.Log.Error(err.Error())
			return "", err
		}

		for cypherCursor.Next() {
			row, err := cypherCursor.GetRow()
			if err != nil {
				common.Log.Error(err.Error())
				return "", err
			}

			cntEntity := row[0].(*age.SimpleEntity)
			cnt := int(cntEntity.AsInt64())
			if cnt == 0 {
				cypherCursor.Close()
				return id.String(), nil
			}
			break
		}
		cypherCursor.Close()
	}

	return "", errors.New("generateResourceGroupID : An ID was generated in UUID format, but duplicates continued to occur")
}
