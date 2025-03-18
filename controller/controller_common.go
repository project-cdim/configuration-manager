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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	cmapi_filter "github.com/project-cdim/configuration-manager/filter"

	"github.com/gin-gonic/gin"
)

// Configuration item names for hardware control features
const (
	CPU               = "CPU"
	Accelerator       = "Accelerator"
	DSP               = "DSP"
	FPGA              = "FPGA"
	GPU               = "GPU"
	UnknownProcessor  = "UnknownProcessor"
	Memory            = "memory"
	Storage           = "storage"
	NetworkInterface  = "networkInterface"
	GraphicController = "graphicController"
	VirtualMedia      = "virtualMedia"
)

// Label names for Vertex
const (
	DB_CPU               = "CPU"
	DB_Accelerator       = "Accelerator"
	DB_DSP               = "DSP"
	DB_FPGA              = "FPGA"
	DB_GPU               = "GPU"
	DB_UnknownProcessor  = "UnknownProcessor"
	DB_Memory            = "Memory"
	DB_Storage           = "Storage"
	DB_NetworkInterface  = "NetworkInterface"
	DB_GraphicController = "GraphicController"
	DB_VirtualMedia      = "VirtualMedia"
)

// StatusToResponse is a map that maps the status code to the response body.
var StatusToResponse = map[int]gin.H{
	http.StatusInternalServerError: {"code": "internalServerError", "message": "Internal Server Error. Contact the administrator."},
	http.StatusBadRequest:          {"code": "badRequest", "message": "Bad Request. Check the request parameters."},
	http.StatusNotFound:            {"code": "notFound", "message": "Not Found. Check the request URL."},
}

// hwResourceType defines a string type for representing various hardware resource categories.
type hwResourceType string

// convertToDBLabel converts a JSON label to its corresponding database label.
// It takes a hwResourceType and returns the matching database label string and nil error for known types.
// For unknown types, it returns an empty string and an error indicating an unexpected type.
func (rt hwResourceType) convertToDBLabel() (string, error) {
	switch rt {
	case CPU:
		return DB_CPU, nil
	case Accelerator:
		return DB_Accelerator, nil
	case DSP:
		return DB_DSP, nil
	case FPGA:
		return DB_FPGA, nil
	case GPU:
		return DB_GPU, nil
	case UnknownProcessor:
		return DB_UnknownProcessor, nil
	case Memory:
		return DB_Memory, nil
	case Storage:
		return DB_Storage, nil
	case NetworkInterface:
		return DB_NetworkInterface, nil
	case GraphicController:
		return DB_GraphicController, nil
	case VirtualMedia:
		return DB_VirtualMedia, nil
	default:
		return "", fmt.Errorf("unexpected type in JSON. type(%v)", rt)
	}
}

// unmarshalRequestBodyForMap reads the JSON from the request body, unmarshals it into a map[string]any, and returns the map.
// It returns an error if reading the request body fails or if the JSON is not in the correct format.
func unmarshalRequestBodyForMap(c *gin.Context) (map[string]any, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	return res, nil
}

// unmarshalRequestBodyForSlice reads the JSON from the request body, unmarshals it into a slice of map[string]any, and returns the slice.
// It returns an error if reading the request body fails, or if the JSON is not in the correct format.
func unmarshalRequestBodyForSlice(c *gin.Context) ([]map[string]any, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	var res []map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	return res, nil
}

// getBoolQueryParam retrieves a boolean query parameter from the given gin.Context.
// It checks if the query parameter value matches "true" or "false" (case insensitive).
// If the value matches "true", it returns true. If the value matches "false", it returns false.
// If the value does not match either, it returns an error.
//
// Parameters:
//   - c: *gin.Context - The gin context from which to retrieve the query parameter.
//   - name: string - The name of the query parameter to retrieve.
//
// Returns:
//   - bool: The boolean value of the query parameter.
//   - error: An error if the query parameter value is not "true" or "false".
func getBoolQueryParam(c *gin.Context, name string) (bool, error) {
	v := c.DefaultQuery(name, "false")
	if strings.EqualFold(v, cmapi_filter.TrueStr) {
		return true, nil
	} else if strings.EqualFold(v, cmapi_filter.FalseStr) {
		return false, nil
	}

	return false, fmt.Errorf("query parameter value error. name(%v) value(%v)", name, v)
}

// convertErrorResponse converts an error response containing the specified status and details.
// It retrieves the response map corresponding to the status and adds the details to the map if available.
// It returns the converted response map.
func convertErrorResponse(status int, details ...string) gin.H {
	res := StatusToResponse[status]
	if len(details) > 0 {
		res["details"] = details[0]
	}
	return res
}

// logResponseBody logs the response body at the debug level using the common.LoggerApp.
// The response body is formatted as a string and included in the log message.
//
// Parameters:
//
//	res (any): The response body to be logged. It can be of any type.
func logResponseBody(res any) {
	common.Log.Debug(fmt.Sprintf("Response Body : %v", res))
}
