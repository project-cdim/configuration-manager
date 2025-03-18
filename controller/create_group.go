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
	cmapi_model_group "github.com/project-cdim/configuration-manager/model/group"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_group "github.com/project-cdim/configuration-manager/repository/group"

	"github.com/gin-gonic/gin"
)

// CreateGroup is a handler function to create a new group.
// It converts the request body to a map, performs validation,
// creates the group, and saves it to the database.
// On success, it returns the created group object.
//
// Parameters:
//   - c: gin.Context - Request context
//
// Response:
//   - On success: HTTP status 201 (Created) and the created group object
//   - On validation error: HTTP status 400 (Bad Request)
//   - On server error: HTTP status 500 (Internal Server Error)
func CreateGroup(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "CreateGroup"

	properties, err := unmarshalRequestBodyForMap(c)
	if err != nil {
		errorDatial := "unmarshalRequestBodyForMap error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Validation of the requestBody
	if !cmapi_model_group.ValidateProperty(properties) {
		errorDatial := "Validation error"
		common.Log.Error(fmt.Sprintf("%s %s", funcName, errorDatial), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	group := cmapi_model_group.NewGroupWithCreateTimeStampsNow(properties)
	repository := cmapi_repository_group.NewCreateGroupRepository()
	res, err := cmapi_repository.RelaySet(&repository, &group)
	if err != nil {
		errorDatial := "RelaySet error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusCreated, res)
}
