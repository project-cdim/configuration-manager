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

package annotation_repository

import (
	"fmt"
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/model"
	resource_repository "github.com/project-cdim/configuration-manager/repository/resource"
)

const updateAnnotation string = `
	MATCH (vrs {deviceID: '%s'})-[:Have]->(van:Annotation)
	WHERE %s
	SET van = %s
`

const updateAnnotationWhereParts string = "'%s' IN labels(vrs)"

// updateAnnotationConstructsWhereClause constructs a WHERE clause for Cypher queries to filter vertices
// based on resource types. It generates an OR condition that checks if any of the resource types
// from ResourceTypeList exist as labels on the 'vrs' vertices variable.
// Returns a string containing the complete WHERE clause condition.
func updateAnnotationConstructsWhereClause() string {
	var whereClauses []string
	for _, resourceType := range resource_repository.ResourceTypeList {
		whereClauses = append(whereClauses, fmt.Sprintf(updateAnnotationWhereParts, resourceType))
	}
	return strings.Join(whereClauses, " OR ")
}

// UpdateAnnotationRepository is a struct that holds the device IDs to be updated.
// It is used to update the annotations for the specified devices.
type UpdateAnnotationRepository struct {
	deviceIDs []string
}

// NewUpdateAnnotationRepository creates a new UpdateAnnotationRepository with the given device IDs.
// It returns an UpdateAnnotationRepository instance.
func NewUpdateAnnotationRepository(deviceIDs []string) UpdateAnnotationRepository {
	return UpdateAnnotationRepository{
		deviceIDs: deviceIDs,
	}
}

// Set updates annotations in the database based on the provided model and device IDs.
// It converts the model to a Cypher property map and then iterates through the device IDs,
// executing an update query for each device ID.
//
// Parameters:
//   - cmdb: A database connection implementing the database.CmDb interface.
//   - model: A model implementing the model.CmModelMapper interface, representing the annotation data.
//
// Returns:
//   - A map[string]any representing the updated annotation object, or nil if an error occurs.
//   - An error if any operation fails during the update process.
func (uar *UpdateAnnotationRepository) Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error) {
	annotationObject := model.ToObject()

	cypherProperty, err := common.Map2CypherProperty(annotationObject)
	if err != nil {
		return nil, err
	}

	whereClauses := updateAnnotationConstructsWhereClause()
	if uar.deviceIDs != nil {
		for _, deviceIDs := range uar.deviceIDs {
			common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", updateAnnotation, deviceIDs, whereClauses, cypherProperty))
			_, err = cmdb.CmDbExecCypher(0, updateAnnotation, deviceIDs, whereClauses, cypherProperty)
			if err != nil {
				return nil, err
			}
		}
	}

	return annotationObject, nil
}
