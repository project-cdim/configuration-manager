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
	"reflect"
	"testing"
)

func TestMap2CypherProperty(t *testing.T) {
	type args struct {
		input map[string]any
	}
	type testST1 struct {
		F1 string `default:"aaa"`
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"insert_{}",
			args{
				map[string]any{"key1": "value1"},
			},
			"{key1:\"value1\"}",
			false,
		},
		{
			"insert_,and{}",
			args{
				map[string]any{"key1": "value1", "key2": "value2"},
			},
			"{key1:\"value1\",key2:\"value2\"}",
			false,
		},
		{
			"unexpected",
			args{map[string]any{"key1": "value1", "key2": &testST1{}}},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map2CypherProperty(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map2CypherProperty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Map2CypherProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice2CypherProperty(t *testing.T) {
	type args struct {
		input []any
	}
	type testST1 struct {
		F1 string `default:"aaa"`
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"insert_[]",
			args{
				[]any{"value1"},
			},
			"[\"value1\"]",
			false,
		},
		{
			"insert_,and[]",
			args{
				[]any{"value1", "value2"},
			},
			"[\"value1\",\"value2\"]",
			false,
		},
		{
			"unexpected",
			args{
				[]any{&testST1{}},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Slice2CypherProperty(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice2CypherProperty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Slice2CypherProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber2string(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"int", args{int(1)}, "1", false},
		{"int8", args{int8(2)}, "2", false},
		{"int16", args{int16(3)}, "3", false},
		{"int32", args{int32(4)}, "4", false},
		{"int64", args{int64(5)}, "5", false},
		{"uint", args{uint(6)}, "6", false},
		{"uint8", args{uint8(7)}, "7", false},
		{"uint16", args{uint16(8)}, "8", false},
		{"uint32", args{uint32(9)}, "9", false},
		{"uint64", args{uint64(10)}, "10", false},
		{"uintptr", args{uintptr(11)}, "11", false},
		{"float32", args{float32(12)}, "12", false},
		{"float64", args{float64(13)}, "13", false},
		{"bool", args{true}, "", true},
		{"string", args{"test string"}, "", true},
		{"array", args{[2]string{"1", "2"}}, "", true},
		{"slice", args{[]string{"1", "2"}}, "", true},
		{"map", args{map[string]string{"1": "2"}}, "", true},
		{"unexpected", args{nil}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := number2string(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number2string() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number2string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAny2anyslice(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{"slice", args{createTestValue_Slice()}, []any{"value1", "value2"}},
		{"notslice", args{"aaa"}, []any{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Any2anyslice(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Any2anyslice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertByType(t *testing.T) {
	type args struct {
		anyValue any
	}
	type testST1 struct {
		F1 string `default:"aaa"`
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"string", args{"test string"}, "\"test string\"", false},
		{"int", args{int(1)}, "1", false},
		{"int8", args{int8(2)}, "2", false},
		{"int16", args{int16(3)}, "3", false},
		{"int32", args{int32(4)}, "4", false},
		{"int64", args{int64(5)}, "5", false},
		{"uint", args{uint(6)}, "6", false},
		{"uint8", args{uint8(7)}, "7", false},
		{"uint16", args{uint16(8)}, "8", false},
		{"uint32", args{uint32(9)}, "9", false},
		{"uint64", args{uint64(10)}, "10", false},
		{"uintptr", args{uintptr(11)}, "11", false},
		{"float32", args{float32(12)}, "12", false},
		{"float64", args{float64(13)}, "13", false},
		{"bool", args{true}, "true", false},
		{"array", args{createTestValue_Array()}, "", true},
		{"slice", args{createTestValue_Slice()}, "[\"value1\",\"value2\"]", false},
		{"slice,array", args{[]any{createTestValue_Array(), createTestValue_Array()}}, "", true},
		{"slice,slice", args{[]any{createTestValue_Slice(), createTestValue_Slice()}}, "[[\"value1\",\"value2\"],[\"value1\",\"value2\"]]", false},
		{"slice,map", args{[]any{createTestValue_Map()}}, "[{key1:\"value1\",key2:\"value2\"}]", false},
		{"map", args{createTestValue_Map()}, "{key1:\"value1\",key2:\"value2\"}", false},
		{"map,array", args{map[string]any{"key1": createTestValue_Array(), "key2": createTestValue_Array()}}, "", true},
		{"map,slice", args{map[string]any{"key1": createTestValue_Slice(), "key2": createTestValue_Slice()}}, "{key1:[\"value1\",\"value2\"],key2:[\"value1\",\"value2\"]}", false},
		{"map,map", args{map[string]any{"key1": createTestValue_Map(), "key2": createTestValue_Map()}}, "{key1:{key1:\"value1\",key2:\"value2\"},key2:{key1:\"value1\",key2:\"value2\"}}", false},
		{"nil", args{nil}, "null", false},
		{"unexpected_otherType", args{map[string]any{"key1": createTestValue_Slice(), "key2": &testST1{}}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertByType(tt.args.anyValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertByType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createTestValue_Array() [2]any {
	res := [2]any{"value1", "value2"}
	return res
}

func createTestValue_Slice() []any {
	res := []any{"value1", "value2"}
	return res
}

func createTestValue_Map() map[string]any {
	res := map[string]any{"key1": "value1", "key2": "value2"}
	return res
}

func Test_Nil2EmptyFromMap(t *testing.T) {
	var valNilSlice []any
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			"値nil配列の場合空配列を返す map{ccc: []}",
			args{
				map[string]any{"key1": valNilSlice},
			},
			map[string]any{"key1": []any{}},
		},
		{
			"値無配列 map{ccc: [{}]}",
			args{
				map[string]any{"key1": []any{}},
			},
			map[string]any{"key1": []any{}},
		},
		{
			"値string map{aa:aaa1}",
			args{
				map[string]any{"key1": "value1"},
			},
			map[string]any{"key1": "value1"},
		},
		{
			"値string, 二つ map{aa:aaa1 bb:bbb1}",
			args{
				createTestValue_Map(),
			},
			createTestValue_Map(),
		},
		{
			"値string配列 map{ccc: [aaa1, bbb1]}",
			args{
				map[string]any{"key1": createTestValue_Slice()},
			},
			map[string]any{"key1": createTestValue_Slice()},
		},
		{
			"値map配列 map{ccc: [map{aa:aaa2 bb:bbb2}, map{aa:aaa3 bb:bbb3}]}",
			args{
				map[string]any{"key1": []any{createTestValue_Map(), createTestValue_Map()}},
			},
			map[string]any{"key1": []any{createTestValue_Map(), createTestValue_Map()}},
		},
		{
			"値map map{ccc: map{aa:aaa2 bb:bbb2}}",
			args{
				map[string]any{"key1": createTestValue_Map()},
			},
			map[string]any{"key1": createTestValue_Map()},
		},
		{
			"値mapでmapの値が配列 map{ccc: map{aa:[aaa1, bbb1] bb:[aaa2, bbb2]}}",
			args{
				map[string]any{"key1": map[string]any{"aaa": createTestValue_Slice(), "bbb": createTestValue_Slice()}},
			},
			map[string]any{"key1": map[string]any{"aaa": createTestValue_Slice(), "bbb": createTestValue_Slice()}},
		},
		{
			"値mapでmapの値が値nil配列 map{ccc: map{aa:[] bb:[aaa2, bbb2]}}",
			args{
				map[string]any{"key1": map[string]any{"aaa": valNilSlice, "bbb": createTestValue_Slice()}},
			},
			map[string]any{"key1": map[string]any{"aaa": []any{}, "bbb": createTestValue_Slice()}},
		},
		{
			"値mapでmapの値が値無配列 map{ccc: map{aa:[{}] bb:[aaa2, bbb2]}}",
			args{
				map[string]any{"key1": map[string]any{"aaa": []any{}, "bbb": createTestValue_Slice()}},
			},
			map[string]any{"key1": map[string]any{"aaa": []any{}, "bbb": createTestValue_Slice()}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Nil2EmptyFromMap(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nil2EmptyFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Nil2EmptyFromSlice(t *testing.T) {
	var valNilSlice []any
	type args struct {
		input []any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			"One-dimensional array with empty value",
			args{
				[]any{},
			},
			[]any{},
		},
		{
			"One-dimensional array with one string value",
			args{
				[]any{"value1"},
			},
			[]any{"value1"},
		},
		{
			"One-dimensional array with two string values",
			args{
				createTestValue_Slice(),
			},
			createTestValue_Slice(),
		},
		{
			"One-dimensional array with map values",
			args{
				[]any{createTestValue_Map(), createTestValue_Map()},
			},
			[]any{createTestValue_Map(), createTestValue_Map()},
		},
		{
			"Two-dimensional array with string values",
			args{
				[]any{createTestValue_Slice(), createTestValue_Slice()},
			},
			[]any{createTestValue_Slice(), createTestValue_Slice()},
		},
		{
			"Two-dimensional array with nil value",
			args{
				[]any{createTestValue_Slice(), valNilSlice},
			},
			[]any{createTestValue_Slice(), []any{}},
		},
		{
			"Two-dimensional array with empty value",
			args{
				[]any{createTestValue_Slice(), []any{}},
			},
			[]any{createTestValue_Slice(), []any{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Nil2EmptyFromSlice(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nil2EmptyFromSlice() = %v:%T:%t, want %v:%T:%t", got, got, reflect.ValueOf(got).IsNil(), tt.want, tt.want, reflect.ValueOf(tt.want).IsNil())
			}
		})
	}
}

func Test_nil2Empty(t *testing.T) {
	var valNilSlice []any
	type testST1 struct {
		F1 string `default:"aaa"`
	}
	type args struct {
		anyValue any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"string", args{"test string"}, "test string"},
		{"int", args{int(1)}, 1},
		{"bool", args{true}, true},
		{"array", args{createTestValue_Array()}, createTestValue_Array()},
		{"slice", args{createTestValue_Slice()}, createTestValue_Slice()},
		{"slice,lenZero", args{[]any{}}, []any{}},
		{"slice,array", args{[]any{createTestValue_Array(), createTestValue_Array()}}, []any{createTestValue_Array(), createTestValue_Array()}},
		{"slice,slice", args{[]any{createTestValue_Slice(), createTestValue_Slice()}}, []any{createTestValue_Slice(), createTestValue_Slice()}},
		{"slice,nilslice", args{[]any{createTestValue_Slice(), valNilSlice}}, []any{createTestValue_Slice(), []any{}}}, // nil
		{"slice,karaslice", args{[]any{createTestValue_Slice(), []any{}}}, []any{createTestValue_Slice(), []any{}}},    // kara
		{"slice,map", args{[]any{createTestValue_Map()}}, []any{createTestValue_Map()}},
		{"map", args{createTestValue_Map()}, createTestValue_Map()},
		{"map,array", args{map[string]any{"key1": createTestValue_Array(), "key2": createTestValue_Array()}}, map[string]any{"key1": createTestValue_Array(), "key2": createTestValue_Array()}},
		{"map,slice", args{map[string]any{"key1": createTestValue_Slice(), "key2": createTestValue_Slice()}}, map[string]any{"key1": createTestValue_Slice(), "key2": createTestValue_Slice()}},
		{"map,nilslice", args{map[string]any{"key1": valNilSlice, "key2": createTestValue_Slice()}}, map[string]any{"key1": []any{}, "key2": createTestValue_Slice()}}, // nil
		{"map,karaslice", args{map[string]any{"key1": []any{}, "key2": createTestValue_Slice()}}, map[string]any{"key1": []any{}, "key2": createTestValue_Slice()}},    // kara
		{"map,map", args{map[string]any{"key1": createTestValue_Map(), "key2": createTestValue_Map()}}, map[string]any{"key1": createTestValue_Map(), "key2": createTestValue_Map()}},
		{"unexpected_otherType", args{map[string]any{"key1": createTestValue_Slice(), "key2": &testST1{}}}, map[string]any{"key1": createTestValue_Slice(), "key2": &testST1{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nil2Empty(tt.args.anyValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nil2Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
