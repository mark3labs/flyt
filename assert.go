package flyt

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Assert provides safe type assertion with a default value
func Assert[T any](v any, defaultVal T) T {
	if v == nil {
		return defaultVal
	}
	if typed, ok := v.(T); ok {
		return typed
	}
	return defaultVal
}

// MustAssert performs type assertion and panics on failure
func MustAssert[T any](v any) T {
	if v == nil {
		panic(fmt.Sprintf("MustAssert: got nil, expected %T", *new(T)))
	}
	typed, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("MustAssert: type assertion failed, got %T, expected %T", v, *new(T)))
	}
	return typed
}

// TryAssert attempts type assertion and returns result with ok flag
func TryAssert[T any](v any) (T, bool) {
	var zero T
	if v == nil {
		return zero, false
	}
	typed, ok := v.(T)
	return typed, ok
}

// Convert attempts to convert between compatible types using JSON
func Convert[T any](v any) (T, error) {
	var result T

	// Direct type match
	if typed, ok := v.(T); ok {
		return typed, nil
	}

	// Check if types are assignable
	if v != nil {
		srcType := reflect.TypeOf(v)
		dstType := reflect.TypeOf(result)
		if srcType.AssignableTo(dstType) {
			reflect.ValueOf(&result).Elem().Set(reflect.ValueOf(v))
			return result, nil
		}
	}

	// Fallback to JSON conversion
	data, err := json.Marshal(v)
	if err != nil {
		return result, fmt.Errorf("convert: failed to marshal: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("convert: failed to unmarshal: %w", err)
	}

	return result, nil
}

// MustConvert converts between types or panics
func MustConvert[T any](v any) T {
	result, err := Convert[T](v)
	if err != nil {
		panic(fmt.Sprintf("MustConvert failed: %v", err))
	}
	return result
}

// AssertMap safely asserts to map[string]any
func AssertMap(v any) (map[string]any, bool) {
	if v == nil {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

// AssertSlice safely asserts to []any
func AssertSlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	s, ok := v.([]any)
	return s, ok
}

// AssertString safely asserts to string
func AssertString(v any) (string, bool) {
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// ChainAssert allows chaining type assertions with fallback
type ChainAssert struct {
	value  any
	result any
	ok     bool
}

// NewChainAssert creates a new assertion chain
func NewChainAssert(v any) *ChainAssert {
	return &ChainAssert{value: v}
}

// Try attempts an assertion in the chain
func (c *ChainAssert) Try(fn func(any) (any, bool)) *ChainAssert {
	if c.ok {
		return c
	}
	c.result, c.ok = fn(c.value)
	return c
}

// OrElse provides a default value if all assertions failed
func (c *ChainAssert) OrElse(defaultVal any) any {
	if c.ok {
		return c.result
	}
	return defaultVal
}

// Result returns the result and success flag
func (c *ChainAssert) Result() (any, bool) {
	return c.result, c.ok
}
