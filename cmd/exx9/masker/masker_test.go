package masker

import (
	"encoding/json"
	"reflect"
	"testing"
)

func compareJSONMaps(got, want string) bool {
	var gotMap, wantMap map[string]interface{}

	if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(want), &wantMap); err != nil {
		return false
	}

	return reflect.DeepEqual(gotMap, wantMap)
}

func TestMaskByKeys(t *testing.T) {
	type args struct {
		jsonData string
		keys     []string
		strategy []maskingStrategy
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid input with one key to mask",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				keys:     []string{"email"},
			},
			want:    `{"name":"John","email":"****","age":30}`,
			wantErr: false,
		},
		{
			name: "Valid input with no keys to mask",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				keys:     []string{},
			},
			want:    `{"name":"John","email":"john@example.com","age":30}`,
			wantErr: false,
		},
		{
			name: "Invalid JSON input",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com" "age": 30}`,
				keys:     []string{"email"},
			},
			want:    `error during masking by keys: error parsing JSON: invalid character '"' after object key:value pair`,
			wantErr: true,
		},
		{
			name: "Custom masking strategy",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				keys:     []string{"age"},
				strategy: []maskingStrategy{func(_ interface{}) interface{} { return "XX" }},
			},
			want:    `{"name":"John","email":"john@example.com","age":"XX"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaskByKeys(tt.args.jsonData, tt.args.keys, tt.args.strategy...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskByKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if (err != nil) && !compareJSONMaps(got, tt.want) {
				t.Errorf("MaskByKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskByIndexes(t *testing.T) {
	type args struct {
		jsonData string
		indexes  []int
		strategy []maskingStrategy
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid input with one index to mask",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				indexes:  []int{1},
			},
			want:    `{"name":"John","email":"****","age":30}`,
			wantErr: false,
		},
		{
			name: "Valid input with multiple indexes to mask",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				indexes:  []int{0, 2},
			},
			want:    `{"name":"****","email":"john@example.com","age":"****"}`,
			wantErr: false,
		},
		{
			name: "Index out of bounds",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				indexes:  []int{3},
			},
			want:    "index 3 is out of bounds for keys",
			wantErr: true,
		},
		{
			name: "Custom masking strategy",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com", "age": 30}`,
				indexes:  []int{2},
				strategy: []maskingStrategy{func(_ interface{}) interface{} { return "XX" }},
			},
			want:    `{"name":"John","email":"john@example.com","age":"XX"}`,
			wantErr: false,
		},
		{
			name: "Invalid JSON input",
			args: args{
				jsonData: `{"name": "John", "email": "john@example.com" "age": 30}`,
				indexes:  []int{1},
			},
			want:    `error during masking by indexes: error parsing JSON: invalid character '"' after object key:value pair`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaskByIndexes(tt.args.jsonData, tt.args.indexes, tt.args.strategy...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskByIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !compareJSONMaps(got, tt.want) {
				t.Errorf("MaskByIndexes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
