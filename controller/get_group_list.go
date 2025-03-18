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

// GetGroupList is a handler that retrieves the requested group list and returns a JSON response.
//
// Parameters:
// - c: gin.Context. The context of the HTTP request and response.
//
// Processing flow:
// 1. Logs the start of the request.
// 2. Retrieves the query parameter "withResources".
// 3. Creates a repository based on the "withResources" parameter.
// 4. Retrieves the group list from the repository.
// 5. Serializes the retrieved group list into JSON format.
// 6. Returns the serialized data as a response.
// 7. Logs the completion of the request.
//
// Error handling:
// - If the retrieval of the query parameter fails, returns 400 Bad Request.
// - If the retrieval of the group list fails, returns 500 Internal Server Error.
// - If the JSON serialization fails, returns 500 Internal Server Error.
func GetGroupList(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetGroupList"

	filter := cmapi_filter.NewNoFilter()

	// Retrieve query parameter: withResources
	withResources, err := getBoolQueryParam(c, "withResources")
	if err != nil {
		errorDatial := "getBoolQueryParam error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	repository := cmapi_repository_group.NewGroupListRepository(withResources)
	groups, err := cmapi_repository.RelayFindList(&repository, filter)
	if err != nil {
		// In case of an error during the retrieval or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFindList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	res := gin.H{
		"count":          len(groups),
		"resourceGroups": groups,
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
