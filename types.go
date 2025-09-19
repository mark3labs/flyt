package flyt

// PrepResult represents common prep phase results
type PrepResult struct {
	Data   any
	Config map[string]any
	Items  []any
}

// ExecResult represents common exec phase results
type ExecResult struct {
	Success bool
	Data    any
	Error   string
	Results []any
}

// BatchPrepResult for batch processing nodes
type BatchPrepResult struct {
	Items     []any
	BatchSize int
	Workers   int
	CountKey  string
}

// BatchExecResult for batch processing results
type BatchExecResult struct {
	Results    []any
	Errors     []error
	Successful int
	Failed     int
}

// ExtractPrepResult safely extracts PrepResult from any
func ExtractPrepResult(v any) (*PrepResult, bool) {
	if v == nil {
		return nil, false
	}
	if result, ok := v.(*PrepResult); ok {
		return result, true
	}
	if result, ok := v.(PrepResult); ok {
		return &result, true
	}
	return nil, false
}

// ExtractExecResult safely extracts ExecResult from any
func ExtractExecResult(v any) (*ExecResult, bool) {
	if v == nil {
		return nil, false
	}
	if result, ok := v.(*ExecResult); ok {
		return result, true
	}
	if result, ok := v.(ExecResult); ok {
		return &result, true
	}
	return nil, false
}

// MustExtractPrepResult extracts PrepResult or panics
func MustExtractPrepResult(v any) *PrepResult {
	result, ok := ExtractPrepResult(v)
	if !ok {
		panic("failed to extract PrepResult")
	}
	return result
}

// MustExtractExecResult extracts ExecResult or panics
func MustExtractExecResult(v any) *ExecResult {
	result, ok := ExtractExecResult(v)
	if !ok {
		panic("failed to extract ExecResult")
	}
	return result
}
