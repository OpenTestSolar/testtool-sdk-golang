package model

import (
	"time"
)

// ResultType is an enumeration of possible test result types.
type ResultType string

// Enum values for ResultType
const (
	ResultTypeUnknown    ResultType = "UNKNOWN"
	ResultTypeSucceed    ResultType = "SUCCEED"
	ResultTypeFailed     ResultType = "FAILED"
	ResultTypeLoadFailed ResultType = "LOAD_FAILED"
	ResultTypeIgnored    ResultType = "IGNORED"
	ResultTypeRunning    ResultType = "RUNNING"
	ResultTypeWaiting    ResultType = "WAITING"
)

// LogLevel is an enumeration of possible log levels.
type LogLevel string

// Enum values for LogLevel
const (
	LogLevelTrace LogLevel = "VERBOSE"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARNNING"
	LogLevelError LogLevel = "ERROR"
)

// AttachmentType is an enumeration of possible attachment types.
type AttachmentType string

// Enum values for AttachmentType
const (
	AttachmentTypeFile   AttachmentType = "FILE"
	AttachmentTypeURL    AttachmentType = "URL"
	AttachmentTypeIFrame AttachmentType = "IFRAME"
)

// TestCaseAssertError represents an assertion error in a test case.
type TestCaseAssertError struct {
	Expect  string
	Actual  string
	Message string
}

// TestCaseRuntimeError represents a runtime error in a test case.
type TestCaseRuntimeError struct {
	Summary string
	Detail  string
}

// Attachment represents an attachment in a test case log.
type Attachment struct {
	Name           string
	Url            string
	AttachmentType AttachmentType
}

// TestCaseLog represents a log entry for a test case.
type TestCaseLog struct {
	Time         time.Time
	Level        LogLevel
	Content      string
	AssertError  *TestCaseAssertError
	RuntimeError *TestCaseRuntimeError
	Attachments  []Attachment
}

// TestCaseStep represents a step in a test case.
type TestCaseStep struct {
	StartTime  time.Time
	Title      string
	ResultType ResultType
	EndTime    *time.Time
	Logs       []TestCaseLog
}

// TestResult represents the result of a test case.
type TestResult struct {
	Test       TestCase
	StartTime  time.Time
	ResultType ResultType
	Message    string
	EndTime    *time.Time
	Steps      []TestCaseStep
}
