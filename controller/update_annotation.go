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
        
package controller

import (
	"fmt"
	"net/http"

	"github.com/project-cdim/configuration-manager/common"
	cmapi_filter "github.com/project-cdim/configuration-manager/filter"
	cmapi_model_annotation "github.com/project-cdim/configuration-manager/model/annotation"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_annotation "github.com/project-cdim/configuration-manager/repository/annotation"
	cmapi_repository_resource "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/gin-gonic/gin"
)

// UpdateAnnotation handles the HTTP request to update the annotation of a resource.
// It performs the following steps:
// 1. Reads the JSON body of the request and converts it into a map representing the annotation properties.
// 2. Checks if the resource associated with the given ID exists.
// 3. If the resource exists, it updates the annotation with the new properties.
// 4. Returns the updated annotation as a JSON response.
//
// Parameters:
// - c *gin.Context: The context of the Gin HTTP request, which contains the request and response objects.
//
// Responses:
// - 200 OK with the updated annotation in JSON format if the update is successful.
// - Appropriate error response (e.g., 404 Not Found, 500 Internal Server Error) if the resource does not exist or if there is an error during the update process.
//
// This function is part of the REST API that allows clients to update annotations for resources in the system.
func UpdateAnnotation(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "UpdateAnnotation"

	id := c.Param("id")
	// Reads the JSON from the RequestBody and expands it into a Map variable
	annotationProperties, err := unmarshalRequestBodyForMap(c)
	if err != nil {
		errorDatial := "unmarshalRequestBodyForMap error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Checks if the resource associated with the target annotation exists by searching once
	filter := cmapi_filter.NewNoFilter()
	getRepository := cmapi_repository_resource.NewResourceRepository(id)
	resource, err := cmapi_repository.RelayFind(&getRepository, filter)
	if err != nil {
		// In case of an error during the graph DB search or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if resource == nil {
		// If the target resource for update did not exist
		errorDatial := "The target resource for update did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id), false)
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	annotation := cmapi_model_annotation.NewAnnotation()
	annotation.Properties = annotationProperties
	repository := cmapi_repository_annotation.NewUpdateAnnotationRepository(id)
	res, err := cmapi_repository.RelaySet(&repository, &annotation)
	if err != nil {
		errorDatial := "RelaySet error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
