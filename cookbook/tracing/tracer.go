package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bytedance/sonic"
	langfuse "github.com/cloudwego/eino-ext/libs/acl/langfuse"
	"github.com/google/uuid"
	"github.com/mark3labs/flyt"
)

// Tracer wraps Langfuse client for tracing Flyt workflows
type Tracer struct {
	client      langfuse.Langfuse
	enabled     bool
	traceID     string
	currentSpan *Span
	spans       []*Span
}

// Span represents a traced span of execution
type Span struct {
	tracer    *Tracer
	name      string
	startTime time.Time
	metadata  map[string]any
	spanID    string
	traceID   string
	parentID  string
	input     any
}

// Trace represents the main trace
type Trace struct {
	tracer    *Tracer
	name      string
	startTime time.Time
	metadata  map[string]any
}

// NewTracer creates a new tracer instance
func NewTracer() *Tracer {
	// Check if Langfuse is configured
	host := os.Getenv("LANGFUSE_HOST")
	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")

	enabled := host != "" && publicKey != "" && secretKey != ""

	var client langfuse.Langfuse
	if enabled {
		// Create Langfuse client with the new API
		client = langfuse.NewLangfuse(
			host,
			publicKey,
			secretKey,
			// You can add options here like:
			// langfuse.WithFlushInterval(5*time.Second),
			// langfuse.WithMaxRetry(3),
		)
	}

	return &Tracer{
		client:  client,
		enabled: enabled,
		spans:   make([]*Span, 0),
	}
}

// StartTrace starts a new trace
func (t *Tracer) StartTrace(name string, metadata map[string]any) *Trace {
	trace := &Trace{
		tracer:    t,
		name:      name,
		startTime: time.Now(),
		metadata:  metadata,
	}

	if t.enabled && t.client != nil {
		// Generate trace ID
		t.traceID = uuid.New().String()

		// Marshal metadata to JSON string if needed
		metadataJSON, _ := sonic.MarshalString(metadata)

		// Create trace with new API
		traceBody := &langfuse.TraceEventBody{
			BaseEventBody: langfuse.BaseEventBody{
				ID:       t.traceID,
				Name:     name,
				MetaData: metadata,
			},
			TimeStamp: time.Now(),
			Input:     metadataJSON,
		}

		_, err := t.client.CreateTrace(traceBody)
		if err != nil {
			log.Printf("Failed to create trace: %v", err)
		}
	} else {
		log.Printf("🔍 [TRACE START] %s - %v", name, metadata)
	}

	return trace
}

// End ends the trace
func (tr *Trace) End(err error) {
	duration := time.Since(tr.startTime)

	if tr.tracer.enabled && tr.tracer.client != nil {
		// Update trace metadata
		metadata := tr.metadata
		if metadata == nil {
			metadata = make(map[string]any)
		}
		metadata["duration_ms"] = duration.Milliseconds()
		if err != nil {
			metadata["error"] = err.Error()
		}

		// Create a completion event to mark the trace end
		eventName := tr.name + "_completed"
		eventBody := &langfuse.EventEventBody{
			BaseObservationEventBody: langfuse.BaseObservationEventBody{
				BaseEventBody: langfuse.BaseEventBody{
					Name:     eventName,
					MetaData: metadata,
				},
				TraceID:   tr.tracer.traceID,
				StartTime: time.Now(),
			},
		}

		_, eventErr := tr.tracer.client.CreateEvent(eventBody)
		if eventErr != nil {
			log.Printf("Failed to end trace: %v", eventErr)
		}
	} else {
		status := "SUCCESS"
		if err != nil {
			status = fmt.Sprintf("ERROR: %v", err)
		}
		log.Printf("🔍 [TRACE END] %s - Duration: %v - Status: %s", tr.name, duration, status)
	}
}

// AddMetadata adds metadata to the trace
func (tr *Trace) AddMetadata(metadata map[string]any) {
	if tr.metadata == nil {
		tr.metadata = make(map[string]any)
	}
	for k, v := range metadata {
		tr.metadata[k] = v
	}
}

// StartSpan starts a new span with input tracking
func (t *Tracer) StartSpan(name string, metadata map[string]any, input any) *Span {
	span := &Span{
		tracer:    t,
		name:      name,
		startTime: time.Now(),
		metadata:  metadata,
		spanID:    uuid.New().String(),
		traceID:   t.traceID,
		input:     input,
	}

	if t.currentSpan != nil {
		span.parentID = t.currentSpan.spanID
	}

	t.currentSpan = span
	t.spans = append(t.spans, span)

	if t.enabled && t.client != nil {
		// Convert input to JSON string
		inputJSON, _ := sonic.MarshalString(input)

		// Create span with new API
		spanBody := &langfuse.SpanEventBody{
			BaseObservationEventBody: langfuse.BaseObservationEventBody{
				BaseEventBody: langfuse.BaseEventBody{
					ID:       span.spanID,
					Name:     name,
					MetaData: metadata,
				},
				TraceID:             t.traceID,
				ParentObservationID: span.parentID,
				Input:               inputJSON,
				StartTime:           time.Now(),
			},
		}

		_, err := t.client.CreateSpan(spanBody)
		if err != nil {
			log.Printf("Failed to create span: %v", err)
		}
	} else {
		log.Printf("  📍 [SPAN START] %s - metadata: %v, input: %v", name, metadata, input)
	}

	return span
}

// EndWithOutput ends the span with output tracking
func (s *Span) EndWithOutput(output any, err error) {
	duration := time.Since(s.startTime)

	if s.tracer.enabled && s.tracer.client != nil {
		// Update span metadata
		metadata := s.metadata
		if metadata == nil {
			metadata = make(map[string]any)
		}
		metadata["duration_ms"] = duration.Milliseconds()
		if err != nil {
			metadata["error"] = err.Error()
		}

		// Convert output to JSON string
		outputJSON, _ := sonic.MarshalString(output)

		// End the span with new API
		spanBody := &langfuse.SpanEventBody{
			BaseObservationEventBody: langfuse.BaseObservationEventBody{
				BaseEventBody: langfuse.BaseEventBody{
					ID:       s.spanID,
					MetaData: metadata,
				},
				TraceID: s.traceID,
				Output:  outputJSON,
			},
			EndTime: time.Now(),
		}

		endErr := s.tracer.client.EndSpan(spanBody)
		if endErr != nil {
			log.Printf("Failed to end span: %v", endErr)
		}
	} else {
		status := "SUCCESS"
		if err != nil {
			status = fmt.Sprintf("ERROR: %v", err)
		}
		log.Printf("  📍 [SPAN END] %s - Duration: %v - Status: %s - Output: %v", s.name, duration, status, output)
	}

	// Reset current span to parent if this was the current span
	if s.tracer.currentSpan == s {
		// Find parent span
		var parent *Span
		for _, span := range s.tracer.spans {
			if span.spanID == s.parentID {
				parent = span
				break
			}
		}
		s.tracer.currentSpan = parent
	}
}

// End ends the span (calls EndWithOutput with nil output)
func (s *Span) End(err error) {
	s.EndWithOutput(nil, err)
}

// AddMetadata adds metadata to the span
func (s *Span) AddMetadata(metadata map[string]any) {
	if s.metadata == nil {
		s.metadata = make(map[string]any)
	}
	for k, v := range metadata {
		s.metadata[k] = v
	}
}

// Flush flushes all pending traces
func (t *Tracer) Flush(ctx context.Context) {
	if t.enabled && t.client != nil {
		t.client.Flush()
		log.Println("📊 Traces flushed to Langfuse")
	}
}

// marshalSharedStore converts SharedStore to a JSON-serializable map
func marshalSharedStore(store *flyt.SharedStore) map[string]any {
	result := make(map[string]any)
	// Note: SharedStore doesn't expose a method to iterate all keys,
	// so we'd need to know the keys in advance or modify the implementation
	// For now, we'll just return a placeholder
	result["type"] = "SharedStore"
	return result
}

// toJSONString safely converts any value to a JSON string
func toJSONString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}
