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
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_group "github.com/project-cdim/configuration-manager/repository/group"
	cmapi_repository_resource "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/gin-gonic/gin"
)

// AssignResourceToGroup handles the update of a resource group.
// It expects a JSON body containing a single group ID to update the resource with.
// If multiple group IDs are provided, it returns a BadRequest error.
// It performs the following steps:
// 1. Logs the start of the request.
// 2. Binds the JSON body to a slice of strings.
// 3. Validates that only one group ID is provided.
// 4. Retrieves the resource by its ID from the repository.
// 5. Checks if the resource exists, returning a NotFound error if it does not.
// 6. Retrieves the group by its ID from the repository.
// 7. Checks if the group exists, returning a BadRequest error if it does not.
// 8. Updates the resource with the new group ID.
// 9. Logs the successful completion of the request.
// 10. Returns the updated resource in the response.
//
// Parameters:
// - c: The Gin context, which provides request and response handling.
//
// Responses:
// - 200 OK: The resource was successfully updated.
// - 400 BadRequest: The request body was invalid or multiple group IDs were provided.
// - 404 NotFound: The resource or group was not found.
// - 500 InternalServerError: An error occurred during the update process.
func AssignResourceToGroup(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "AssignResourceToGroup"

	id := c.Param("id")
	var targetGroups []string
	if err := c.ShouldBindJSON(&targetGroups); err != nil {
		errorDatial := "BindJson error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Multiple resource groups specified in the request body are not allowed
	if len(targetGroups) > 1 {
		errorDatial := "multiple groups specified error"
		common.Log.Warn(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	filter := cmapi_filter.NewNoFilter()
	resourceRepository := cmapi_repository_resource.NewResourceRepository(id)
	resource, err := cmapi_repository.RelayFind(&resourceRepository, filter)
	if err != nil {
		// In case of an error during the graph DB search or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if resource == nil {
		// If the target group for update did not exist
		errorDatial := "The target resource for update did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id))
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	device := resource["device"].(map[string]any)
	deviceType := hwResourceType(device["type"].(string))
	dbDeciceType, err := deviceType.convertToDBLabel()
	if err != nil {
		errorDatial := "device type is invalid"
		common.Log.Error(fmt.Sprintf("%s %s : %v", funcName, errorDatial, deviceType), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	groupRepository := cmapi_repository_group.NewGroupRepository(targetGroups[0], false)
	group, err := cmapi_repository.RelayFind(&groupRepository, filter)
	if err != nil {
		// In case of an error during the graph DB search or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if group == nil {
		// If the target group for update did not exist
		errorDatial := "The target group for update did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, targetGroups[0]))
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	repository := cmapi_repository_resource.NewAssignResourceToGroupRepository(id, dbDeciceType, targetGroups)
	groupIDs, err := cmapi_repository.RelaySet(&repository, nil)
	if err != nil {
		errorDatial := "RelaySet error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	res := gin.H{
		"count":            len(groupIDs["resourceGroupIDs"].([]string)),
		"resourceGroupIDs": groupIDs["resourceGroupIDs"],
	}
	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
