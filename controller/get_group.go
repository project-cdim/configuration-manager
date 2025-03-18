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

// GetGroup retrieves the group information for the specified ID and returns it in JSON format.
//
// Parameters:
//   - c: gin.Context. Represents the request context.
//
// Processing flow:
//  1. Logs the request URL path and method.
//  2. Initializes the filter and repository.
//  3. Retrieves the group information from the repository.
//  4. If an error occurs, logs the error message and returns a response with HTTP status 500.
//  5. If the group information is not found, logs a warning message and returns a response with HTTP status 404.
//  6. Serializes the group information into JSON format.
//  7. If an error occurs during serialization, logs the error message and returns a response with HTTP status 500.
//  8. If processing completes successfully, returns the group information as a response with HTTP status 200.
func GetGroup(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetGroup"

	id := c.Param("id")
	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_group.NewGroupRepository(id, true)
	res, err := cmapi_repository.RelayFind(&repository, filter)
	if err != nil {
		// In case of an error during the retrieval or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if res == nil {
		errorDatial := "No search results"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id))
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
