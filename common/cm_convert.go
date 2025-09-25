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

package common

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Map2CypherProperty converts a map[string]any to a Cypher query property string.
// It iterates over the input map, sorts the keys alphabetically, and then constructs
// a string representation of the map in the format of Cypher properties. For each key-value
// pair, the key is added without double quotes, followed by a colon, and then the value is
// converted to a string suitable for Cypher query syntax. If the value conversion fails,
// an error is returned. The resulting string is enclosed in curly braces.
//
// Example:
// Given an input map[string]any{"name": "John", "age": 30},
// the function returns "{age:30,name:\"John\"}", nil.
func Map2CypherProperty(input map[string]any) (string, error) {
	result := "{"

	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		// For the second and subsequent times, concatenate the item values with a comma.
		if i > 0 {
			result += ","
		}
		// Use the key without adding double quotation marks, append a colon, and use it as the output value.
		result += k + ":"

		convertResult, err := convertByType(input[k])
		if err != nil {
			return "", err
		}
		result += convertResult

	}
	result += "}"

	return result, nil
}

// Slice2CypherProperty converts a slice of any type into a string representation
// suitable for Cypher properties. It handles various types of input, including
// simple types like strings and complex types like structs, by delegating to
// convertByType function for each element in the slice. The output is formatted
// as a Cypher array, with appropriate conversions applied to each element to ensure
// it is compatible with Cypher query syntax.
//
// Output examples include simple arrays like ["memory","CPU"], arrays of objects
// like [{type: "memory",deviceID: "res202"}], and nested arrays like [["memory","CPU"], ["memory","CPU"]].
//
// Parameters:
// - input: A slice of any type that needs to be converted into a Cypher property array.
//
// Returns:
// - A string representation of the input slice formatted as a Cypher property array.
// - An error if any element in the slice cannot be converted to a suitable Cypher property.
func Slice2CypherProperty(input []any) (string, error) {
	result := "["

	for i, v := range input {
		// For the second and subsequent times, concatenate the item values with a comma.
		if i > 0 {
			result += ","
		}

		convertResult, err := convertByType(v)
		if err != nil {
			return "", err
		}
		result += convertResult
	}

	result += "]"
	return result, nil
}

// number2string converts a numeric input of various types to its string representation.
// It supports conversion for integer types (Int, Int8, Int16, Int32, Int64, Uint, Uint8,
// Uint16, Uint32, Uint64, Uintptr) and floating-point types (Float32, Float64). The function
// uses reflection to determine the input type and then formats it accordingly as a string.
//
// If the input is an integer type, it is formatted without any decimal point. If the input
// is a floating-point type, it is formatted using the 'g' verb to ensure the most compact
// and efficient representation. If the input type is not a supported number type, the function
// returns an error indicating a conversion issue.
//
// Parameters:
// - input: The numeric value to be converted to a string. It can be any numeric type.
//
// Returns:
//   - A string representation of the input number.
//   - An error if the input is not a supported numeric type, including a specific error code
//     and message indicating the nature of the conversion error.
func number2string(input any) (string, error) {
	switch reflect.ValueOf(input).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", input), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", input), nil
	default:
		return "", fmt.Errorf("input type is not a supported numeric type. type(%v), value(%v)", reflect.ValueOf(input).Kind(), reflect.ValueOf(input))
	}
}

// Any2anyslice takes an interface{} value and attempts to convert it into a slice of interface{}.
// If the input value is already a slice, it iterates through the slice and adds each element to a new
// slice of interface{}, preserving the order of elements. If the input value is not a slice, it returns
// an empty slice of interface{}. This function is useful for converting slices of a specific type to
// slices of interface{}, allowing for more generic handling of collections of values.
//
// Parameters:
// - anyValue: The value to be converted into a slice of interface{}. Can be of any type.
//
// Returns:
// - A slice of interface{} containing the elements of the input slice if the input was a slice.
// - An empty slice of interface{} if the input was not a slice.
func Any2anyslice(anyValue any) []any {
	reflectValue := reflect.ValueOf(anyValue)
	if reflectValue.Kind() != reflect.Slice {
		return []any{}
	}

	result := make([]any, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		result[i] = reflectValue.Index(i).Interface()
	}

	return result
}

// convertByType takes an input of any type and converts it into a string representation suitable for Cypher queries.
// It handles various data types including nil, string, numeric types, boolean, slices, and maps by applying specific
// formatting rules for each type. This function is designed to ensure that the converted values are compatible with
// Cypher query syntax, facilitating the construction of dynamic queries.
//
// The conversion rules are as follows:
// - nil values are converted to the string "null".
// - Strings are enclosed in double quotation marks.
// - Numeric types are converted to their string representation without quotation marks.
// - Boolean values are converted to "true" or "false" without quotation marks.
// - Slices are recursively processed to generate a Cypher array representation.
// - Maps are recursively processed to generate a Cypher map representation.
// - For unsupported types, an error is returned indicating a marshal type error.
//
// Parameters:
// - anyValue: The value to be converted, which can be of any type.
//
// Returns:
// - A string representation of the input value formatted for Cypher query compatibility.
// - An error if the input type is unsupported or if an error occurs during the conversion of slices or maps.
func convertByType(anyValue any) (string, error) {
	// If the value is null in JSON, return the string "null" as the value to specify in the Cypher query.
	if anyValue == nil {
		return "null", nil
	}

	result := ""
	switch reflect.ValueOf(anyValue).Kind() {
	case reflect.String:
		// Use the value of v by enclosing it in double quotation marks as the output value.
		result = fmt.Sprintf("%q", anyValue.(string))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		// Use the value of v without enclosing it in double quotation marks as the output value.
		valueStr, _ := number2string(anyValue)
		result = valueStr

	case reflect.Bool:
		// Use the value of v without adding double quotation marks as the output value.
		result = strconv.FormatBool(anyValue.(bool))

	case reflect.Slice:
		// Generated value when the array value is a string: Example: ["memory", "CPU"]
		arrayRecursiveResult, err := Slice2CypherProperty(anyValue.([]any))
		if err != nil {
			return "", err
		}
		result = arrayRecursiveResult

	case reflect.Map:
		// Generated value: Example: {state:"Enabled", health:"OK"}
		mapRecursiveResult, err := Map2CypherProperty(anyValue.(map[string]any))
		if err != nil {
			return "", err
		}
		result = mapRecursiveResult

	default:
		// Return an error for unexpected types (since there are no places within cmapi that generate Arrays, Arrays will also result in an error)
		return "", fmt.Errorf("input type is unsupported. type(%v), value(%v)", reflect.ValueOf(anyValue).Kind(), reflect.ValueOf(anyValue))
	}

	return result, nil
}

// Nil2EmptyFromMap iterates through a map and converts any nil values within the map to empty values.
// This function is useful for preparing maps for operations that do not handle nil values well, ensuring
// that all keys have non-nil values. The conversion is done in-place, modifying the original map.
//
// The function handles various types of values stored in the map, including nested maps and slices, by
// recursively applying the conversion to ensure that no nil values are left in any nested structures.
//
// Input examples include:
// - map{aa: "aaa1", bb: "bbb1"}
// - map{ccc: []}
// - map{ccc: [map{aa: "aaa2", bb: "bbb2"}, map{aa: "aaa3", bb: "bbb3"}]}
// - map{ccc: map{aa: "aaa1", bb: "bbb1"}}
// - map{ccc: map{aa: ["aaa1", "bbb1"], bb: []}}
//
// Parameters:
// - input: The map to be processed, where any nil values will be converted to empty values.
//
// Returns:
// - The original map with nil values converted to empty values.
func Nil2EmptyFromMap(input map[string]any) map[string]any {
	for k, v := range input {
		input[k] = nil2Empty(v)
	}
	return input
}

// Nil2EmptyFromSlice iterates over a slice and converts any nil values within the slice to empty values.
// This function is particularly useful for data normalization before processing or storing, ensuring that
// all elements in the slice are non-nil. The conversion is applied in-place, directly modifying the original
// slice. It leverages a helper function, nil2Empty, to perform the actual conversion for each element.
//
// Input examples include:
// - Simple slices like ["aa", "bb"] where no conversion is needed.
// - Slices of maps like [{aa: "aaa1", bb: "bbb1"}, {aa: "aaa2", bb: "bbb2"}] where maps may contain nil values.
// - Nested slices like [["aa1", "aa2"], []] where inner slices may be nil or empty.
//
// Parameters:
// - input: The slice to be processed, where any nil values will be converted to empty values.
//
// Returns:
// - The original slice with nil values converted to empty values.
func Nil2EmptyFromSlice(input []any) []any {
	for i, v := range input {
		input[i] = nil2Empty(v)
	}
	return input
}

// nil2Empty checks the type of the input value and converts nil slices and maps to their empty equivalents.
// For slices, if the input is nil, it returns an empty slice of type []any. For maps, it returns an empty map
// of type map[string]any if the input is nil. This function is recursive for slices, ensuring that all nested
// slices are also converted from nil to empty. For maps, it leverages the Nil2EmptyFromMap function to recursively
// convert nil values within the map to empty values. For all other types, the input value is returned unchanged.
//
// This function is useful for data normalization, ensuring that data structures do not contain nil values, which
// can be problematic for certain operations or when encoding to formats like JSON.
//
// Parameters:
// - anyValue: The value to be checked and potentially converted from nil to an empty structure.
//
// Returns:
// - An equivalent value with nil slices and maps converted to their empty counterparts.
// - For types other than slice or map, the original value is returned unchanged.
func nil2Empty(anyValue any) any {
	switch reflect.TypeOf(anyValue).Kind() {
	case reflect.Slice:
		if reflect.ValueOf(anyValue).IsNil() {
			return []any{}
		} else {
			return Nil2EmptyFromSlice(anyValue.([]any))
		}
	case reflect.Map:
		return Nil2EmptyFromMap(anyValue.(map[string]any))
	default:
		return anyValue
	}
}

// UnquoteRecursive performs recursive unquoting on the input value.
//
// This function must be called during the response phase.
// In both request and response phase, control characters such as tabs undergo escape sequence processing (e.g., \t).
// If this function is not called, escape sequences will be duplicated (e.g., \\t).
//
// This function processes the following types:
//   - string: Unquotes using strconv.Unquote
//   - slice/array: Recursively applies UnquoteRecursive to each element
//   - map: Recursively applies UnquoteRecursive to each value (keys are treated as strings)
//   - other types: Returns as-is
//
// Parameters:
//   - input: The value to be unquoted
//
// Returns:
//   - any: The unquoted value. Returns nil if nil is passed or if unquoting fails
//   - error: An error if unquoting fails for a string value, otherwise nil
func UnquoteRecursive(input any) (any, error) {
	if input == nil {
		return nil, nil
	}

	val := reflect.ValueOf(input)
	kind := val.Kind()

	switch kind {
	case reflect.String:
		s := `"` + val.String() + `"`
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			// This error should not occur because the string to be unquoted is always quoted beforehand.
			Log.Error(err.Error())
			return nil, err
		}
		return unquoted, nil
	case reflect.Slice, reflect.Array:
		newSlice := make([]any, val.Len())

		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)

			// Convert element to any type and recursively call itself
			unquotedElem, err := UnquoteRecursive(elem.Interface())
			if err != nil {
				return nil, err
			}
			newSlice[i] = unquotedElem
		}
		return newSlice, nil
	case reflect.Map:
		newMap := make(map[string]any, val.Len())

		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			// Convert value to any type and recursively call itself
			unquotedValue, err := UnquoteRecursive(v.Interface())
			if err != nil {
				return nil, err
			}
			newMap[k.String()] = unquotedValue
		}
		return newMap, nil
	default:
		return input, nil
	}
}
