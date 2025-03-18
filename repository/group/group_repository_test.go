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

	"github.com/project-cdim/configuration-manager/common"
)

func TestNewGroupRepository(t *testing.T) {
	type args struct {
		groupID       string
		withResources bool
	}
	tests := []struct {
		name string
		args args
		want GroupRepository
	}{
		{
			"Normal case: Create an instance of the GroupRepository struct",
			args{common.DefaultGroupId, true},
			GroupRepository{
				common.DefaultGroupId,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroupRepository(tt.args.groupID, tt.args.withResources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupRepository_Find(t *testing.T) {
	t.Skip("not test")
}
