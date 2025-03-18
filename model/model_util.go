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
        
package model

import "time"

// CurrentTimeISO8601 returns the current time in UTC formatted according to ISO 8601 standard.
// The format used is: YYYY-MM-DDThh:mm:ssZ.
func CurrentTimeISO8601() string {
	// Get the current time in UTC time zone
	now := time.Now().UTC()
	// ISO 8601 format: YYYY-MM-DDThh:mm:ssZ
	res := now.Format("2006-01-02T15:04:05Z07:00")
	return res
}

// ValidateISO8601 checks if the provided string is in the ISO 8601 format.
// It returns true if the string is a valid ISO 8601 date-time, otherwise false.
//
// The expected format is "2006-01-02T15:04:05Z".
func ValidateISO8601(s string) bool {
	layout := "2006-01-02T15:04:05Z07:00"
	_, err := time.Parse(layout, s)
	return err == nil
}
