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

func TestNewResourceList(t *testing.T) {
	tests := []struct {
		name string
		want ResourceList
	}{
		{
			"Normal case: Create an instance of the Resource struct",
			ResourceList{[]Resource{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceList_ToObject(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: Successfully convert Resource struct to map",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{"node001"},
						false,
					},
					{
						map[string]any{"deviceID": "002"},
						annotation_model.Annotation{Properties: map[string]any{"available": false}},
						[]string{"00002"},
						[]string{"node002"},
						true,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
					"nodeIDs":          []string{"node001"},
					"detected":         false,
				},
				{
					"device":           map[string]any{"deviceID": "002"},
					"annotation":       map[string]any{"available": false},
					"resourceGroupIDs": []string{"00002"},
					"nodeIDs":          []string{"node002"},
					"detected":         true,
				},
			},
		},
		{
			"Normal case: If an empty Resource struct is included, convert to map excluding the empty",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{"node001"},
						false,
					},
					{
						map[string]any{},
						annotation_model.Annotation{Properties: map[string]any{}},
						[]string{},
						[]string{},
						false,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
					"nodeIDs":          []string{"node001"},
					"detected":         false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := &ResourceList{
				Resources: tt.fields.Resources,
			}
			if got := rl.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceList.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceList_ToObject4Node(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: Successfully convert Resource struct to map",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{},
						false,
					},
					{
						map[string]any{"deviceID": "002"},
						annotation_model.Annotation{Properties: map[string]any{"available": false}},
						[]string{"00002"},
						[]string{},
						true,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
					"detected":         false,
				},
				{
					"device":           map[string]any{"deviceID": "002"},
					"annotation":       map[string]any{"available": false},
					"resourceGroupIDs": []string{"00002"},
					"detected":         true,
				},
			},
		},
		{
			"Normal case: If an empty Resource struct is included, convert to map excluding the empty",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{},
						false,
					},
					{
						map[string]any{},
						annotation_model.Annotation{Properties: map[string]any{}},
						[]string{},
						[]string{},
						false,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
					"detected":         false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := &ResourceList{
				Resources: tt.fields.Resources,
			}
			if got := rl.ToObject4Node(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceList.ToObject4Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceList_ToObject4Unused(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: Successfully convert Resource struct to map",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{},
						false,
					},
					{
						map[string]any{"deviceID": "002"},
						annotation_model.Annotation{Properties: map[string]any{"available": false}},
						[]string{"00002"},
						[]string{},
						true,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
				},
				{
					"device":           map[string]any{"deviceID": "002"},
					"annotation":       map[string]any{"available": false},
					"resourceGroupIDs": []string{"00002"},
				},
			},
		},
		{
			"Normal case: If an empty Resource struct is included, convert to map excluding the empty",
			fields{
				[]Resource{
					{
						map[string]any{"deviceID": "001"},
						annotation_model.Annotation{Properties: map[string]any{"available": true}},
						[]string{"00001"},
						[]string{},
						false,
					},
					{
						map[string]any{},
						annotation_model.Annotation{Properties: map[string]any{}},
						[]string{},
						[]string{},
						false,
					},
				},
			},
			[]map[string]any{
				{
					"device":           map[string]any{"deviceID": "001"},
					"annotation":       map[string]any{"available": true},
					"resourceGroupIDs": []string{"00001"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := &ResourceList{
				Resources: tt.fields.Resources,
			}
			if got := rl.ToObject4Unused(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceList.ToObject4Unused() = %v, want %v", got, tt.want)
			}
		})
	}
}
