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
        
package annotation_model

import (
	"reflect"
	"testing"
)

func TestNewAnnotation(t *testing.T) {
	tests := []struct {
		name string
		want Annotation
	}{
		{
			"Normal case: Create an instance of the Annotation struct",
			Annotation{
				Properties: map[string]any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAnnotation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotation_ToObject(t *testing.T) {
	tests := []struct {
		name   string
		fields Annotation
		want   map[string]any
	}{
		{
			"Normal case: Successfully convert Annotation struct to map",
			Annotation{
				Properties: map[string]any{
					"available": true,
				},
			},
			map[string]any{
				"available": true,
			},
		},
		{
			"Empty properties: Convert Annotation struct with empty properties to map",
			Annotation{
				Properties: map[string]any{},
			},
			map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Annotation{
				Properties: tt.fields.Properties,
			}
			if got := a.ToObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Annotation.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
