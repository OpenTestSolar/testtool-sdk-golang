package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/OpenTestSolar/testtool-sdk-golang/model"
)

func generateLoadResults() *model.LoadResult {
	return &model.LoadResult{
		Tests: []*model.TestCase{
			{
				Name: "test1",
				Attributes: map[string]string{
					"key": "value",
				},
			},
		},
		LoadErrors: []*model.LoadError{
			{
				Name:    "test2",
				Message: "load failed",
			},
		},
	}
}

// TestReporter_ReportLoadResult tests the ReportLoadResult method
func TestReporter_ReportLoadResult(t *testing.T) {
	// Prepare test data
	loadResult := generateLoadResults()

	// Create a temporary file to simulate the pipe
	tmpFile, err := os.CreateTemp("", "pipe")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a ReporterClient instance
	reporter, err := NewReporterClient()
	if err != nil {
		t.Fatalf("Failed to create reporter: %v", err)
	}
	defer reporter.Close()

	// Set reporter's pipeIO to the temporary file
	reporter.(*ReporterClient).pipeIO = tmpFile

	// Call ReportLoadResult method
	err = reporter.ReportLoadResult(loadResult)
	if err != nil {
		t.Fatalf("ReportLoadResult failed: %v", err)
	}

	// Verify result
	verifyPipeData(t, tmpFile)
}

func generateCaseResult() *model.TestResult {
	return &model.TestResult{
		Test: &model.TestCase{
			Name:       "test1",
			Attributes: map[string]string{},
		},
		StartTime:  time.Now(),
		ResultType: model.ResultTypeSucceed,
		Message:    "test passed",
		EndTime:    time.Now(),
		Steps: []*model.TestCaseStep{
			{
				StartTime:  time.Now(),
				Title:      "step1",
				ResultType: model.ResultTypeSucceed,
				EndTime:    time.Now(),
				Logs: []*model.TestCaseLog{
					{
						Time:    time.Now(),
						Level:   model.LogLevelInfo,
						Content: "step1 passed",
						AssertError: &model.TestCaseAssertError{
							Expect:  "expect",
							Actual:  "actual",
							Message: "assert failed",
						},
						RuntimeError: &model.TestCaseRuntimeError{
							Summary: "runtime error",
							Detail:  "runtime error detail",
						},
						Attachments: []*model.Attachment{
							{
								Name:           "attachment1",
								Url:            "http://example.com/attachment1",
								AttachmentType: model.AttachmentTypeFile,
							},
						},
					},
				},
			},
		},
	}
}

// TestReporter_ReportCaseResult tests the ReportCaseResult method
func TestReporter_ReportCaseResult(t *testing.T) {
	// Prepare test data
	caseResult := generateCaseResult()

	// Create a temporary file to simulate the pipe
	tmpFile, err := os.CreateTemp("", "pipe")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a Reporter instance
	reporter, err := NewReporterClient()
	if err != nil {
		t.Fatalf("Failed to create reporter: %v", err)
	}
	defer reporter.Close()

	// Set reporter's pipeIO to the temporary file
	reporter.(*ReporterClient).pipeIO = tmpFile

	// Call ReportCaseResult method
	err = reporter.ReportCaseResult(caseResult)
	if err != nil {
		t.Fatalf("ReportCaseResult failed: %v", err)
	}

	// Verify result
	verifyPipeData(t, tmpFile)
}

func verifyCaseResultFields(t *testing.T, result map[string]interface{}) {
	verifyDateTimeFormat := func(obj map[string]interface{}, key string) {
		timeStr, ok := obj[key].(string)
		if !ok {
			t.Fatal("time field is missing or not a string")
		}
		if _, err := time.Parse(model.DateTimeFormat, timeStr); err != nil {
			t.Fatalf("time field does not match the expected format: %v", err)
		}
	}
	for _, field := range []string{
		"StartTime",
		"EndTime",
	} {
		verifyDateTimeFormat(result, field)
	}
	steps := result["Steps"].([]interface{})
	for _, step := range steps {
		step := step.(map[string]interface{})
		for _, field := range []string{
			"StartTime",
			"EndTime",
		} {
			verifyDateTimeFormat(step, field)
		}
		if _, ok := step["Logs"]; ok {
			testCaseLogs := step["Logs"].([]interface{})
			for _, testCaseLog := range testCaseLogs {
				testCaseLog := testCaseLog.(map[string]interface{})
				level := testCaseLog["Level"].(string)
				if level != "INFO" {
					t.Errorf("Incorrect log level: %s", level)
				}
				content := testCaseLog["Content"].(string)
				if content != "step1 passed" {
					t.Errorf("Incorrect log content: %s", content)
				}
				verifyDateTimeFormat(testCaseLog, "Time")
			}
		}
	}
	resultType := result["ResultType"].(string)
	if resultType != "SUCCEED" {
		t.Errorf("Incorrect result type: %s", resultType)
	}
	message := result["Message"].(string)
	if message != "test passed" {
		t.Errorf("Incorrect message: %s", message)
	}
}

func verifyLoadResultFields(t *testing.T, result map[string]interface{}) {
	tests := result["Tests"].([]interface{})
	for _, test := range tests {
		test := test.(map[string]interface{})
		name := test["Name"].(string)
		if name != "test1" {
			t.Error("Incorrect test name")
		}
		attr := test["Attributes"].(map[string]interface{})
		for key, value := range attr {
			value = value.(string)
			if key == "key" && value == "value" {
				continue
			}
			t.Error("Incorrect test attribute")
		}
	}
	loadErrors := result["LoadErrors"].([]interface{})
	for _, loadError := range loadErrors {
		loadError := loadError.(map[string]interface{})
		name := loadError["Name"].(string)
		if name != "test2" {
			t.Error("Incorrect test name")
		}
		message := loadError["Message"].(string)
		if message != "load failed" {
			t.Error("Incorrect message")
		}
	}
}

// verifyPipeData reads data from the pipe file and verifies if it meets expectations
func verifyPipeData(t *testing.T, pipeFile *os.File) {
	// Re-open the temporary file to read from the beginning
	file, err := os.Open(pipeFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	// Read and verify magic number
	var magicNumber uint32
	err = binary.Read(file, binary.LittleEndian, &magicNumber)
	if err != nil {
		t.Fatalf("Failed to read magic number: %v", err)
	}
	if magicNumber != MagicNumber {
		t.Errorf("Expected magic number %v, but got %v", MagicNumber, magicNumber)
	}

	// Read and verify data length
	var length uint32
	err = binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		t.Fatalf("Failed to read length: %v", err)
	}

	// Read and verify data
	data := make([]byte, length)
	_, err = file.Read(data)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}
	var result map[string]interface{}
	fmt.Printf("Raw json string: %s\n", string(data))
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal data to obj: %v", err)
	}
	if _, ok := result["ResultType"]; ok {
		verifyCaseResultFields(t, result)
	} else {
		verifyLoadResultFields(t, result)
	}
}
