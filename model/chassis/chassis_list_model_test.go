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
        
package chassis_model

import (
	"reflect"
	"testing"

	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewChassisList(t *testing.T) {
	tests := []struct {
		name string
		want ChassisList
	}{
		{
			"Normal case: Create an instance of the ChassisList struct",
			ChassisList{
				[]Chassis{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChassisList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChassisList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChassisList_ToObject(t *testing.T) {
	type fields struct {
		Chassis []Chassis
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When ChassisList struct has all normal Chassis structs, convert all to map",
			fields{
				[]Chassis{
					{
						map[string]any{"id": "aaa"},
						resource_model.NewResourceList(),
						cxlswitch_model.NewCXLSwitchList(),
					},
					{
						map[string]any{"id": "bbb"},
						resource_model.NewResourceList(),
						cxlswitch_model.NewCXLSwitchList(),
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
			"Normal case: When ChassisList struct includes abnormal Chassis structs (no id or id length is 0), convert to map excluding abnormal Chassis structs",
			fields{
				[]Chassis{
					{
						map[string]any{"aaa": "bbb"},
						resource_model.NewResourceList(),
						cxlswitch_model.NewCXLSwitchList(),
					},
					{
						map[string]any{"id": ""},
						resource_model.NewResourceList(),
						cxlswitch_model.NewCXLSwitchList(),
					},
					{
						map[string]any{"id": "ccc"},
						resource_model.NewResourceList(),
						cxlswitch_model.NewCXLSwitchList(),
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
			cl := &ChassisList{
				Chassis: tt.fields.Chassis,
			}
			if got := cl.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChassisList.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
