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
	cmapi_filter_resource "github.com/project-cdim/configuration-manager/filter/resource"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_resource "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/gin-gonic/gin"
)

// GetUnusedResourceList retrieves a list of unused resources based on the provided resource group IDs.
// It starts by logging the start of the request. It then extracts the 'resourceGroupID' query parameters from the request URL.
// A filter is created to specify the criteria for unused resources, which includes setting the availability and visibility flags to true,
// and including the specified resource group IDs. A repository for the resource list is instantiated with a flag indicating only unused resources should be considered.
// The function then attempts to find the list of resources that match the filter criteria. If an error occurs during this retrieval process,
// it logs the error and returns an error response. On successful retrieval, it constructs a response object containing the count of resources found
// and the list of resources themselves. It attempts to marshal this response object into JSON. If marshaling fails, it logs the error and returns an error response.
// Otherwise, it logs the marshaled JSON for debugging purposes. Finally, it logs the successful completion of the request and returns the marshaled JSON as a response
// with a 200 OK status.
func GetUnusedResourceList(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetUnusedResourceList"

	query := c.Request.URL.Query()
	resourceGroupIDs := query["resourceGroupID"]

	filter := cmapi_filter_resource.NewResourceUnusedFilter(resourceGroupIDs)
	repository := cmapi_repository_resource.NewResourceListRepository(true)

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
