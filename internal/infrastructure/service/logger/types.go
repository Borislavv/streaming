package logger

import (
	"time"
)

const (
	InfoLevelReadable      = "INFO"
	DebugLevelReadable     = "DEBUG"
	WarningLevelReadable   = "WARNING"
	ErrorLevelReadable     = "ERROR"
	CriticalLevelReadable  = "CRITICAL"
	EmergencyLevelReadable = "EMERGENCY"
	InfoLevel              = iota
	DebugLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
	EmergencyLevel

	DebugLogType = "debug"
	InfoLogType  = "info"
	ErrorLogType = "error"
)

type LoggableError interface {
	Error() string
	Level() int
}

type RequestIdAware interface {
	RequestID() string
	SetRequestID(id string)
}

type introspectedError interface {
	Date() time.Time
	Error() string
	File() string
	Func() string
	Line() int
	Level() int
	Type() string
	RequestId() string
	SetRequestId(id string)
}

type introspectionError struct {
	Dt time.Time `json:"date"`
	Rq string    `json:"requestID,omitempty"`
	Tp string    `json:"type"`
	Mg string    `json:"message"`
	Fl string    `json:"file"`
	Fn string    `json:"function"`
	Ln int       `json:"line"`
}

func (e *introspectionError) Date() time.Time {
	return e.Dt
}
func (e *introspectionError) Error() string {
	return e.Mg
}
func (e *introspectionError) File() string {
	return e.Fl
}
func (e *introspectionError) Func() string {
	return e.Fn
}
func (e *introspectionError) Line() int {
	return e.Ln
}
func (e *introspectionError) Level() int {
	return ErrorLevel
}
func (e *introspectionError) Type() string {
	return ErrorLogType
}
func (e *introspectionError) RequestId() string {
	return e.Rq
}
func (e *introspectionError) SetRequestId(id string) {
	e.Rq = id
}

type infoLevelError struct{ introspectionError }

func (e *infoLevelError) Level() int {
	return InfoLevel
}
func (e *infoLevelError) Type() string {
	return InfoLogType
}

type debugLevelError struct{ introspectionError }

func (e *debugLevelError) Level() int {
	return DebugLevel
}
func (e *debugLevelError) Type() string {
	return DebugLogType
}

type warningLevelError struct{ introspectionError }

func (e *warningLevelError) Level() int {
	return WarningLevel
}
func (e *warningLevelError) Type() string {
	return ErrorLogType
}

type errorLevelError struct{ introspectionError }

func (e *errorLevelError) Level() int {
	return ErrorLevel
}
func (e *errorLevelError) Type() string {
	return ErrorLogType
}

type criticalLevelError struct{ introspectionError }

func (e *criticalLevelError) Level() int {
	return CriticalLevel
}
func (e *criticalLevelError) Type() string {
	return ErrorLogType
}

type emergencyLevelError struct{ introspectionError }

func (e *emergencyLevelError) Level() int {
	return EmergencyLevel
}
func (e *emergencyLevelError) Type() string {
	return ErrorLogType
}
