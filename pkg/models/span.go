package models

import (
	"time"
)

// Span represents a single trace span
type Span struct {
	TraceID       string            `json:"trace_id"`
	SpanID        string            `json:"span_id"`
	ParentSpanID  string            `json:"parent_span_id,omitempty"`
	ServiceName   string            `json:"service_name"`
	OperationName string            `json:"operation_name"`
	StartTime     time.Time         `json:"start_time"`
	Duration      time.Duration     `json:"duration"`
	Tags          map[string]string `json:"tags"`
	Status        SpanStatus        `json:"status"`

	// HTTP-specific fields
	HTTPMethod string `json:"http_method,omitempty"`
	HTTPPath   string `json:"http_path,omitempty"`
	HTTPStatus int    `json:"http_status,omitempty"`

	// Kubernetes-specific
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
}

// SpanStatus represents the status of a span
type SpanStatus struct {
	Code    SpanStatusCode `json:"code"`
	Message string         `json:"message,omitempty"`
}

// SpanStatusCode represents span status codes
type SpanStatusCode int

const (
	StatusCodeOK SpanStatusCode = iota
	StatusCodeError
	StatusCodeTimeout
)

// Trace represents a collection of spans that form a trace
type Trace struct {
	TraceID   string        `json:"trace_id"`
	Spans     []*Span       `json:"spans"`
	Duration  time.Duration `json:"duration"`
	Services  []string      `json:"services"`
	StartTime time.Time     `json:"start_time"`
	Status    TraceStatus   `json:"status"`
}

// TraceStatus represents the overall status of a trace
type TraceStatus struct {
	Code    TraceStatusCode `json:"code"`
	Message string          `json:"message,omitempty"`
}

// TraceStatusCode represents trace status codes
type TraceStatusCode int

const (
	TraceStatusOK TraceStatusCode = iota
	TraceStatusError
	TraceStatusPartial
)