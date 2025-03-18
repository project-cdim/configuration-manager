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
        
package rack_model

import (
	"reflect"
	"testing"

	chassis_model "github.com/project-cdim/configuration-manager/model/chassis"
)

func TestNewRack(t *testing.T) {
	tests := []struct {
		name string
		want Rack
	}{
		{
			"Normal case: Create an instance of the Rack struct",
			Rack{map[string]any{}, chassis_model.ChassisList{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRack_Validate(t *testing.T) {
	type fields struct {
		Properties map[string]any
		Chassis    chassis_model.ChassisList
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
				chassis_model.NewChassisList(),
			},
			false,
		},
		{
			"Error case: Return false if the value of the id element in Properties is empty",
			fields{
				map[string]any{"id": ""},
				chassis_model.NewChassisList(),
			},
			false,
		},
		{
			"Normal case: Return true if Properties has an id element and its value is not empty",
			fields{
				map[string]any{"id": "test"},
				chassis_model.NewChassisList(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rack{
				Properties: tt.fields.Properties,
				Chassis:    tt.fields.Chassis,
			}
			if got := r.Validate(); got != tt.want {
				t.Errorf("Rack.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRack_ToObject(t *testing.T) {
	type fields struct {
		Properties map[string]any
		Chassis    chassis_model.ChassisList
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal case: Successfully convert Rack struct to map",
			fields{
				map[string]any{"id": "test"},
				chassis_model.NewChassisList(),
			},
			map[string]any{
				"id":      "test",
				"chassis": []map[string]any{},
			},
		},
		{
			"Error case: Return nil for a Rack struct without an id",
			fields{
				map[string]any{"aaa": "bbb"},
				chassis_model.NewChassisList(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rack{
				Properties: tt.fields.Properties,
				Chassis:    tt.fields.Chassis,
			}
			if got := r.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rack.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
