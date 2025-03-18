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
        
package node_repository

import (
	"testing"

	"github.com/apache/age/drivers/golang/age"
)

func Test_compareByNode(t *testing.T) {
	type args struct {
		records     [][]age.Entity
		nodeIdx     int
		resourceIdx int
		i           int
		j           int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Normal case: Lexicographically compare the node IDs as the first sort key, return true if 'the ID of the former node < the ID of the latter node'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Node", map[string]any{"id": "Node01"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				0, 1,
				0, 1,
			},
			true,
		},
		{
			"Normal case: Lexicographically compare the node IDs as the first sort key, return false if 'the ID of the former node >= the ID of the latter node'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Node", map[string]any{"id": "Node01"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				0, 1,
				0, 1,
			},
			false,
		},
		{
			"Normal case: If the first sort key is the same, lexicographically compare the deviceID as the second sort key, return true if 'the deviceID of the former < the deviceID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU01"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				0, 1,
				0, 1,
			},
			true,
		},
		{
			"Normal case: If the first sort key is the same, lexicographically compare the deviceID as the second sort key, return false if 'the deviceID of the former >= the deviceID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Node", map[string]any{"id": "Node02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU01"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				0, 1,
				0, 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareByNode(tt.args.records, tt.args.nodeIdx, tt.args.resourceIdx, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("compareByNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
