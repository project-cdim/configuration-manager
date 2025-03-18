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
        
package cxlswitch_model

import (
	"reflect"
	"testing"

	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewCXLSwitch(t *testing.T) {
	tests := []struct {
		name string
		want CXLSwitch
	}{
		{
			"Normal case: Create an instance of the CXLSwitch struct",
			CXLSwitch{map[string]any{}, resource_model.ResourceList{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCXLSwitch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCXLSwitch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitch_Validate(t *testing.T) {
	type fields struct {
		Properties map[string]any
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"Error case: Return false if Properties does not have an id element",
			fields{
				map[string]any{"aaa": "bbb"},
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if the value of the id element in Properties is empty",
			fields{
				map[string]any{"id": ""},
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Normal case: Return true if Properties has an id element and its value is not empty",
			fields{
				map[string]any{"id": "test"},
				resource_model.NewResourceList(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CXLSwitch{
				Properties: tt.fields.Properties,
				Resources:  tt.fields.Resources,
			}
			if got := c.Validate(); got != tt.want {
				t.Errorf("CXLSwitch.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitch_ToObject(t *testing.T) {
	type fields struct {
		Properties map[string]any
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Successfully convert CXLSwitch struct to map",
			fields{
				map[string]any{"id": "test"},
				resource_model.NewResourceList(),
			},
			map[string]any{
				"id":        "test",
				"resources": []map[string]any{},
			},
		},
		{
			"Error case: Return nil for a Node struct without an id",
			fields{
				map[string]any{"aaa": "bbb"},
				resource_model.NewResourceList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CXLSwitch{
				Properties: tt.fields.Properties,
				Resources:  tt.fields.Resources,
			}
			if got := c.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CXLSwitch.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitch_ToObject4Chassis(t *testing.T) {
	type fields struct {
		Properties map[string]any
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Successfully convert CXLSwitch struct to map",
			fields{
				map[string]any{"id": "test"},
				resource_model.NewResourceList(),
			},
			map[string]any{
				"id": "test",
			},
		},
		{
			"Error case: Return nil for a Node struct without an id",
			fields{
				map[string]any{"aaa": "bbb"},
				resource_model.NewResourceList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CXLSwitch{
				Properties: tt.fields.Properties,
				Resources:  tt.fields.Resources,
			}
			if got := c.ToObject4Chassis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CXLSwitch.ToObject4Chassis() = %v, want %v", got, tt.want)
			}
		})
	}
}
