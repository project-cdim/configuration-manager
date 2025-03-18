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

func TestNewResourceRepository(t *testing.T) {
	type args struct {
		deviceID string
	}
	tests := []struct {
		name string
		args args
		want ResourceRepository
	}{
		{
			"Normal case: Creates an instance of the ResourceRepository structure ('001')",
			args{"001"},
			ResourceRepository{
				"001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceRepository(tt.args.deviceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRepository_Find(t *testing.T) {
	t.Skip("not test")
}

func Test_getQueryResource(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Normal case",
			queryResource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getQueryResource("res101"); got != tt.want {
				t.Errorf("getQueryResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

const queryResource string = `
MATCH (vrs: CPU{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: Accelerator{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: DSP{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: FPGA{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: GPU{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: UnknownProcessor{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: Memory{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: Storage{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: NetworkInterface{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: GraphicController{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END
UNION ALL
MATCH (vrs: VirtualMedia{deviceID: 'res101'})
OPTIONAL MATCH (vrs)-[ehv: Have]->(van)
OPTIONAL MATCH (vrsg)-[ein: Include]->(vrs)
OPTIONAL MATCH (vrs)-[endt: NotDetected]->(vndd: NotDetectedDevice)
OPTIONAL MATCH (vnd)-[ecm: Compose]->(vrs)
RETURN vrs, 
	CASE WHEN van IS NULL THEN {id:-1, label:"dummy", properties: {}}::vertex ELSE van END, 
	COLLECT(vrsg.id), 
	COLLECT(vnd.id),
	CASE WHEN endt IS NULL THEN true ELSE false END`
