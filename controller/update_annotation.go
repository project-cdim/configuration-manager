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
	cmapi_model_annotation "github.com/project-cdim/configuration-manager/model/annotation"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"
	cmapi_repository_annotation "github.com/project-cdim/configuration-manager/repository/annotation"
	cmapi_repository_resource "github.com/project-cdim/configuration-manager/repository/resource"

	"github.com/gin-gonic/gin"
)

// UpdateAnnotation handles the update of annotations for a given resource ID.
//
// It retrieves the resource based on the provided ID, updates its annotations with the properties
// provided in the request body, and persists the changes. The function handles different scenarios
// based on the resource type (CPU or non-CPU) and the presence of non-removable devices associated
// with the resource.
//
// When non-removable devices are associated with the resource, annotations will be updated for all
// related resources as a batch operation. For CPU resources, the annotations are updated for both
// the CPU itself and its non-removable devices. For non-CPU resources with associated CPU resources,
// the annotations are updated for both the non-CPU resource and any related resources through the CPU.
//
// Parameters:
//
// c: *gin.Context - The Gin context containing the request and response information.
//   The request context must include the resource ID as a parameter.
//
// Responses:
//
// 200 OK: The annotation was successfully updated. The response body contains the updated resource data.
// 400 Bad Request: The request body could not be unmarshaled, or the provided data is invalid.
//   The response body contains an error message.
// 404 Not Found: The resource with the given ID does not exist. The response body contains an error message.
// 500 Internal Server Error: An error occurred while retrieving or updating the resource in the database.
//   The response body contains an error message.
func UpdateAnnotation(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "UpdateAnnotation"

	id := c.Param("id")
	// Reads the JSON from the RequestBody and expands it into a Map variable
	annotationProperties, err := unmarshalRequestBodyForMap(c)
	if err != nil {
		errorDatial := "unmarshalRequestBodyForMap error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Checks if the resource associated with the target annotation exists by searching once
	filter := cmapi_filter.NewNoFilter()
	getRepository := cmapi_repository_resource.NewResourceRepository(id)
	resource, err := cmapi_repository.RelayFind(&getRepository, filter)
	if err != nil {
		// In case of an error during the graph DB search or array storage process,
		// outputs JSON containing the error code and error message to the ResponseBody and terminates
		errorDatial := "RelayFind error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	if resource == nil {
		// If the target resource for update did not exist
		errorDatial := "The target resource for update did not exist"
		common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, errorDatial, id), false)
		c.JSON(http.StatusNotFound, convertErrorResponse(http.StatusNotFound, errorDatial))
		return
	}

	// Collect the device IDs to be updated.
	deviceIDs := make([]string, 0, 10)
	device := resource["device"].(map[string]any)
	nonRemovableDeviceIDs := getNonRemovableDeviceIDs(device)

	if len(nonRemovableDeviceIDs) == 0 {
		deviceIDs = append(deviceIDs, id)
	} else {
		if device["type"].(string) == CPU {
			deviceIDs = append(deviceIDs, id)
			deviceIDs = append(deviceIDs, nonRemovableDeviceIDs...)
		} else {
			if len(nonRemovableDeviceIDs) > 1 {
				// Since it is assumed that there is only one nonRemovableDevices element, a warning is output if there are two or more.
				common.Log.Warn(fmt.Sprintf("%s %s [id : %v]", funcName, "For resources other than CPU, there were two or more elements in the nonRemovableDevices element.", id), false)
			}

			cpuDeviceID := nonRemovableDeviceIDs[0]
			getRepository = cmapi_repository_resource.NewResourceRepository(cpuDeviceID)
			cpuResource, err := cmapi_repository.RelayFind(&getRepository, filter)
			if err != nil {
				errorDatial := "RelayFind error"
				common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
				c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
				return
			}

			cpuDevice := cpuResource["device"].(map[string]any)
			cpuNonRemovableDeviceIDs := getNonRemovableDeviceIDs(cpuDevice)

			deviceIDs = append(deviceIDs, cpuDeviceID)
			deviceIDs = append(deviceIDs, cpuNonRemovableDeviceIDs...)
		}
	}

	annotation := cmapi_model_annotation.NewAnnotation()
	annotation.Properties = annotationProperties
	repository := cmapi_repository_annotation.NewUpdateAnnotationRepository(deviceIDs)
	res, err := cmapi_repository.RelaySet(&repository, &annotation)
	if err != nil {
		errorDatial := "RelaySet error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	c.JSON(http.StatusOK, res)
}

// getNonRemovableDeviceIDs extracts the deviceIDs of non-removable devices from a device map.
// It expects the device map to contain a "constraints" field, which is a map.
// The "constraints" map should contain a "nonRemovableDevices" field, which is a list of maps.
// Each map in the "nonRemovableDevices" list should contain a "deviceID" field, which is a string.
// If any of these fields are missing or have the wrong type, a warning is logged and an empty list is returned.
// It returns a list of deviceIDs of non-removable devices.
func getNonRemovableDeviceIDs(device map[string]any) []string {
	constraints, ok := device["constraints"].(map[string]any)
	if !ok {
		common.Log.Warn(fmt.Sprintf("constraints field is missing or not a map. resource(%v)", device))
		return []string{}
	}

	nonRemovableDevices, ok := constraints["nonRemovableDevices"].([]any)
	if !ok {
		common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices field is missing or not a list. resource(%v)", device))
		return []string{}
	}

	if len(nonRemovableDevices) == 0 {
		common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices contains no elements. resource(%v)", device))
		return []string{}
	}

	res := []string{}
	for i, nonRemovableDevice := range nonRemovableDevices {
		nonRemovableDeviceMap, ok := nonRemovableDevice.(map[string]any)
		if !ok {
			common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices[%d] is not a map. resource(%v)", i, device))
			return []string{}
		}

		deviceID, ok := nonRemovableDeviceMap["deviceID"].(string)
		if !ok {
			common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices[%d]/deviceID field is missing or not a string. resource(%v)", i, device))
			return []string{}
		}

		res = append(res, deviceID)
	}

	return res
}
