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
	cmapi_repository_rack "github.com/project-cdim/configuration-manager/repository/rack"

	"github.com/gin-gonic/gin"
)

// GetRack retrieves a specific rack by its ID and the list of chassis associated with the rack,
// as well as devices (switches, resources) associated with each chassis. It starts by logging the start
// of the request. It then retrieves the 'detail' query parameter to determine the level of detail required
// in the response. Using a no-filter approach, it attempts to retrieve the rack from the repository.
// If an error occurs during retrieval, it logs the error and returns an error response. If the rack is not found,
// it logs a warning and returns a 404 Not Found response. On successful retrieval, it attempts to marshal the rack
// into JSON. If marshaling fails, it logs the error and returns an error response. Otherwise, it logs the marshaled
// JSON for debugging, logs the successful completion of the request, and returns the rack as a JSON response with a 200 OK status.
func GetRack(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetRack"

	id := c.Param("id")
	// Retrieve query parameter: detail
	detail, err := getBoolQueryParam(c, "detail")
	if err != nil {
		errorDatial := "getBoolQueryParam error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_rack.NewRackRepository(id, detail)
	res, err := cmapi_repository.RelayFind(&repository, filter)
	if err != nil {
		// Outputs JSON containing the error code and error message to the ResponseBody and terminates
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
