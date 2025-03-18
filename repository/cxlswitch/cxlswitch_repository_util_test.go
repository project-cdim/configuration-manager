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
        
package cxlswitch_repository

import (
	"testing"

	"github.com/apache/age/drivers/golang/age"
)

func Test_compareByCXLSwitch(t *testing.T) {
	type args struct {
		records     [][]age.Entity
		switchIdx   int
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
			"Normal case: Lexicographically compare the CXLSwitch ID as the first sort key, return true if 'the CXLSwitch ID of the former < the CXLSwitch ID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "CXLSwitch", map[string]any{"id": "Switch01"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
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
			"Normal case: Lexicographically compare the CXLSwitch ID as the first sort key, return false if 'the CXLSwitch ID of the former >= the CXLSwitch ID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "CXLSwitch", map[string]any{"id": "Switch01"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
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
			"Normal case: With the first sort key being the same, lexicographically compare the deviceID as the second sort key, return true if 'the deviceID of the former < the deviceID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU01"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
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
			"Normal case: With the first sort key being the same, lexicographically compare the deviceID as the second sort key, return false if 'the deviceID of the former >= the deviceID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(11, "CPU", map[string]any{"deviceID": "CPU02"}),
						age.NewVertex(12, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "CXLSwitch", map[string]any{"id": "Switch02"}),
						age.NewVertex(21, "CPU", map[string]any{"deviceID": "CPU01"}),
						age.NewVertex(22, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
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
			if got := compareByCXLSwitch(tt.args.records, tt.args.switchIdx, tt.args.resourceIdx, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("compareByCXLSwitch() = %v, want %v", got, tt.want)
			}
		})
	}
}
