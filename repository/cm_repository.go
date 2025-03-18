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
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/filter"
	"github.com/project-cdim/configuration-manager/model"

	"github.com/apache/age/drivers/golang/age"
)

// RelayFindList returns a list of models that match the conditions.
// It is a generic function, so the model returned changes based on the parameter.
func RelayFindList[T RepositoryListFinder](repo T, filter filter.CmFilter) ([]map[string]any, error) {
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

	return res, nil
}

// RelayFind returns a model that matches the conditions.
// It is a generic function, so the model returned changes based on the parameter.
func RelayFind[T RepositoryFinder](repo T, filter filter.CmFilter) (map[string]any, error) {
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

	return res, nil
}

// RelaySet registers a model.
// It is a generic function, so the model passed changes based on the parameter.
func RelaySet[T RepositorySetter](repo T, model model.CmModelMapper) (map[string]any, error) {
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

// RelayDelete delete a model.
// It is a generic function, so the model passed changes based on the parameter.
func RelayDelete[T RepositoryDeleter](repo T) error {
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

// ExtractEntityString converts a SimpleEntity to a string and returns it.
func ExtractEntityString(entity *age.SimpleEntity) string {
	return entity.AsStr()
}

// ExtractEntitySlice converts a SimpleEntity to a slice of strings and returns it.
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
