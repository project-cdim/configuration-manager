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
)

// RepositoryListFinder is the interface to retrieve the model list.
type RepositoryListFinder interface {
	FindList(cmdb database.CmDb, filter filter.CmFilter) ([]map[string]any, error)
}

// RepositoryFinder is the interface to retrieve the model.
type RepositoryFinder interface {
	Find(cmdb database.CmDb, filter filter.CmFilter) (map[string]any, error)
}

// RepositorySetter is the interface for insert or update to repository.
type RepositorySetter interface {
	Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error)
}

// RepositoryDeleter is the interface for delete to repository.
type RepositoryDeleter interface {
	Delete(cmdb database.CmDb) error
}
