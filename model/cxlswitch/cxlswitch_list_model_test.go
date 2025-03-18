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

func TestNewCXLSwitchList(t *testing.T) {
	tests := []struct {
		name string
		want CXLSwitchList
	}{
		{
			"Normal case: Create an instance of the CXLSwitchList struct",
			CXLSwitchList{
				[]CXLSwitch{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCXLSwitchList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCXLSwitchList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitchList_ToObject(t *testing.T) {
	type fields struct {
		CXLSwitches []CXLSwitch
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When the CXLSwitchList struct has all normal CXLSwitch structs, convert all to map",
			fields{
				[]CXLSwitch{
					{
						map[string]any{"id": "aaa"},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": "bbb"},
						resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":        "aaa",
					"resources": []map[string]any{},
				},
				{
					"id":        "bbb",
					"resources": []map[string]any{},
				},
			},
		},
		{
			"Normal case: When the CXLSwitchList struct includes abnormal (no id or id length is 0) CXLSwitch structs, convert to map excluding the abnormal CXLSwitch structs",
			fields{
				[]CXLSwitch{
					{
						map[string]any{"aaa": "bbb"},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": ""},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": "ccc"},
						resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":        "ccc",
					"resources": []map[string]any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &CXLSwitchList{
				CXLSwitches: tt.fields.CXLSwitches,
			}
			if got := nl.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CXLSwitchList.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitchList_ToObject4Chassis(t *testing.T) {
	type fields struct {
		CXLSwitches []CXLSwitch
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When the CXLSwitchList struct has all normal CXLSwitch structs, convert all to map",
			fields{
				[]CXLSwitch{
					{
						map[string]any{"id": "aaa"},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": "bbb"},
						resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id": "aaa",
				},
				{
					"id": "bbb",
				},
			},
		},
		{
			"Normal case: When the CXLSwitchList struct includes abnormal (no id or id length is 0) CXLSwitch structs, convert to map excluding the abnormal CXLSwitch structs",
			fields{
				[]CXLSwitch{
					{
						map[string]any{"aaa": "bbb"},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": ""},
						resource_model.NewResourceList(),
					},
					{
						map[string]any{"id": "ccc"},
						resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id": "ccc",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &CXLSwitchList{
				CXLSwitches: tt.fields.CXLSwitches,
			}
			if got := nl.ToObject4Chassis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CXLSwitchList.ToObject4Chassis() = %v, want %v", got, tt.want)
			}
		})
	}
}
