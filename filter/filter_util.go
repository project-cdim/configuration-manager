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

// Value specifying resource available or unavailable in search condition.
const (
	TrueStr  = "true"  // true: Search only available resources
	FalseStr = "false" // false: Search only unavailable resources
	AllStr   = "all"   // all: Search all resources available and unavailable
)

// Strings specifying whether a resource is available or unavailable in the search condition.
var ResourceUseValueType = []string{
	TrueStr,  // true: Search only for resources that are available
	FalseStr, // false: Search only for resources that are unavailable
	AllStr,   // all: Search for all resources, both available and unavailable
}

// ResourceUseType is value to determine available resources.
type ResourceUseType int

const (
	ResourceUseAvailable   ResourceUseType = iota // available resources
	ResourceUseUnAvailable                        // unavailable resources
	ResourceUseAll                                // all resources
)

// Logical operators for search conditions.
const (
	AND = "and"
	OR  = "or"
	NOT = "not"
)

// ExpressionValueType is Logical operators for search conditions
var ExpressionValueType = []string{
	AND, // and
	OR,  // or
	NOT, // not
}

// ResourceExpressionType is a logical operator for multiple values in a single condition item.
type ResourceExpressionType int

const (
	RresourceExpressionAnd ResourceExpressionType = iota
	RresourceExpressionNot
	RresourceExpressionOr
)
