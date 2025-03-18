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
	cmapi_repository_cxlswitch "github.com/project-cdim/configuration-manager/repository/cxlswitch"

	"github.com/gin-gonic/gin"
)

// GetCxlSwitch retrieves a specific CXL switch and the list of resources associated with each CXL switch.
// It logs the start and end of the request, handles any errors by logging them and returning an appropriate JSON response.
// If the CXL switch is not found, it returns a 404 status with a message.
// On success, it marshals the CXL switch into JSON and returns it in the response body.
func GetCxlSwitch(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetCxlSwitch"

	id := c.Param("id")
	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_cxlswitch.NewCXLSwitchRepository(id)
	res, err := cmapi_repository.RelayFind(&repository, filter)
	if err != nil {
		// Outputs JSON containing the error code and error message to the ResponseBody and terminates.
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
