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

	annotation_model "github.com/project-cdim/configuration-manager/model/annotation"
	cxlswitch_model "github.com/project-cdim/configuration-manager/model/cxlswitch"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewChassis(t *testing.T) {
	tests := []struct {
		name string
		want Chassis
	}{
		{
			"Normal case: Create an instance of the Chassis struct",
			Chassis{
				map[string]any{},
				resource_model.ResourceList{},
				cxlswitch_model.CXLSwitchList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChassis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChassis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChassis_Validate(t *testing.T) {
	type fields struct {
		Properties  map[string]any
		Resources   resource_model.ResourceList
		CXLSwitches cxlswitch_model.CXLSwitchList
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
				cxlswitch_model.NewCXLSwitchList(),
			},
			false,
		},
		{
			"Error case: Return false if the value of the id element in Properties is empty",
			fields{
				map[string]any{"id": ""},
				resource_model.NewResourceList(),
				cxlswitch_model.NewCXLSwitchList(),
			},
			false,
		},
		{
			"Normal case: Return true if Properties has an id element and its value is not empty",
			fields{
				map[string]any{"id": "test"},
				resource_model.NewResourceList(),
				cxlswitch_model.NewCXLSwitchList(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chassis{
				Properties:  tt.fields.Properties,
				Resources:   tt.fields.Resources,
				CXLSwitches: tt.fields.CXLSwitches,
			}
			if got := c.Validate(); got != tt.want {
				t.Errorf("Chassis.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChassis_ToObject(t *testing.T) {
	type fields struct {
		Properties  map[string]any
		Resources   resource_model.ResourceList
		CXLSwitches cxlswitch_model.CXLSwitchList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Successfully convert Chassis struct to map",
			fields{
				map[string]any{"id": "test"},
				createResourceList(),
				createCXLSwitchList(),
			},
			map[string]any{
				"id":        "test",
				"resources": createResources(),
			},
		},
		{
			"Error case: Return nil for a Chassis struct without an id",
			fields{
				map[string]any{"aaa": "bbb"},
				createResourceList(),
				createCXLSwitchList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chassis{
				Properties:  tt.fields.Properties,
				Resources:   tt.fields.Resources,
				CXLSwitches: tt.fields.CXLSwitches,
			}
			if got := c.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chassis.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createResourceList() resource_model.ResourceList {
	resourcelist := resource_model.NewResourceList()
	resource := resource_model.Resource{
		Device:           map[string]any{"deviceID": "aaa"},
		Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
		ResourceGroupIDs: []string{"g-abb", "g-acc"},
		NodeIDs:          []string{"n-add", "n-aee"},
		Detected:         false,
	}
	resourcelist.Resources = append(resourcelist.Resources, resource)
	resource = resource_model.Resource{
		Device:           map[string]any{"deviceID": "bbb"},
		Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": false}},
		ResourceGroupIDs: []string{"g-bbb", "g-bcc"},
		NodeIDs:          []string{"n-bdd", "n-bee"},
		Detected:         true,
	}
	resourcelist.Resources = append(resourcelist.Resources, resource)
	return resourcelist
}

func createCXLSwitchList() cxlswitch_model.CXLSwitchList {
	cxlSwitchList := cxlswitch_model.NewCXLSwitchList()
	cxlSwitch1 := cxlswitch_model.CXLSwitch{
		Properties: map[string]any{"id": "a-bbb"},
		Resources:  resource_model.ResourceList{},
	}
	cxlSwitchList.CXLSwitches = append(cxlSwitchList.CXLSwitches, cxlSwitch1)
	cxlSwitch2 := cxlswitch_model.CXLSwitch{
		Properties: map[string]any{"id": "b-bbb"},
		Resources:  resource_model.ResourceList{},
	}
	cxlSwitchList.CXLSwitches = append(cxlSwitchList.CXLSwitches, cxlSwitch2)
	return cxlSwitchList
}

func createResources() []map[string]any {
	return []map[string]any{
		{
			"id": "a-bbb",
		},
		{
			"id": "b-bbb",
		},
		{
			"device":           map[string]any{"deviceID": "aaa"},
			"annotation":       map[string]any{"available": true},
			"resourceGroupIDs": []string{"g-abb", "g-acc"},
			"nodeIDs":          []string{"n-add", "n-aee"},
			"detected":         false,
		},
		{
			"device":           map[string]any{"deviceID": "bbb"},
			"annotation":       map[string]any{"available": false},
			"resourceGroupIDs": []string{"g-bbb", "g-bcc"},
			"nodeIDs":          []string{"n-bdd", "n-bee"},
			"detected":         true,
		},
	}
}
