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
        
package rack_repository

import (
	"testing"

	"github.com/apache/age/drivers/golang/age"
)

func Test_compareByRack(t *testing.T) {
	type args struct {
		records    [][]age.Entity
		chassisIdx int
		deviceIdx  int
		i          int
		j          int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Normal case: Numerically compares the first sort key, unitPosition, and returns true if 'the former unitPosition < the latter unitPosition'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(3)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			true,
		},
		{
			"Normal case: Numerically compares the first sort key, unitPosition, and returns true if 'the former unitPosition < the latter unitPosition' (ensuring it's not a lexicographical comparison)",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(11)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			true,
		},
		{
			"Normal case: Numerically compares the first sort key, unitPosition, and returns false if 'the former unitPosition >= the latter unitPosition'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(1)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			false,
		},
		{
			"Normal case: With the first sort key being the same, lexicographically compares the second sort key, type, and returns true if 'the former type < the latter type'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "Memory", map[string]any{"deviceID": "Mem02", "type": "Memory"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			true,
		},
		{
			"Normal case: With the first sort key being the same, lexicographically compares the second sort key, type, and returns false if 'the former type >= the latter type'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "Memory", map[string]any{"deviceID": "Mem02", "type": "Memory"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			false,
		},
		{
			"Normal case: With the first two sort keys being the same, lexicographically compares the third sort key, deviceID, and returns true if 'the former deviceID < the latter deviceID'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU01", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			true,
		},
		{
			"Normal case: If up to the second sort key is the same, lexicographically compare the deviceID as the third sort key, return false if 'the deviceID of the former >= the deviceID of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"deviceID": "CPU02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "CPU", map[string]any{"deviceID": "CPU01", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			false,
		},
		{
			"Normal case: If up to the third sort key is the same, lexicographically compare the id as the fourth sort key, return true if 'the id of the former < the id of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"id": "01", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "CPU", map[string]any{"id": "02", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			true,
		},
		{
			"Normal case: If up to the third sort key is the same, lexicographically compare the id as the fourth sort key, return false if 'the id of the former >= the id of the latter'",
			args{
				[][]age.Entity{
					{
						age.NewVertex(10, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(11, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(12, "CPU", map[string]any{"id": "02", "type": "CPU"}),
						age.NewVertex(13, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
					{
						age.NewVertex(20, "Rack", map[string]any{"id": "Rack02"}),
						age.NewVertex(21, "Chassis", map[string]any{"id": "Chassis02", "unitPosition": int64(2)}),
						age.NewVertex(22, "CPU", map[string]any{"id": "01", "type": "CPU"}),
						age.NewVertex(23, "Annotation", map[string]any{}),
						age.NewSimpleEntity([]any{}),
						age.NewSimpleEntity(true),
					},
				},
				1, 2,
				0, 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareByRack(tt.args.records, tt.args.chassisIdx, tt.args.deviceIdx, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("compareByRack() = %v, want %v", got, tt.want)
			}
		})
	}
}
