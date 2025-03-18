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

	"github.com/project-cdim/configuration-manager/model"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name string
		want Group
	}{
		{
			"Normal case: Create an instance of the Group struct",
			Group{
				Id:         "",
				Properties: map[string]any{},
				CreatedAt:  "",
				UpdatedAt:  "",
				Resources:  resource_model.NewResourceList(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGroupWithCreateTimeStampsNow(t *testing.T) {
	type args struct {
		properties map[string]any
	}
	tests := []struct {
		name string
		args args
		want Group
	}{
		{
			"Normal case: Create an instance of the Group struct",
			args{
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
			},
			Group{
				Id: "",
				Properties: map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				CreatedAt: "",
				UpdatedAt: "",
				Resources: resource_model.NewResourceList(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGroupWithCreateTimeStampsNow(tt.args.properties)
			if got.Id != "" {
				t.Errorf("Group.NewGroupWithCreateTimeStampsNow() Id = %v", got.Id)
			}
			if !reflect.DeepEqual(got.Properties, tt.want.Properties) {
				t.Errorf("Group.NewGroupWithCreateTimeStampsNow() Properties = %v, want %v", got.Properties, tt.want.Properties)
			}
			if !model.ValidateISO8601(got.CreatedAt) {
				t.Errorf("Group.NewGroupWithCreateTimeStampsNow() CreatedAt = %v", got.CreatedAt)
			}
			if !model.ValidateISO8601(got.UpdatedAt) {
				t.Errorf("Group.NewGroupWithCreateTimeStampsNow() UpdatedAt = %v", got.UpdatedAt)
			}
			if !reflect.DeepEqual(got.Resources, tt.want.Resources) {
				t.Errorf("Group.NewGroupWithCreateTimeStampsNow() Resources = %v", got.Resources)
			}
		})
	}
}

func TestNewGroupForUpdate(t *testing.T) {
	type args struct {
		groupFromDb map[string]any
		properties  map[string]any
	}
	tests := []struct {
		name string
		args args
		want Group
	}{
		{
			"Normal case: Create an instance of the Group struct",
			args{
				map[string]any{
					"id":        "test ID",
					"createdAt": "2021-01-01T00:00:00Z",
				},
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
			},
			Group{
				Id: "test ID",
				Properties: map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				CreatedAt: "2021-01-01T00:00:00Z",
				UpdatedAt: "",
				Resources: resource_model.NewResourceList(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGroupForUpdate(tt.args.groupFromDb, tt.args.properties)
			if got.Id != tt.want.Id {
				t.Errorf("Group.NewGroupForUpdate() Id = %v, want %v", got.Id, tt.want.Id)
			}
			if !reflect.DeepEqual(got.Properties, tt.want.Properties) {
				t.Errorf("Group.NewGroupForUpdate() Properties = %v, want %v", got.Properties, tt.want.Properties)
			}
			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("Group.NewGroupForUpdate() CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
			if !model.ValidateISO8601(got.UpdatedAt) {
				t.Errorf("Group.NewGroupForUpdate() UpdatedAt = %v", got.UpdatedAt)
			}
			if !reflect.DeepEqual(got.Resources, tt.want.Resources) {
				t.Errorf("Group.NewGroupForUpdate() Resources = %v", got.Resources)
			}
		})
	}
}

func TestGroup_SetCreatedAt(t *testing.T) {
	type fields struct {
		Id         string
		Properties map[string]any
		CreatedAt  string
		UpdatedAt  string
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"Normal case: Set CreatedAt",
			fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Id:         tt.fields.Id,
				Properties: tt.fields.Properties,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				Resources:  tt.fields.Resources,
			}
			g.createTimeStampsNow()
			if !model.ValidateISO8601(g.CreatedAt) {
				t.Errorf("Group.SetUpdatedAt() CreatedAt = %v", g.CreatedAt)
			}
			if !model.ValidateISO8601(g.UpdatedAt) {
				t.Errorf("Group.SetUpdatedAt() UpdatedAt = %v", g.UpdatedAt)
			}
		})
	}
}

func TestGroup_SetUpdatedAt(t *testing.T) {
	type fields struct {
		Id         string
		Properties map[string]any
		CreatedAt  string
		UpdatedAt  string
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"Normal case: Set UpdateAt",
			fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Id:         tt.fields.Id,
				Properties: tt.fields.Properties,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				Resources:  tt.fields.Resources,
			}
			g.updateTimeStampsNow()
			if model.ValidateISO8601(g.CreatedAt) {
				t.Errorf("Group.SetUpdatedAt() CreatedAt = %v", g.CreatedAt)
			}
			if !model.ValidateISO8601(g.UpdatedAt) {
				t.Errorf("Group.SetUpdatedAt() UpdatedAt = %v", g.UpdatedAt)
			}
		})
	}
}

func TestGroup_Validate(t *testing.T) {
	type fields struct {
		Id         string
		Properties map[string]any
		CreatedAt  string
		UpdatedAt  string
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"Error case: Return false if Properties name field is not exist",
			fields{
				"group01",
				map[string]any{
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties name field is not string",
			fields{
				"group01",
				map[string]any{
					"name":        123,
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties name field length is 0",
			fields{
				"group01",
				map[string]any{
					"name":        "",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties name field length is over 64",
			fields{
				"group01",
				map[string]any{
					"name":        "12345678901234567890123456789012345678901234567890123456789012345",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties description field is not exist",
			fields{
				"group01",
				map[string]any{
					"name": "group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties description field is not string",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": 123,
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if Properties description field length is over 256",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if CreatedAt is not in ISO8601 format",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Error case: Return false if UpdatedAt is not in ISO8601 format",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00",
				resource_model.NewResourceList(),
			},
			false,
		},
		{
			"Normal case: Return true if all validations pass",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Id:         tt.fields.Id,
				Properties: tt.fields.Properties,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				Resources:  tt.fields.Resources,
			}
			if got := g.Validate(); got != tt.want {
				t.Errorf("Group.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_ToObject(t *testing.T) {
	type fields struct {
		Id         string
		Properties map[string]any
		CreatedAt  string
		UpdatedAt  string
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Return a map representation of the Group struct",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			map[string]any{
				"id":          "group01",
				"name":        "group01",
				"description": "This is group01",
				"createdAt":   "2021-01-01T00:00:00Z",
				"updatedAt":   "2021-01-01T00:00:00Z",
			},
		},
		{
			"Error case: Return nil if the not valid",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00",
				resource_model.NewResourceList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Id:         tt.fields.Id,
				Properties: tt.fields.Properties,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				Resources:  tt.fields.Resources,
			}
			if got := g.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_ToObjectWithResources(t *testing.T) {
	type fields struct {
		Id         string
		Properties map[string]any
		CreatedAt  string
		UpdatedAt  string
		Resources  resource_model.ResourceList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Return a map representation of the Group struct",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00Z",
				resource_model.NewResourceList(),
			},
			map[string]any{
				"id":          "group01",
				"name":        "group01",
				"description": "This is group01",
				"createdAt":   "2021-01-01T00:00:00Z",
				"updatedAt":   "2021-01-01T00:00:00Z",
				"resources":   []map[string]any{},
			},
		},
		{
			"Error case: Return nil if the not valid",
			fields{
				"group01",
				map[string]any{
					"name":        "group01",
					"description": "This is group01",
				},
				"2021-01-01T00:00:00Z",
				"2021-01-01T00:00:00",
				resource_model.NewResourceList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Id:         tt.fields.Id,
				Properties: tt.fields.Properties,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				Resources:  tt.fields.Resources,
			}
			if got := g.ToObjectWithResources(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.ToObjectWithResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
