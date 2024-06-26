package model

import (
	"encoding/json"
	"time"
)

// DateTimeFormat is the format string for JSON datetime representation.
const DateTimeFormat = "2006-01-02T15:04:05.000Z"

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
	Expect  string `json:"Expect"`
	Actual  string `json:"Actual"`
	Message string `json:"Message"`
}

// TestCaseRuntimeError represents a runtime error in a test case.
type TestCaseRuntimeError struct {
	Summary string `json:"Summary"`
	Detail  string `json:"Detail"`
}

// Attachment represents an attachment in a test case log.
type Attachment struct {
	Name           string         `json:"Name"`
	Url            string         `json:"Url"`
	AttachmentType AttachmentType `json:"AttachmentType"`
}

// TestCaseLog represents a log entry for a test case.
type TestCaseLog struct {
	Time         time.Time             `json:"-"`
	Level        LogLevel              `json:"Level"`
	Content      string                `json:"Content"`
	AssertError  *TestCaseAssertError  `json:"AssertError"`
	RuntimeError *TestCaseRuntimeError `json:"RuntimeError"`
	Attachments  []*Attachment         `json:"Attachments"`
}

// MarshalJSON implements the json.Marshaler interface for TestCaseLog.
func (tcl TestCaseLog) MarshalJSON() ([]byte, error) {
	type Alias TestCaseLog
	return json.Marshal(&struct {
		Time string `json:"Time"`
		*Alias
	}{
		Time:  tcl.Time.UTC().Format(DateTimeFormat),
		Alias: (*Alias)(&tcl),
	})
}

// TestCaseStep represents a step in a test case.
type TestCaseStep struct {
	StartTime  time.Time      `json:"-"`
	Title      string         `json:"Title"`
	ResultType ResultType     `json:"ResultType"`
	EndTime    time.Time      `json:"-"`
	Logs       []*TestCaseLog `json:"Logs"`
}

// MarshalJSON implements the json.Marshaler interface for TestCaseStep.
func (tcs TestCaseStep) MarshalJSON() ([]byte, error) {
	type Alias TestCaseStep
	return json.Marshal(&struct {
		StartTime string `json:"StartTime"`
		EndTime   string `json:"EndTime"`
		*Alias
	}{
		StartTime: tcs.StartTime.UTC().Format(DateTimeFormat),
		EndTime:   tcs.EndTime.UTC().Format(DateTimeFormat),
		Alias:     (*Alias)(&tcs),
	})
}

// TestResult represents the result of a test case.
type TestResult struct {
	Test       *TestCase       `json:"Test"`
	StartTime  time.Time       `json:"-"`
	ResultType ResultType      `json:"ResultType"`
	Message    string          `json:"Message"`
	EndTime    time.Time       `json:"-"`
	Steps      []*TestCaseStep `json:"Steps"`
}

// MarshalJSON implements the json.Marshaler interface for TestResult.
func (tr TestResult) MarshalJSON() ([]byte, error) {
	type Alias TestResult
	return json.Marshal(&struct {
		StartTime string `json:"StartTime"`
		EndTime   string `json:"EndTime"`
		*Alias
	}{
		StartTime: tr.StartTime.UTC().Format(DateTimeFormat),
		EndTime:   tr.EndTime.UTC().Format(DateTimeFormat),
		Alias:     (*Alias)(&tr),
	})
}
