package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/mark3labs/flyt"
)

// Tracer wraps Langfuse client for tracing Flyt workflows
type Tracer struct {
	client      *langfuse.Langfuse
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

	var client *langfuse.Langfuse
	if enabled {
		ctx := context.Background()
		client = langfuse.New(ctx)
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

		// Create Langfuse trace
		now := time.Now()
		_, err := t.client.Trace(&model.Trace{
			ID:        t.traceID,
			Name:      name,
			Metadata:  convertToModelM(metadata),
			Timestamp: &now,
		})
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
		// Update trace with end time and duration
		metadata := tr.metadata
		if metadata == nil {
			metadata = make(map[string]any)
		}
		metadata["duration_ms"] = duration.Milliseconds()
		if err != nil {
			metadata["error"] = err.Error()
		}

		// Add an event to mark the end
		eventName := tr.name + "_completed"
		now := time.Now()
		_, eventErr := tr.tracer.client.Event(&model.Event{
			TraceID:   tr.tracer.traceID,
			Name:      eventName,
			Metadata:  convertToModelM(metadata),
			StartTime: &now,
		}, nil)
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

// StartSpan starts a new span
func (t *Tracer) StartSpan(name string, metadata map[string]any) *Span {
	span := &Span{
		tracer:    t,
		name:      name,
		startTime: time.Now(),
		metadata:  metadata,
		spanID:    uuid.New().String(),
		traceID:   t.traceID,
	}

	if t.currentSpan != nil {
		span.parentID = t.currentSpan.spanID
	}

	t.currentSpan = span
	t.spans = append(t.spans, span)

	if t.enabled && t.client != nil {
		// Create Langfuse span
		now := time.Now()
		langfuseSpan := &model.Span{
			ID:        span.spanID,
			TraceID:   t.traceID,
			Name:      name,
			Metadata:  convertToModelM(metadata),
			StartTime: &now,
		}

		if span.parentID != "" {
			langfuseSpan.ParentObservationID = span.parentID
		}

		var parentIDPtr *string
		if span.parentID != "" {
			parentIDPtr = &span.parentID
		}
		_, err := t.client.Span(langfuseSpan, parentIDPtr)
		if err != nil {
			log.Printf("Failed to create span: %v", err)
		}
	} else {
		log.Printf("  📍 [SPAN START] %s - %v", name, metadata)
	}

	return span
}

// End ends the span
func (s *Span) End(err error) {
	duration := time.Since(s.startTime)

	if s.tracer.enabled && s.tracer.client != nil {
		// End the span
		metadata := s.metadata
		if metadata == nil {
			metadata = make(map[string]any)
		}
		metadata["duration_ms"] = duration.Milliseconds()
		if err != nil {
			metadata["error"] = err.Error()
		}

		now := time.Now()
		_, endErr := s.tracer.client.SpanEnd(&model.Span{
			ID:       s.spanID,
			TraceID:  s.traceID,
			EndTime:  &now,
			Metadata: convertToModelM(metadata),
		})
		if endErr != nil {
			log.Printf("Failed to end span: %v", endErr)
		}
	} else {
		status := "SUCCESS"
		if err != nil {
			status = fmt.Sprintf("ERROR: %v", err)
		}
		log.Printf("  📍 [SPAN END] %s - Duration: %v - Status: %s", s.name, duration, status)
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
		t.client.Flush(ctx)
		log.Println("📊 Traces flushed to Langfuse")
	}
}

// convertToModelM converts map[string]any to model.M
func convertToModelM(m map[string]any) model.M {
	if m == nil {
		return nil
	}

	result := make(model.M)
	for k, v := range m {
		// Convert the value to a string representation if needed
		switch val := v.(type) {
		case string:
			result[k] = val
		case int, int32, int64, float32, float64, bool:
			result[k] = fmt.Sprintf("%v", val)
		case *flyt.SharedStore:
			// Handle SharedStore specially
			result[k] = "SharedStore{...}"
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}
	return result
}
