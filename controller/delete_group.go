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

	"github.com/gin-gonic/gin"
)

// DeleteGroup is a handler that deletes the specified group.
// If there are resources associated with the group, the group cannot be deleted.
// The default group cannot be deleted.
//
// Parameters:
// - c: gin.Context, the request context
//
// Process:
// 1. Logs the start of the request.
// 2. If the default group is specified, deletion is not allowed, and an error response is returned.
// 3. Searches for the group based on the specified group ID. If an error occurs, an error response is returned.
// 4. If the group does not exist, logs a warning and returns a 404 error response.
// 5. If the group contains resources, deletion is not allowed, logs a warning, and returns a 400 error response.
// 6. Deletes the group. If an error occurs, an error response is returned.
// 7. Logs the completion of the request and returns a 204 response.
func DeleteGroup(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "DeleteGroup"

	id := c.Param("id")
	// The default group cannot be deleted.
	if id == common.DefaultGroupId {
		errorDatial := "Default group specified error"
		common.Log.Error(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	filter := cmapi_filter.NewNoFilter()
	getRepository := cmapi_repository_group.NewGroupRepository(id, true)
	group, err := cmapi_repository.RelayFind(&getRepository, filter)
	if err != nil {
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if group == nil {
		errorDatial := "The target group for delete did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id))
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	// If there are resources associated with the group, the group cannot be deleted.
	resources := group["resources"].([]any)
	if len(resources) > 0 {
		errorDatial := "Group has resources error"
		common.Log.Warn(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	repository := cmapi_repository_group.NewDeleteGroupRepository(id)
	err = cmapi_repository.RelayDelete(&repository)
	if err != nil {
		errorDatial := "RelayDelete error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusNoContent, nil)
}
