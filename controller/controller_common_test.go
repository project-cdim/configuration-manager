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
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHwResourceType_convertToDBLabel(t *testing.T) {
	tests := []struct {
		name    string
		rt      hwResourceType
		want    string
		wantErr bool
	}{
		{
			"Normal case: Convert CPU",
			hwResourceType("CPU"),
			"CPU",
			false,
		},
		{
			"Normal case: Convert Accelerator",
			hwResourceType("Accelerator"),
			"Accelerator",
			false,
		},
		{
			"Normal case: Convert DSP",
			hwResourceType("DSP"),
			"DSP",
			false,
		},
		{
			"Normal case: Convert FPGA",
			hwResourceType("FPGA"),
			"FPGA",
			false,
		},
		{
			"Normal case: Convert GPU",
			hwResourceType("GPU"),
			"GPU",
			false,
		},
		{
			"Normal case: Convert UnknownProcessor",
			hwResourceType("UnknownProcessor"),
			"UnknownProcessor",
			false,
		},
		{
			"Normal case: Convert memory",
			hwResourceType("memory"),
			"Memory",
			false,
		},
		{
			"Normal case: Convert storage",
			hwResourceType("storage"),
			"Storage",
			false,
		},
		{
			"Normal case: Convert networkInterface",
			hwResourceType("networkInterface"),
			"NetworkInterface",
			false,
		},
		{
			"Normal case: Convert graphicController",
			hwResourceType("graphicController"),
			"GraphicController",
			false,
		},
		{
			"Normal case: Convert virtualMedia",
			hwResourceType("virtualMedia"),
			"VirtualMedia",
			false,
		},
		{
			"Error case: Unexpected type",
			hwResourceType("aaa"),
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rt.convertToDBLabel()
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceType.convertToDBLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResourceType.convertToDBLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshalRequestBodyForMap(t *testing.T) {
	t.Skip("not test")
}

func Test_unmarshalRequestBodyForSlice(t *testing.T) {
	t.Skip("not test")
}

func Test_getBoolQueryParam(t *testing.T) {
	type args struct {
		c    *gin.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test true value",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test false value",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test True value",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test fAlse value",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param4",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test invalid value",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param5",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "test not exists element",
			args: args{
				c:    setupTestGinContext("param1=true&param2=false&param3=True&param4=fAlse&param5=aaa"),
				name: "param6",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBoolQueryParam(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBoolQueryParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBoolQueryParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		details []string
		want    gin.H
	}{
		{
			name:    "Internal Server Error without details",
			status:  http.StatusInternalServerError,
			details: nil,
			want:    gin.H{"code": "internalServerError", "message": "Internal Server Error. Contact the administrator."},
		},
		{
			name:    "Bad Request without details",
			status:  http.StatusBadRequest,
			details: nil,
			want:    gin.H{"code": "badRequest", "message": "Bad Request. Check the request parameters."},
		},
		{
			name:    "Not Found without details",
			status:  http.StatusNotFound,
			details: nil,
			want:    gin.H{"code": "notFound", "message": "Not Found. Check the request URL."},
		},
		{
			name:    "Internal Server Error with details",
			status:  http.StatusInternalServerError,
			details: []string{"Something went wrong"},
			want:    gin.H{"code": "internalServerError", "message": "Internal Server Error. Contact the administrator.", "details": "Something went wrong"},
		},
		{
			name:    "Bad Request with details",
			status:  http.StatusBadRequest,
			details: []string{"Invalid input"},
			want:    gin.H{"code": "badRequest", "message": "Bad Request. Check the request parameters.", "details": "Invalid input"},
		},
		{
			name:    "Not Found with details",
			status:  http.StatusNotFound,
			details: []string{"Resource not found"},
			want:    gin.H{"code": "notFound", "message": "Not Found. Check the request URL.", "details": "Resource not found"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertErrorResponse(tt.status, tt.details...)
			for key, value := range tt.want {
				if got[key] != value {
					t.Errorf("convertErrorResponse() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_logResponseBody(t *testing.T) {
	t.Skip("not test")
}

// setupTestGinContext creates a new gin.Context for testing purposes.
// It sets the gin mode to TestMode, creates a new HTTP GET request with the provided query string,
// and returns the initialized gin.Context.
//
// Parameters:
//
//	queryString - The query string to be appended to the test URL.
//
// Returns:
//
//	*gin.Context - The initialized gin.Context with the test request.
func setupTestGinContext(queryString string) *gin.Context {
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/test?%s", queryString), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c
}
