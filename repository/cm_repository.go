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

package repository

import (
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	"github.com/project-cdim/configuration-manager/model"

	"github.com/apache/age/drivers/golang/age"
)

// RelayFindList relays the FindList operation to the provided RepositoryListFinder.
// It starts a database transaction, calls the FindList method on the repository,
// and handles transaction management (commit/rollback) and connection closing.
//
// Parameters:
//   - repo: A RepositoryListFinder instance that implements the FindList method.
//   - filter: A filter.CmFilter instance containing the filter criteria for the FindList operation.
//
// Returns:
//   - A slice of any, representing the list of found items.
//   - An error, if any occurred during the operation.
func RelayFindList(repo RepositoryListFinder, filter filter.CmFilter) ([]any, error) {
	cmdb := database.NewCmDb()

	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		return nil, err
	}
	defer cmdb.CmDbDisconnection()

	res, err := repo.FindList(cmdb, filter)
	if err != nil {
		return nil, err
	}

	unquoted, err := common.UnquoteRecursive(res)
	if err != nil {
		// This error should not occur because the string to be unquoted is always quoted beforehand.
		return nil, err
	}

	unquotedSlice, ok := unquoted.([]any)
	if !ok {
		common.Log.Error(fmt.Sprintf("RelayFindList error: expected []any, got %T", unquoted))
		return nil, fmt.Errorf("RelayFindList error: expected []any, got %T", unquoted)
	}

	return unquotedSlice, nil
}

// RelayFind finds configuration data using the provided RepositoryFinder and filter.
// It manages a database transaction and ensures proper disconnection.
//
// Parameters:
//   - repo: A RepositoryFinder instance used to find the configuration data.
//   - filter: A filter.CmFilter instance used to specify the search criteria.
//
// Returns:
//   - A map[string]any containing the found configuration data, or nil if an error occurred.
//   - An error if any error occurred during the process, otherwise nil.
func RelayFind(repo RepositoryFinder, filter filter.CmFilter) (map[string]any, error) {
	cmdb := database.NewCmDb()

	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		return nil, err
	}
	defer cmdb.CmDbDisconnection()

	res, err := repo.Find(cmdb, filter)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil // No results found, return nil
	}

	unquoted, err := common.UnquoteRecursive(res)
	if err != nil {
		// This error should not occur because the string to be unquoted is always quoted beforehand.
		return nil, err
	}

	unquotedMap, ok := unquoted.(map[string]any)
	if !ok {
		common.Log.Error(fmt.Sprintf("RelayFind error: expected map[string]any, got %T", unquoted))
		return nil, fmt.Errorf("RelayFind error: expected map[string]any, got %T", unquoted)
	}

	return unquotedMap, nil
}

// RelaySet sets the configuration model using the provided repository setter and model mapper.
// It starts a database transaction, sets the configuration, commits the transaction, and returns the result.
// If any error occurs during the process, it rolls back the transaction.
//
// Parameters:
//   - repo: RepositorySetter interface for setting the configuration model.
//   - model: model.CmModelMapper containing the configuration data.
//
// Returns:
//   - map[string]any: The result of setting the configuration model.
//   - error: An error if any occurred during the process.
func RelaySet(repo RepositorySetter, model model.CmModelMapper) (map[string]any, error) {
	cmdb := database.NewCmDb()

	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		return nil, err
	}
	defer cmdb.CmDbDisconnection()

	res, err := repo.Set(cmdb, model)
	if err != nil {
		cmdb.CmDbRollback()
		return nil, err
	}

	err = cmdb.CmDbCommit()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// RelayDelete deletes a repository using the provided RepositoryDeleter.
// It manages a database transaction, ensuring atomicity of the delete operation.
//
// Parameters:
//   - repo: A RepositoryDeleter interface that provides the Delete method.
//
// Returns:
//   - error: An error if any operation fails during the process, including starting, committing, or rolling back the transaction. Returns nil if the deletion is successful.
func RelayDelete(repo RepositoryDeleter) error {
	cmdb := database.NewCmDb()

	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		return err
	}
	defer cmdb.CmDbDisconnection()

	err = repo.Delete(cmdb)
	if err != nil {
		cmdb.CmDbRollback()
		return err
	}

	err = cmdb.CmDbCommit()
	if err != nil {
		return err
	}

	return nil
}

// ExtractEntityString extracts the string representation of a SimpleEntity.
//
// Parameters:
//   - entity: A pointer to the SimpleEntity to extract the string from.
//
// Returns:
//   - string: The string representation of the SimpleEntity.
func ExtractEntityString(entity *age.SimpleEntity) string {
	return entity.AsStr()
}

// ExtractEntitySlice extracts a slice of strings from a SimpleEntity.
//
// Parameters:
//   - entity: A pointer to an age.SimpleEntity from which to extract the string slice.
//
// Returns:
//   - []string: A slice containing the string values extracted from the entity.
//     Only string values are included in the result. Non-string values are ignored.
func ExtractEntitySlice(entity *age.SimpleEntity) []string {
	anyArr := entity.AsArr()

	res := []string{}
	for _, anyValue := range anyArr {
		strValue, ok := anyValue.(string)
		if ok {
			res = append(res, strValue)
		}
	}
	return res
}
