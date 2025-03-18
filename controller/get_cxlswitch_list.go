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

// GetCxlSwitchList retrieves a list of all CXL switches and the list of resources associated with each node.
// It logs the start and end of the request, handles any errors by logging them and returning an appropriate JSON response.
// On success, it marshals the list of CXL switches into JSON and returns it in the response body.
func GetCxlSwitchList(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetCxlSwitchList"

	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_cxlswitch.NewCXLSwitchListRepository()
	cxlswitches, err := cmapi_repository.RelayFindList(&repository, filter)
	if err != nil {
		// Outputs JSON containing the error code and error message to the ResponseBody and terminates.
		errorDatial := "RelayFindList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	res := gin.H{
		"count":       len(cxlswitches),
		"CXLSwitches": cxlswitches,
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	// Sets the return value
	c.JSON(http.StatusOK, res)
}
