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
        
package group_repository

import (
	"reflect"
	"testing"
)

func TestNewDeleteGroupRepository(t *testing.T) {
	type args struct {
		groupID string
	}
	tests := []struct {
		name string
		args args
		want DeleteGroupRepository
	}{
		{
			"Normal case: Create an instance of the DeleteGroupRepository struct",
			args{"001"},
			DeleteGroupRepository{
				"001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteGroupRepository(tt.args.groupID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteGroupRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteGroupRepository_Delete(t *testing.T) {
	t.Skip("not test")
}
