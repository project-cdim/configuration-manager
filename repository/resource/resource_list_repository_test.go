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

package resource_repository

import (
	"reflect"
	"testing"
)

func TestNewResourceListRepository(t *testing.T) {
	type args struct {
		detail bool
	}
	tests := []struct {
		name string
		args args
		want ResourceListRepository
	}{
		{
			"Normal case: Creates an instance of the ResourceListRepository structure (argument: true)",
			args{true},
			ResourceListRepository{
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceListRepository(tt.args.detail); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceListRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceListRepository_FindList(t *testing.T) {
	t.Skip("not test")
}

func Test_getQueryResourceList(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Normal case",
			queryResourceList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getQueryResourceList(); got != tt.want {
				t.Errorf("getQueryResourceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

const queryResourceList string = `
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs:%s)
OPTIONAL MATCH (vrs)-[:Have]->(van)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt:NotDetected]->(: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[:Compose]->(vrs)
RETURN vrs,
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END,
	COLLECT(vrsg.id),
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`
