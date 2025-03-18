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
        
package main

import (
	"testing"

	"github.com/project-cdim/configuration-manager/common"
)

const (
	urlBase string = "/" + common.ProjectName + "/api/v1"
)

func Test_main(t *testing.T) {
	t.Skip("not test")
}

// func Test_Req_GetResource(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response []int
// 	}{
// 		{
// 			"Normal case and API call check: Resource exists or not (existence depends on the state of the database)",
// 			urlBase + "/resources/res101",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if !slices.Contains(tt.response, resRecorder.Code) {
// 				t.Errorf("Test_Req_GetResource() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetResourceList(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response int
// 	}{
// 		{
// 			"Normal case and API call check: Specify true for the detail query parameter",
// 			urlBase + "/resources?detail=true",
// 			"GET",
// 			http.StatusOK,
// 		},
// 		{
// 			"Normal case: Specify false for the detail query parameter",
// 			urlBase + "/resources?detail=false",
// 			"GET",
// 			http.StatusOK,
// 		},
// 		{
// 			"Normal case: Specify TRUE for the detail query parameter",
// 			urlBase + "/resources?detail=TRUE",
// 			"GET",
// 			http.StatusOK,
// 		},
// 		{
// 			"Normal case: Specify FALSE for the detail query parameter",
// 			urlBase + "/resources?detail=FALSE",
// 			"GET",
// 			http.StatusOK,
// 		},
// 		{
// 			"Normal case: detail query parameter not specified",
// 			urlBase + "/resources",
// 			"GET",
// 			http.StatusOK,
// 		},
// 		{
// 			"Abnormal case: Specify an invalid value (other than true/false) for the detail query parameter",
// 			urlBase + "/resources?detail=aaa",
// 			"GET",
// 			http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if tt.response != resRecorder.Code {
// 				t.Errorf("Test_Req_GetResourceList() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetNode(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response []int
// 	}{
// 		{
// 			"Normal case and API call check: Node exists or not (existence depends on the state of the database)",
// 			urlBase + "/nodes/res101",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if !slices.Contains(tt.response, resRecorder.Code) {
// 				t.Errorf("Test_Req_GetNode() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetNodeList(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response int
// 	}{
// 		{
// 			"Normal case and API call check",
// 			urlBase + "/nodes",
// 			"GET",
// 			http.StatusOK,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if tt.response != resRecorder.Code {
// 				t.Errorf("Test_Req_GetNodeList() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetCxlSwitch(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response []int
// 	}{
// 		{
// 			"Normal case and API call check: Switch exists or not (existence depends on the state of the database)",
// 			urlBase + "/cxlswitches/CXL11",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if !slices.Contains(tt.response, resRecorder.Code) {
// 				t.Errorf("Test_Req_GetCxlSwitch() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetCxlSwitchList(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response int
// 	}{
// 		{
// 			"Normal case and API call check",
// 			urlBase + "/cxlswitches",
// 			"GET",
// 			http.StatusOK,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if tt.response != resRecorder.Code {
// 				t.Errorf("Test_Req_GetCxlSwitchList() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_RegisterDevice(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		body     string
// 		response int
// 	}{
// 		{
// 			"Abnormal case: Specify an unexpected type in the request body (other than []map[string]any)",
// 			urlBase + "/devices",
// 			"POST",
// 			"1",
// 			http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			var reqBody io.Reader = nil
// 			if len(tt.body) > 0 {
// 				reqBody = strings.NewReader(tt.body)
// 			}
// 			req, _ := http.NewRequest(tt.method, tt.uri, reqBody)
// 			router.ServeHTTP(resRecorder, req)
// 			if tt.response != resRecorder.Code {
// 				t.Errorf("Test_Req_RegisterDevice() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetRack(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response []int
// 	}{
// 		{
// 			"Normal case and API call check: Specify true for the detail query parameter. Rack exists or not (existence depends on the state of the database)",
// 			urlBase + "/racks/rack11?detail=true",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 		{
// 			"Normal case: Specify false for the detail query parameter. Rack exists or not (existence depends on the state of the database)",
// 			urlBase + "/racks/rack11?detail=false",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 		{
// 			"Normal case: Specify TRUE for the detail query parameter. Rack exists or not (existence depends on the state of the database)",
// 			urlBase + "/racks/rack11?detail=TRUE",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 		{
// 			"Normal case: Specify FALSE for the detail query parameter. Rack exists or not (existence depends on the state of the database)",
// 			urlBase + "/racks/rack11?detail=FALSE",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 		{
// 			"Normal case: detail query parameter not specified. Rack exists or not (existence depends on the state of the database)",
// 			urlBase + "/racks/rack11",
// 			"GET",
// 			[]int{http.StatusOK, http.StatusNotFound},
// 		},
// 		{
// 			"Abnormal case: Specify an invalid value (other than true/false) for the detail query parameter",
// 			urlBase + "/racks/rack11?detail=aaa",
// 			"GET",
// 			[]int{http.StatusBadRequest},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if !slices.Contains(tt.response, resRecorder.Code) {
// 				t.Errorf("Test_Req_GetRack() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }

// func Test_Req_GetAvailableResourceList(t *testing.T) {
// 	router := setupEngine()

// 	tests := []struct {
// 		name     string
// 		uri      string
// 		method   string
// 		response []int
// 	}{
// 		{
// 			"Normal case and API call check",
// 			urlBase + "/resources/available",
// 			"GET",
// 			[]int{http.StatusOK},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resRecorder := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tt.method, tt.uri, nil)
// 			router.ServeHTTP(resRecorder, req)
// 			if !slices.Contains(tt.response, resRecorder.Code) {
// 				t.Errorf("Test_Req_GetAvailableResourceList() = %v, want %v", resRecorder.Code, tt.response)
// 			}
// 		})
// 	}
// }
