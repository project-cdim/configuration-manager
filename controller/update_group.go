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
	cmapi_model_group "github.com/project-cdim/configuration-manager/model/group"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_group "github.com/project-cdim/configuration-manager/repository/group"

	"github.com/gin-gonic/gin"
)

// UpdateGroup is a handler function to update group information.
// This function processes the following steps:
// 1. Converts the request body to a map.
// 2. Prohibits updating the default group.
// 3. Validates the request body.
// 4. Retrieves existing group information based on the group ID.
// 5. Returns a 404 error if the target group for update does not exist.
// 6. Creates and updates new group information.
// 7. Returns a response with a 200 status code if the update is successful.
//
// Parameters:
// - c: gin.Context - Request context
//
// Response:
// - On success: 200 status code with the updated group information
// - On validation error: 400 status code
// - If the group does not exist: 404 status code
// - On server error: 500 status code
func UpdateGroup(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "UpdateGroup"

	id := c.Param("id")
	properties, err := unmarshalRequestBodyForMap(c)
	if err != nil {
		errorDatial := "unmarshalRequestBodyForMap error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Updating the default group is not allowed
	if id == common.DefaultGroupId {
		errorDatial := "Default group specified error"
		common.Log.Error(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Validation of the requestBody
	if !cmapi_model_group.ValidateProperty(properties) {
		errorDatial := "Validation error"
		common.Log.Error(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	filter := cmapi_filter.NewNoFilter()
	getRepository := cmapi_repository_group.NewGroupRepository(id, false)
	groupFromDb, err := cmapi_repository.RelayFind(&getRepository, filter)
	if err != nil {
		// In case of an error during the graph DB search or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if groupFromDb == nil {
		// If the target group for update did not exist
		errorDatial := "The target group for update did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id), false)
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	group := cmapi_model_group.NewGroupForUpdate(groupFromDb, properties)
	repository := cmapi_repository_group.NewUpdateGroupRepository()
	res, err := cmapi_repository.RelaySet(&repository, &group)
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
