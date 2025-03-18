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
        
package group_model

import (
	"reflect"
	"testing"

	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewGroupList(t *testing.T) {
	tests := []struct {
		name string
		want GroupList
	}{
		{
			"Normal case: Create an instance of the GroupList struct",
			GroupList{
				[]Group{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroupList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupList_ToObject(t *testing.T) {
	type fields struct {
		Groups []Group
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When GroupList struct has all valid Node structs, convert all to map",
			fields{
				[]Group{
					{
						Id:         "group01",
						Properties: map[string]any{"name": "group01", "description": "group01"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
					{
						Id:         "group02",
						Properties: map[string]any{"name": "group02", "description": "group02"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":          "group01",
					"name":        "group01",
					"description": "group01",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
				},
				{
					"id":          "group02",
					"name":        "group02",
					"description": "group02",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
				},
			},
		},
		{
			"Normal case: When GroupList struct includes invalid Group structs, convert excluding the invalid Node structs to map",
			fields{
				[]Group{
					{
						Id:         "group01",
						Properties: map[string]any{},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
					{
						Id:         "group02",
						Properties: map[string]any{"name": "group02", "description": "group02"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":          "group02",
					"name":        "group02",
					"description": "group02",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gl := &GroupList{
				Groups: tt.fields.Groups,
			}
			if got := gl.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupList.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupList_ToObjectWithResources(t *testing.T) {
	type fields struct {
		Groups []Group
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When GroupList struct has all valid Node structs, convert all to map",
			fields{
				[]Group{
					{
						Id:         "group01",
						Properties: map[string]any{"name": "group01", "description": "group01"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
					{
						Id:         "group02",
						Properties: map[string]any{"name": "group02", "description": "group02"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":          "group01",
					"name":        "group01",
					"description": "group01",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
					"resources":   []map[string]any{},
				},
				{
					"id":          "group02",
					"name":        "group02",
					"description": "group02",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
					"resources":   []map[string]any{},
				},
			},
		},
		{
			"Normal case: When GroupList struct includes invalid Group structs, convert excluding the invalid Node structs to map",
			fields{
				[]Group{
					{
						Id:         "group01",
						Properties: map[string]any{},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
					{
						Id:         "group02",
						Properties: map[string]any{"name": "group02", "description": "group02"},
						CreatedAt:  "2021-01-01T00:00:00Z",
						UpdatedAt:  "2021-01-01T00:00:00Z",
						Resources:  resource_model.NewResourceList(),
					},
				},
			},
			[]map[string]any{
				{
					"id":          "group02",
					"name":        "group02",
					"description": "group02",
					"createdAt":   "2021-01-01T00:00:00Z",
					"updatedAt":   "2021-01-01T00:00:00Z",
					"resources":   []map[string]any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gl := &GroupList{
				Groups: tt.fields.Groups,
			}
			if got := gl.ToObjectWithResources(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupList.ToObjectWithResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
