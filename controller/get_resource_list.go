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

// GetResourceList retrieves all resources from the system. It begins by logging the start of the request.
// A filter is created to specify the criteria for retrieving resources, with a flag indicating whether to include
// detailed information. The 'detail' query parameter is extracted from the request to determine the level of detail
// required in the response. A repository for managing resource lists is then instantiated with this detail level.
// The function attempts to find all resources that match the filter criteria. If an error occurs during this process,
// it logs the error and returns an error response. On successful retrieval, a response object is constructed containing
// the count of resources found and the list of resources themselves. This response object is then marshaled into JSON.
// If marshaling fails, the error is logged and an error response is returned. Otherwise, the marshaled JSON is logged
// for debugging purposes. Finally, the function logs the successful completion of the request and returns the marshaled
// JSON as a response with a 200 OK status.
func GetResourceList(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetResourceList"

	filter := cmapi_filter.NewNoFilter()

	// Retrieve query parameter: detail
	detail, err := getBoolQueryParam(c, "detail")
	if err != nil {
		errorDatial := "getBoolQueryParam error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	repository := cmapi_repository_resource.NewResourceListRepository(detail)
	resources, err := cmapi_repository.RelayFindList(&repository, filter)
	if err != nil {
		// In case of an error during the retrieval or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFindList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	res := gin.H{
		"count":     len(resources),
		"resources": resources,
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	// Sets the return value
	c.JSON(http.StatusOK, res)
}
