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

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	"github.com/project-cdim/configuration-manager/model"
)

const updateAnnotation string = `
	MATCH (vrs {deviceID: '%s'})-[ehv: Have]->(van:Annotation)
	SET van = %s
`

// UpdateAnnotationRepository is a repository that handles the updating of annotations.
// It contains the deviceID which is used to identify the specific device for which the annotations are being updated.
type UpdateAnnotationRepository struct {
	deviceID string
}

// NewUpdateAnnotationRepository creates a new instance of UpdateAnnotationRepository with the provided device ID.
//
// Parameters:
//   - deviceID: A string representing the unique identifier of the device.
//
// Returns:
//   - UpdateAnnotationRepository: A new instance of UpdateAnnotationRepository initialized with the given device ID.
func NewUpdateAnnotationRepository(deviceID string) UpdateAnnotationRepository {
	return UpdateAnnotationRepository{
		deviceID: deviceID,
	}
}

// Set updates the resource annotation for the repository.
// This method is responsible for updating the annotation information associated with a specific device in the repository.
// It takes a CmDb interface for database operations and a CmModelMapper for mapping the model to the database schema.
//
// The method extracts the deviceID from the model, constructs an update query, and executes it against the database.
// If the operation is successful, it returns nil; otherwise, it returns an error.
//
// Note: This method currently includes a workaround for including the deviceID in the annotation object.
// This will be unnecessary once the deviceID is removed from the annotation object in the database schema.
func (uar *UpdateAnnotationRepository) Set(cmdb database.CmDb, model model.CmModelMapper) (map[string]any, error) {
	annotationObject := model.ToObject()

	cypherProperty, err := common.Map2CypherProperty(annotationObject)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(updateAnnotation, uar.deviceID, cypherProperty)

	_, err = cmdb.CmDbExecCypher(0, query)
	if err != nil {
		return nil, err
	}

	return annotationObject, nil
}
