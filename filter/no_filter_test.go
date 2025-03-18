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
        
package filter

import (
	"reflect"
	"testing"
)

func TestNewNoFilter(t *testing.T) {
	tests := []struct {
		name string
		want noFilter
	}{
		{
			"Normal case: Generate an instance of the NoFilter structure",
			noFilter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoFilter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoFilter_FilterByCondition(t *testing.T) {
	type args struct {
		rec map[string]any
		opt []any
	}
	tests := []struct {
		name string
		nsc  noFilter
		args args
		want bool
	}{
		{
			"Normal case: Execute the filter, the result is always true",
			NewNoFilter(),
			args{
				map[string]any{"test": "testVal"},
				[]any{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nsc := tt.nsc
			if got := nsc.FilterByCondition(tt.args.rec, tt.args.opt...); got != tt.want {
				t.Errorf("NoFilter.FilterByCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
