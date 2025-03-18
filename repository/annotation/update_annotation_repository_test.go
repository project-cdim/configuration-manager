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
        
package annotation_repository

import (
	"reflect"
	"testing"
)

func TestNewUpdateAnnotationRepository(t *testing.T) {
	tests := []struct {
		name     string
		deviceID string
		want     UpdateAnnotationRepository
	}{
		{
			name:     "Normal case: Create an instance of the UpdateAnnotationRepository struct",
			deviceID: "device123",
			want:     UpdateAnnotationRepository{deviceID: "device123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUpdateAnnotationRepository(tt.deviceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdateAnnotationRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateAnnotationRepository_Set(t *testing.T) {
	t.Skip("not test")
}
