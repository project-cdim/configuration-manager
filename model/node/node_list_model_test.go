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
        
package node_model

import (
	"reflect"
	"testing"

	resource_model "github.com/project-cdim/configuration-manager/model/resource"
)

func TestNewNodeList(t *testing.T) {
	tests := []struct {
		name string
		want NodeList
	}{
		{
			"Normal case: Create an instance of the NodeList struct",
			NodeList{
				[]Node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeList_ToObject(t *testing.T) {
	type fields struct {
		Nodes []Node
	}
	tests := []struct {
		name   string
		fields fields
		want   []map[string]any
	}{
		{
			"Normal case: When NodeList struct has all valid Node structs, convert all to map",
			fields{
				[]Node{
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
			"Normal case: When NodeList struct includes invalid Node structs (no id or id length is 0), convert excluding the invalid Node structs to map",
			fields{
				[]Node{
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
			nl := &NodeList{
				Nodes: tt.fields.Nodes,
			}
			if got := nl.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeList.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
