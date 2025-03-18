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
	cmapi_repository_resource "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/gin-gonic/gin"
)

// GetResource retrieves a specific resource by its ID. It logs the start of the request and creates a filter for resource retrieval.
// The function then instantiates a repository with the resource ID obtained from the request parameters.
// It attempts to find the resource using the repository and filter. If an error occurs, it logs the error and returns an error response.
// If no resource is found, it logs a warning and returns a 404 Not Found response. On successful retrieval of the resource,
// it attempts to marshal the resource into JSON. If marshaling fails, it logs the error and returns an error response.
// Otherwise, it logs the marshaled JSON for debugging purposes and the successful completion of the request before returning the resource as a JSON response with a 200 OK status.
func GetResource(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetResource"

	id := c.Param("id")
	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_resource.NewResourceRepository(id)
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
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id), false)
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
