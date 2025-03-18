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
	cmapi_repository_node "github.com/project-cdim/configuration-manager/repository/node"

	"github.com/gin-gonic/gin"
)

// GetNodeList retrieves a list of all nodes and the list of resources associated with each node.
// It logs the start of the request, attempts to retrieve the node list using a no-filter approach,
// and handles errors by logging and returning an error response. On success, it marshals the response
// into JSON and sends it back to the client with an HTTP status code of 200 OK.
func GetNodeList(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "GetNodeList"

	filter := cmapi_filter.NewNoFilter()
	repository := cmapi_repository_node.NewNodeListRepository()
	nodes, err := cmapi_repository.RelayFindList(&repository, filter)
	if err != nil {
		errorDatial := "RelayFindList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	res := gin.H{
		"count": len(nodes),
		"nodes": nodes,
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}
