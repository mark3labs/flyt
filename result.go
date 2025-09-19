package flyt

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Result is a wrapper type that provides convenient type assertion methods.
// It's used as return type for Prep and Exec methods to improve developer experience.
// It can hold any value and provides type-safe accessors.
type Result struct {
	value any
}

// NewResult creates a new Result from any value.
func NewResult(v any) Result {
	return Result{value: v}
}

// R is a shorthand for NewResult.
func R(v any) Result {
	return NewResult(v)
}

// AsString retrieves the Result as a string.
// Returns empty string and false if the Result value is nil or not a string.
func (r Result) AsString() (string, bool) {
	if r.value == nil {
		return "", false
	}
	s, ok := r.value.(string)
	return s, ok
}

// AsStringOr retrieves the Result as a string.
// Returns the provided default value if the Result is nil or not a string.
func (r Result) AsStringOr(defaultVal string) string {
	s, ok := r.AsString()
	if !ok {
		return defaultVal
	}
	return s
}

// MustString retrieves the Result as a string.
// Panics if the Result is nil or not a string.
func (r Result) MustString() string {
	s, ok := r.AsString()
	if !ok {
		panic(fmt.Sprintf("Result.MustString: value is not a string (type %T)", r))
	}
	return s
}

// AsInt retrieves the Result as an int.
// Returns 0 and false if the Result value is nil or cannot be converted to int.
// Supports conversion from int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, and float types.
func (r Result) AsInt() (int, bool) {
	if r.value == nil {
		return 0, false
	}

	switch v := r.value.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// AsIntOr retrieves the Result as an int.
// Returns the provided default value if the Result is nil or cannot be converted to int.
func (r Result) AsIntOr(defaultVal int) int {
	i, ok := r.AsInt()
	if !ok {
		return defaultVal
	}
	return i
}

// MustInt retrieves the Result as an int.
// Panics if the Result is nil or cannot be converted to int.
func (r Result) MustInt() int {
	i, ok := r.AsInt()
	if !ok {
		panic(fmt.Sprintf("Result.MustInt: value cannot be converted to int (type %T)", r))
	}
	return i
}

// AsFloat64 retrieves the Result as a float64.
// Returns 0.0 and false if the Result value is nil or cannot be converted to float64.
// Supports conversion from various numeric types.
func (r Result) AsFloat64() (float64, bool) {
	if r.value == nil {
		return 0.0, false
	}

	switch v := r.value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0.0, false
	}
}

// AsFloat64Or retrieves the Result as a float64.
// Returns the provided default value if the Result is nil or cannot be converted to float64.
func (r Result) AsFloat64Or(defaultVal float64) float64 {
	f, ok := r.AsFloat64()
	if !ok {
		return defaultVal
	}
	return f
}

// MustFloat64 retrieves the Result as a float64.
// Panics if the Result is nil or cannot be converted to float64.
func (r Result) MustFloat64() float64 {
	f, ok := r.AsFloat64()
	if !ok {
		panic(fmt.Sprintf("Result.MustFloat64: value cannot be converted to float64 (type %T)", r))
	}
	return f
}

// AsBool retrieves the Result as a bool.
// Returns false and false if the Result value is nil or not a bool.
func (r Result) AsBool() (bool, bool) {
	if r.value == nil {
		return false, false
	}
	b, ok := r.value.(bool)
	return b, ok
}

// AsBoolOr retrieves the Result as a bool.
// Returns the provided default value if the Result is nil or not a bool.
func (r Result) AsBoolOr(defaultVal bool) bool {
	b, ok := r.AsBool()
	if !ok {
		return defaultVal
	}
	return b
}

// MustBool retrieves the Result as a bool.
// Panics if the Result is nil or not a bool.
func (r Result) MustBool() bool {
	b, ok := r.AsBool()
	if !ok {
		panic(fmt.Sprintf("Result.MustBool: value is not a bool (type %T)", r))
	}
	return b
}

// AsSlice retrieves the Result as a []any slice.
// Returns nil and false if the Result value is nil or not a slice.
// Uses ToSlice to convert various slice types to []any.
func (r Result) AsSlice() ([]any, bool) {
	if r.value == nil {
		return nil, false
	}

	// Check if it's already []any
	if slice, ok := r.value.([]any); ok {
		return slice, true
	}

	// Try to convert using ToSlice
	result := ToSlice(r.value)
	// ToSlice wraps non-slice values, so check if it's actually a slice
	if len(result) == 1 && result[0] == r.value {
		// ToSlice wrapped a non-slice value
		return nil, false
	}
	return result, true
}

// AsSliceOr retrieves the Result as a []any slice.
// Returns the provided default value if the Result is nil or not a slice.
func (r Result) AsSliceOr(defaultVal []any) []any {
	s, ok := r.AsSlice()
	if !ok {
		return defaultVal
	}
	return s
}

// MustSlice retrieves the Result as a []any slice.
// Panics if the Result is nil or not a slice.
func (r Result) MustSlice() []any {
	s, ok := r.AsSlice()
	if !ok {
		panic(fmt.Sprintf("Result.MustSlice: value is not a slice (type %T)", r))
	}
	return s
}

// AsMap retrieves the Result as a map[string]any.
// Returns nil and false if the Result value is nil or not a map[string]any.
func (r Result) AsMap() (map[string]any, bool) {
	if r.value == nil {
		return nil, false
	}
	m, ok := r.value.(map[string]any)
	return m, ok
}

// AsMapOr retrieves the Result as a map[string]any.
// Returns the provided default value if the Result is nil or not a map[string]any.
func (r Result) AsMapOr(defaultVal map[string]any) map[string]any {
	m, ok := r.AsMap()
	if !ok {
		return defaultVal
	}
	return m
}

// MustMap retrieves the Result as a map[string]any.
// Panics if the Result is nil or not a map[string]any.
func (r Result) MustMap() map[string]any {
	m, ok := r.AsMap()
	if !ok {
		panic(fmt.Sprintf("Result.MustMap: value is not a map[string]any (type %T)", r))
	}
	return m
}

// Bind binds the Result to a struct using JSON marshaling/unmarshaling.
// This allows for easy conversion of complex types.
// The destination must be a pointer to the target struct.
// Returns an error if the Result value is nil or binding fails.
//
// Example:
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//	var user User
//	err := result.Bind(&user)
func (r Result) Bind(dest any) error {
	if r.value == nil {
		return fmt.Errorf("cannot bind nil Result value")
	}

	// Check if dest is a pointer
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	// If Result value is already the correct type, assign directly
	valType := reflect.TypeOf(r.value)
	destType := rv.Type().Elem()
	if valType == destType {
		rv.Elem().Set(reflect.ValueOf(r.value))
		return nil
	}

	// Otherwise use JSON as intermediate format
	jsonBytes, err := json.Marshal(r.value)
	if err != nil {
		return fmt.Errorf("failed to marshal Result: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %w", err)
	}

	return nil
}

// MustBind is like Bind but panics if binding fails.
// Use this only when binding failure should be considered a programming error.
//
// Example:
//
//	var config Config
//	result.MustBind(&config)  // Panics if binding fails
func (r Result) MustBind(dest any) {
	if err := r.Bind(dest); err != nil {
		panic(fmt.Sprintf("Result.MustBind failed: %v", err))
	}
}

// IsNil checks if the Result value is nil.
func (r Result) IsNil() bool {
	return r.value == nil
}

// Type returns the underlying type of the Result value as a string.
// Returns "nil" if the Result value is nil.
func (r Result) Type() string {
	if r.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%T", r.value)
}

// Value returns the underlying value as any.
// This is useful when you need to pass the Result to functions expecting any.
func (r Result) Value() any {
	return r.value
}

// As attempts to retrieve the Result as the specified type T.
// Returns the typed value and true if successful, or zero value and false if not.
//
// Example:
//
//	if user, ok := result.As[*User](); ok {
//	    // Use typed user
//	}
func As[T any](r Result) (T, bool) {
	var zero T
	if r.value == nil {
		return zero, false
	}
	typed, ok := r.value.(T)
	return typed, ok
}

// MustAs retrieves the Result as the specified type T.
// Panics if the Result is nil or not of type T.
//
// Example:
//
//	user := MustAs[*User](result)
func MustAs[T any](r Result) T {
	typed, ok := As[T](r)
	if !ok {
		panic(fmt.Sprintf("Result.MustAs: value is not of type %T", *new(T)))
	}
	return typed
}
