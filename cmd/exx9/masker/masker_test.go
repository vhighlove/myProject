package masker

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

// compareJSONMaps compares two JSON strings by converting them to maps and using reflect.DeepEqual
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
		err     error // Обробка очікуваної помилки
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
			wantErr: true,
			err:     errors.New("error parsing JSON: invalid character '\"' after object key:value pair"), // Очікувана помилка
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
		{
			name: "Nested keys masking",
			args: args{
				jsonData: `{"person": {"name": "John", "email": "john@example.com", "age": 30}}`,
				keys:     []string{"person/name"},
			},
			want:    `{"person":{"name":"****","email":"john@example.com","age":30}}`,
			wantErr: false,
		},
		{
			name: "Masking array elements by index",
			args: args{
				jsonData: `{"friends": ["Alice", "Bob", "Charlie"]}`,
				keys:     []string{"friends[1]"},
			},
			want:    `{"friends":["Alice","****","Charlie"]}`,
			wantErr: false,
		},
		{
			name: "Masking array range",
			args: args{
				jsonData: `{"friends": ["Alice", "Bob", "Charlie", "David", "Eva"]}`,
				keys:     []string{"friends[1:3]"},
			},
			want:    `{"friends":["Alice","****","****","David","Eva"]}`,
			wantErr: false,
		},
		{
			name: "Invalid range in array",
			args: args{
				jsonData: `{"friends": ["Alice", "Bob", "Charlie"]}`,
				keys:     []string{"friends[3:5]"},
			},
			wantErr: true,
			err:     errors.New("error processing map: invalid range"),
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
				if tt.err != nil && err.Error() != tt.err.Error() {
					t.Errorf("Expected error = %v, got error = %v", tt.err, err) // Обробка: звірка помилок на ідентичність
				}
				return
			}
			if !compareJSONMaps(got, tt.want) {
				t.Errorf("MaskByKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}
