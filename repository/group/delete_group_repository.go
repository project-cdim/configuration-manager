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
)

const deleteGroup = `
	MATCH (vrsg: ResourceGroups {id: '%s'})
	DELETE vrsg
`

// DeleteGroupRepository represents a repository for deleting a group.
// It contains the GroupID which identifies the group to be deleted.
type DeleteGroupRepository struct {
	GroupID string
}

// NewDeleteGroupRepository creates a new instance of DeleteGroupRepository with the specified groupID.
//
// Parameters:
//   - groupID: A string representing the unique identifier of the group to be deleted.
//
// Returns:
//
//	A new instance of DeleteGroupRepository initialized with the provided groupID.
func NewDeleteGroupRepository(groupID string) DeleteGroupRepository {
	return DeleteGroupRepository{
		GroupID: groupID,
	}
}

// Delete removes a group from the database based on the GroupID of the DeleteGroupRepository instance.
// It executes a Cypher query to delete the group and logs the query for debugging purposes.
// If the deletion fails, it returns an error.
//
// Parameters:
//
//	cmdb - An instance of the CmDb database interface.
//
// Returns:
//
//	error - An error object if the deletion fails, otherwise nil.
func (dgr *DeleteGroupRepository) Delete(cmdb database.CmDb) error {
	query := fmt.Sprintf(deleteGroup, dgr.GroupID)
	common.Log.Debug(query)
	_, err := cmdb.CmDbExecCypher(0, query)
	if err != nil {
		return err
	}

	return nil
}
