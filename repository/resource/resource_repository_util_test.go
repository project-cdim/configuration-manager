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
        
package resource_repository

import (
	"reflect"
	"testing"

	annotation_model "github.com/project-cdim/configuration-manager/model/annotation"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"

	"github.com/apache/age/drivers/golang/age"
)

func TestComposeResource(t *testing.T) {
	var valNilSlice []any
	type args struct {
		resVertex        *age.Vertex
		annotationVertex *age.Vertex
		resourceGroupIDs *age.SimpleEntity
		nodeIDs          *age.SimpleEntity
		detected         bool
		detail           bool
	}
	tests := []struct {
		name string
		args args
		want resource_model.Resource
	}{
		{
			"Normal case: When resVertex's Property is empty, returns an empty model",
			args{
				age.NewVertex(10, "label10", map[string]any{}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.NewResource(),
		},
		{
			"Normal case: If the annotationVertex does not have an available element, available defaults to true",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If the annotationVertex has an available element but it's not a bool, available defaults to true",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": "false",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If the annotationVertex has an available element and it's true, available is true",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": true,
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If the annotationVertex has an available element and it's false, available is false",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": false,
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": false}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If detail is true, use all elements of resVertex",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": false,
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": false}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If detail is false, only use 'deviceID, type, status' elements of resVertex",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": false,
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				false,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": false}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: If detail is true and links is a nil slice, the links element in the Device part of the return value is an empty array",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": valNilSlice,
				}),
				age.NewVertex(20, "label20", map[string]any{
					"available": false,
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": false}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: When the argument 'detected' is true, the value of the 'Detected' element in the model should be true",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: When the argument 'detected' is false, the value of the 'Detected' element in the model should be false",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				false,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         false,
			},
		},
		{
			"Normal case: When the argument 'resourceGroupIDs' is an empty slice, the 'ResourceGroupIDs' element in the model should also be an empty slice",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: When the argument 'resourceGroupIDs' is not empty, the value of the 'ResourceGroupIDs' element in the model should be set to the value of the argument",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{"aaa", "bbb"}),
				age.NewSimpleEntity([]any{}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{"aaa", "bbb"},
				NodeIDs:          []string{},
				Detected:         true,
			},
		},
		{
			"Normal case: When the argument 'nodeIDs' is not empty, the value of the 'nodeIDs' element in the model should be set to the value of the argument",
			args{
				age.NewVertex(10, "label10", map[string]any{
					"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"},
				}),
				age.NewVertex(20, "label20", map[string]any{
					"test": "test",
				}),
				age.NewSimpleEntity([]any{"aaa", "bbb"}),
				age.NewSimpleEntity([]any{"ccc", "ddd"}),
				true,
				true,
			},
			resource_model.Resource{
				Device:           map[string]any{"deviceID": "id10", "type": "CPU", "status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"id11", "id12"}},
				Annotation:       annotation_model.Annotation{Properties: map[string]any{"available": true}},
				ResourceGroupIDs: []string{"aaa", "bbb"},
				NodeIDs:          []string{"ccc", "ddd"},
				Detected:         true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeResource(tt.args.resVertex, tt.args.annotationVertex, tt.args.resourceGroupIDs, tt.args.nodeIDs, tt.args.detected, tt.args.detail); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPrimaryDeviceProp(t *testing.T) {
	type args struct {
		prop map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			"Normal case: Only extracts and returns the deviceID, type, and status elements",
			args{
				map[string]any{
					"deviceID": "res001",
					"type":     "CPU",
					"status":   map[string]any{"state": "Enabled", "health": "OK"},
					"dummy":    "dummyVal",
				},
			},
			map[string]any{
				"deviceID": "res001",
				"type":     "CPU",
				"status":   map[string]any{"state": "Enabled", "health": "OK"},
			},
		},
		{
			"Normal case: If there is no status element, it returns without adding a status element",
			args{
				map[string]any{
					"deviceID": "res001",
					"type":     "CPU",
					"dummy":    "dummyVal",
				},
			},
			map[string]any{
				"deviceID": "res001",
				"type":     "CPU",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPrimaryDeviceProp(tt.args.prop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractPrimaryDeviceProp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortResourceList(t *testing.T) {
	type args struct {
		resources []resource_model.Resource
	}
	tests := []struct {
		name string
		args args
		want []resource_model.Resource
	}{
		{
			"Normal case: Sort by resourceType and deviceID",
			args{
				[]resource_model.Resource{
					{Device: map[string]any{"deviceID": "id07", "type": "GPU"}},
					{Device: map[string]any{"deviceID": "id02", "type": "memory"}},
					{Device: map[string]any{"deviceID": "id05", "type": "CPU"}},
					{Device: map[string]any{"deviceID": "id06", "type": "UnknownProcessor"}},
					{Device: map[string]any{"deviceID": "id01", "type": "memory"}},
					{Device: map[string]any{"deviceID": "id03", "type": "CPU"}},
					{Device: map[string]any{"deviceID": "id04", "type": "memory"}},
					{Device: map[string]any{"deviceID": "id08", "type": "cpu"}},
				},
			},
			[]resource_model.Resource{
				{Device: map[string]any{"deviceID": "id03", "type": "CPU"}},
				{Device: map[string]any{"deviceID": "id05", "type": "CPU"}},
				{Device: map[string]any{"deviceID": "id07", "type": "GPU"}},
				{Device: map[string]any{"deviceID": "id06", "type": "UnknownProcessor"}},
				{Device: map[string]any{"deviceID": "id08", "type": "cpu"}},
				{Device: map[string]any{"deviceID": "id01", "type": "memory"}},
				{Device: map[string]any{"deviceID": "id02", "type": "memory"}},
				{Device: map[string]any{"deviceID": "id04", "type": "memory"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if sortResourceList(tt.args.resources); !reflect.DeepEqual(tt.args.resources, tt.want) {
				t.Errorf("sortResourceList() = %v, want %v", tt.args.resources, tt.want)
			}
		})
	}
}
