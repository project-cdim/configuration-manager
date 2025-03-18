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
        
package resource_model

import (
	"reflect"
	"testing"

	annotation_model "github.com/project-cdim/configuration-manager/model/annotation"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		name string
		want Resource
	}{
		{
			"Normal Case: Generates an instance of the Resource struct",
			Resource{map[string]any{}, annotation_model.Annotation{Properties: map[string]any{}}, []string{}, []string{}, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_Validate(t *testing.T) {
	type fields struct {
		Device           map[string]any
		Annotation       annotation_model.Annotation
		ResourceGroupIDs []string
		NodeIDs          []string
		Detected         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"Normal Case: Returns true for a non-empty Resource struct",
			fields{
				map[string]any{"deviceID": "001"},
				annotation_model.Annotation{Properties: map[string]any{"available": true}},
				[]string{"00001"},
				[]string{"node001"},
				false,
			},
			true,
		},
		{
			"Normal Case: Returns false for an empty Resource struct",
			fields{
				map[string]any{},
				annotation_model.Annotation{Properties: map[string]any{}},
				[]string{},
				[]string{},
				false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Device:           tt.fields.Device,
				Annotation:       tt.fields.Annotation,
				ResourceGroupIDs: tt.fields.ResourceGroupIDs,
				NodeIDs:          tt.fields.NodeIDs,
				Detected:         tt.fields.Detected,
			}
			if got := r.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_ToObject(t *testing.T) {
	type fields struct {
		Device           map[string]any
		Annotation       annotation_model.Annotation
		ResourceGroupIDs []string
		NodeIDs          []string
		Detected         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal Case: Successfully converts a Resource struct to a map",
			fields{
				map[string]any{"deviceID": "001"},
				annotation_model.Annotation{Properties: map[string]any{"available": true}},
				[]string{"00001"},
				[]string{"node001"},
				true,
			},
			map[string]any{
				"device":           map[string]any{"deviceID": "001"},
				"annotation":       map[string]any{"available": true},
				"resourceGroupIDs": []string{"00001"},
				"nodeIDs":          []string{"node001"},
				"detected":         true,
			},
		},
		{
			"Normal Case: Returns nil for an empty Resource struct",
			fields{
				map[string]any{},
				annotation_model.Annotation{Properties: map[string]any{}},
				[]string{},
				[]string{},
				false,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Device:           tt.fields.Device,
				Annotation:       tt.fields.Annotation,
				ResourceGroupIDs: tt.fields.ResourceGroupIDs,
				NodeIDs:          tt.fields.NodeIDs,
				Detected:         tt.fields.Detected,
			}
			if got := r.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_ToObject4Node(t *testing.T) {
	type fields struct {
		Device           map[string]any
		Annotation       annotation_model.Annotation
		ResourceGroupIDs []string
		NodeIDs          []string
		Detected         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal Case: Excludes NodeIDs from the output structure even if the input structure has Node elements",
			fields{
				map[string]any{"deviceID": "001"},
				annotation_model.Annotation{Properties: map[string]any{"available": true}},
				[]string{"00001"},
				[]string{},
				true,
			},
			map[string]any{
				"device":           map[string]any{"deviceID": "001"},
				"annotation":       map[string]any{"available": true},
				"resourceGroupIDs": []string{"00001"},
				"detected":         true,
			},
		},
		{
			"Normal Case: Returns nil for an empty Resource struct",
			fields{
				map[string]any{},
				annotation_model.Annotation{Properties: map[string]any{}},
				[]string{},
				[]string{},
				false,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Device:           tt.fields.Device,
				Annotation:       tt.fields.Annotation,
				ResourceGroupIDs: tt.fields.ResourceGroupIDs,
				NodeIDs:          tt.fields.NodeIDs,
				Detected:         tt.fields.Detected,
			}
			if got := r.ToObject4Node(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.ToObject4Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_ToObject4Unused(t *testing.T) {
	type fields struct {
		Device           map[string]any
		Annotation       annotation_model.Annotation
		ResourceGroupIDs []string
		NodeIDs          []string
		Detected         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			"Normal Case: Excludes NodeIDs from the output structure even if the input structure has Node elements",
			fields{
				map[string]any{"deviceID": "001"},
				annotation_model.Annotation{Properties: map[string]any{"available": true}},
				[]string{"00001"},
				[]string{},
				true,
			},
			map[string]any{
				"device":           map[string]any{"deviceID": "001"},
				"annotation":       map[string]any{"available": true},
				"resourceGroupIDs": []string{"00001"},
			},
		},
		{
			"Normal Case: Returns nil for an empty Resource struct",
			fields{
				map[string]any{},
				annotation_model.Annotation{Properties: map[string]any{}},
				[]string{},
				[]string{},
				false,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Device:           tt.fields.Device,
				Annotation:       tt.fields.Annotation,
				ResourceGroupIDs: tt.fields.ResourceGroupIDs,
				NodeIDs:          tt.fields.NodeIDs,
				Detected:         tt.fields.Detected,
			}
			if got := r.ToObject4Unused(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.ToObject4Unused() = %v, want %v", got, tt.want)
			}
		})
	}
}
