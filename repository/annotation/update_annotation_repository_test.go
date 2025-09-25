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

func TestUpdateAnnotationConstructsWhereClause(t *testing.T) {
	want := `'CPU' IN labels(vrs) OR 'Accelerator' IN labels(vrs) OR 'DSP' IN labels(vrs) OR 'FPGA' IN labels(vrs) OR 'GPU' IN labels(vrs) OR 'UnknownProcessor' IN labels(vrs) OR 'Memory' IN labels(vrs) OR 'Storage' IN labels(vrs) OR 'NetworkInterface' IN labels(vrs) OR 'GraphicController' IN labels(vrs) OR 'VirtualMedia' IN labels(vrs)`

	t.Run("verify that the return value is as expected", func(t *testing.T) {
		got := updateAnnotationConstructsWhereClause()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("updateAnnotationConstructsWhereClause() results = %v, want %v", got, want)
		}
	})
}

func TestNewUpdateAnnotationRepository(t *testing.T) {
	deviceIDs := []string{"device1", "device2"}

	repo := NewUpdateAnnotationRepository(deviceIDs)

	if !reflect.DeepEqual(repo.deviceIDs, deviceIDs) {
		t.Errorf("NewUpdateAnnotationRepository().deviceIDs = %v, want %v", repo.deviceIDs, deviceIDs)
	}
}

func TestUpdateAnnotationRepository_Set(t *testing.T) {
	t.Skip("not test")
}
