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
	"reflect"
	"testing"
)

func TestNewNodeRepository(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name string
		args args
		want NodeRepository
	}{
		{
			"Normal case: Create an instance of the NodeRepository struct ('001')",
			args{"001"},
			NodeRepository{
				"001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeRepository(tt.args.nodeID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeRepository_Find(t *testing.T) {
	t.Skip("not test")
}
