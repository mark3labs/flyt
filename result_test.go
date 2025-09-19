package flyt_test

import (
	"encoding/json"
	"testing"

	"github.com/mark3labs/flyt"
)

func TestResultString(t *testing.T) {
	tests := []struct {
		name      string
		result    flyt.Result
		wantStr   string
		wantOk    bool
		wantOr    string
		orDefault string
	}{
		{
			name:      "valid string",
			result:    flyt.R("hello"),
			wantStr:   "hello",
			wantOk:    true,
			wantOr:    "hello",
			orDefault: "default",
		},
		{
			name:      "nil value",
			result:    flyt.R(nil),
			wantStr:   "",
			wantOk:    false,
			wantOr:    "default",
			orDefault: "default",
		},
		{
			name:      "non-string value",
			result:    flyt.R(123),
			wantStr:   "",
			wantOk:    false,
			wantOr:    "default",
			orDefault: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test AsString
			str, ok := tt.result.AsString()
			if str != tt.wantStr || ok != tt.wantOk {
				t.Errorf("AsString() = (%q, %v), want (%q, %v)", str, ok, tt.wantStr, tt.wantOk)
			}

			// Test AsStringOr
			strOr := tt.result.AsStringOr(tt.orDefault)
			if strOr != tt.wantOr {
				t.Errorf("AsStringOr(%q) = %q, want %q", tt.orDefault, strOr, tt.wantOr)
			}

			// Test MustString for valid strings
			if tt.wantOk {
				mustStr := tt.result.MustString()
				if mustStr != tt.wantStr {
					t.Errorf("MustString() = %q, want %q", mustStr, tt.wantStr)
				}
			} else {
				// Should panic for non-strings
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustString() should have panicked")
					}
				}()
				tt.result.MustString()
			}
		})
	}
}

func TestResultInt(t *testing.T) {
	tests := []struct {
		name    string
		result  flyt.Result
		wantInt int
		wantOk  bool
	}{
		{"int", flyt.R(42), 42, true},
		{"int8", flyt.R(int8(42)), 42, true},
		{"int16", flyt.R(int16(42)), 42, true},
		{"int32", flyt.R(int32(42)), 42, true},
		{"int64", flyt.R(int64(42)), 42, true},
		{"uint", flyt.R(uint(42)), 42, true},
		{"uint8", flyt.R(uint8(42)), 42, true},
		{"uint16", flyt.R(uint16(42)), 42, true},
		{"uint32", flyt.R(uint32(42)), 42, true},
		{"uint64", flyt.R(uint64(42)), 42, true},
		{"float32", flyt.R(float32(42.0)), 42, true},
		{"float64", flyt.R(float64(42.0)), 42, true},
		{"nil", flyt.R(nil), 0, false},
		{"string", flyt.R("42"), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, ok := tt.result.AsInt()
			if i != tt.wantInt || ok != tt.wantOk {
				t.Errorf("AsInt() = (%d, %v), want (%d, %v)", i, ok, tt.wantInt, tt.wantOk)
			}

			iOr := tt.result.AsIntOr(99)
			wantOr := tt.wantInt
			if !tt.wantOk {
				wantOr = 99
			}
			if iOr != wantOr {
				t.Errorf("AsIntOr(99) = %d, want %d", iOr, wantOr)
			}
		})
	}
}

func TestResultFloat64(t *testing.T) {
	tests := []struct {
		name    string
		result  flyt.Result
		wantVal float64
		wantOk  bool
	}{
		{"float64", flyt.R(42.5), 42.5, true},
		{"float32", flyt.R(float32(42.5)), 42.5, true},
		{"int", flyt.R(42), 42.0, true},
		{"uint", flyt.R(uint(42)), 42.0, true},
		{"nil", flyt.R(nil), 0.0, false},
		{"string", flyt.R("42.5"), 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, ok := tt.result.AsFloat64()
			if f != tt.wantVal || ok != tt.wantOk {
				t.Errorf("AsFloat64() = (%f, %v), want (%f, %v)", f, ok, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestResultBool(t *testing.T) {
	tests := []struct {
		name    string
		result  flyt.Result
		wantVal bool
		wantOk  bool
	}{
		{"true", flyt.R(true), true, true},
		{"false", flyt.R(false), false, true},
		{"nil", flyt.R(nil), false, false},
		{"string", flyt.R("true"), false, false},
		{"int", flyt.R(1), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, ok := tt.result.AsBool()
			if b != tt.wantVal || ok != tt.wantOk {
				t.Errorf("AsBool() = (%v, %v), want (%v, %v)", b, ok, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestResultSlice(t *testing.T) {
	tests := []struct {
		name    string
		result  flyt.Result
		wantLen int
		wantOk  bool
	}{
		{"[]any", flyt.R([]any{1, 2, 3}), 3, true},
		{"[]string", flyt.R([]string{"a", "b", "c"}), 3, true},
		{"[]int", flyt.R([]int{1, 2, 3}), 3, true},
		{"nil", flyt.R(nil), 0, false},
		{"non-slice", flyt.R("not a slice"), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ok := tt.result.AsSlice()
			if ok != tt.wantOk {
				t.Errorf("AsSlice() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && len(s) != tt.wantLen {
				t.Errorf("AsSlice() len = %d, want %d", len(s), tt.wantLen)
			}
		})
	}
}

func TestResultMap(t *testing.T) {
	tests := []struct {
		name   string
		result flyt.Result
		wantOk bool
	}{
		{"map[string]any", flyt.R(map[string]any{"key": "value"}), true},
		{"nil", flyt.R(nil), false},
		{"non-map", flyt.R("not a map"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := tt.result.AsMap()
			if ok != tt.wantOk {
				t.Errorf("AsMap() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestResultBind(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		result  flyt.Result
		wantErr bool
	}{
		{
			name:    "valid map",
			result:  flyt.R(map[string]any{"id": 123, "name": "Alice"}),
			wantErr: false,
		},
		{
			name:    "valid struct",
			result:  flyt.R(User{ID: 456, Name: "Bob"}),
			wantErr: false,
		},
		{
			name:    "nil value",
			result:  flyt.R(nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user User
			err := tt.result.Bind(&user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && user.ID == 0 {
				t.Error("Bind() didn't populate the struct")
			}
		})
	}

	// Test Bind with non-pointer
	t.Run("non-pointer dest", func(t *testing.T) {
		var user User
		result := flyt.R(map[string]any{"id": 123, "name": "Alice"})
		err := result.Bind(user) // Pass by value, not pointer
		if err == nil {
			t.Error("Bind() should error with non-pointer destination")
		}
	})
}

func TestResultMustBind(t *testing.T) {
	type Config struct {
		Timeout int `json:"timeout"`
	}

	// Test successful bind
	result := flyt.R(map[string]any{"timeout": 30})
	var config Config
	result.MustBind(&config)
	if config.Timeout != 30 {
		t.Errorf("MustBind() didn't populate struct correctly: %+v", config)
	}

	// Test panic on failure
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBind() should have panicked with nil value")
		}
	}()
	nilResult := flyt.R(nil)
	nilResult.MustBind(&config)
}

func TestResultHelpers(t *testing.T) {
	// Test IsNil
	nilResult := flyt.R(nil)
	if !nilResult.IsNil() {
		t.Error("IsNil() should return true for nil value")
	}

	nonNilResult := flyt.R("value")
	if nonNilResult.IsNil() {
		t.Error("IsNil() should return false for non-nil value")
	}

	// Test Type
	if nilResult.Type() != "nil" {
		t.Errorf("Type() = %q, want 'nil'", nilResult.Type())
	}

	stringResult := flyt.R("hello")
	expectedType := "string"
	if stringResult.Type() != expectedType {
		t.Errorf("Type() = %q, want %q", stringResult.Type(), expectedType)
	}

	// Test Value
	val := stringResult.Value()
	if val != "hello" {
		t.Errorf("Value() = %v, want 'hello'", val)
	}
}

func TestResultAs(t *testing.T) {
	type CustomType struct {
		Value string
	}

	custom := &CustomType{Value: "test"}
	result := flyt.R(custom)

	// Test successful type assertion
	if typed, ok := flyt.As[*CustomType](result); !ok || typed.Value != "test" {
		t.Errorf("As[*CustomType]() = (%v, %v), want (%v, true)", typed, ok, custom)
	}

	// Test failed type assertion
	if _, ok := flyt.As[string](result); ok {
		t.Error("As[string]() should have failed for *CustomType")
	}

	// Test with nil
	nilResult := flyt.R(nil)
	if _, ok := flyt.As[*CustomType](nilResult); ok {
		t.Error("As[*CustomType]() should have failed for nil value")
	}
}

func TestResultMustAs(t *testing.T) {
	result := flyt.R("hello")

	// Test successful assertion
	str := flyt.MustAs[string](result)
	if str != "hello" {
		t.Errorf("MustAs[string]() = %q, want 'hello'", str)
	}

	// Test panic on failure
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAs[int]() should have panicked")
		}
	}()
	flyt.MustAs[int](result)
}

func TestResultIntegration(t *testing.T) {
	// Simulate a complex result from a node
	complexData := map[string]any{
		"user": map[string]any{
			"id":   123,
			"name": "Alice",
		},
		"items": []any{"item1", "item2", "item3"},
		"config": map[string]any{
			"enabled": true,
			"timeout": 30,
		},
	}

	result := flyt.R(complexData)

	// Test various accessors
	m, ok := result.AsMap()
	if !ok {
		t.Fatal("AsMap() failed")
	}

	// Access nested data
	if user, ok := m["user"].(map[string]any); ok {
		if user["name"] != "Alice" {
			t.Errorf("user name = %v, want 'Alice'", user["name"])
		}
	} else {
		t.Error("Failed to access user map")
	}

	// Bind to struct
	type Response struct {
		User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Items  []string `json:"items"`
		Config struct {
			Enabled bool `json:"enabled"`
			Timeout int  `json:"timeout"`
		} `json:"config"`
	}

	var response Response
	err := result.Bind(&response)
	if err != nil {
		t.Fatalf("Bind() failed: %v", err)
	}

	if response.User.ID != 123 {
		t.Errorf("User.ID = %d, want 123", response.User.ID)
	}
	if len(response.Items) != 3 {
		t.Errorf("len(Items) = %d, want 3", len(response.Items))
	}
	if !response.Config.Enabled {
		t.Error("Config.Enabled should be true")
	}
}

func TestResultWithJSON(t *testing.T) {
	// Test that Result works with JSON marshaling/unmarshaling
	original := map[string]any{
		"key": "value",
		"num": 42,
	}

	result := flyt.R(original)

	// Marshal the underlying value
	data, err := json.Marshal(result.Value())
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Create new Result and verify
	newResult := flyt.R(decoded)
	m, ok := newResult.AsMap()
	if !ok {
		t.Fatal("AsMap() failed after JSON round-trip")
	}

	if m["key"] != "value" || m["num"].(float64) != 42 {
		t.Errorf("JSON round-trip changed data: %v", m)
	}
}
