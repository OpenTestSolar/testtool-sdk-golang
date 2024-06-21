package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	// Create a ReporterClient instance
	reporter, err := NewReporterClient("./tmp/result.json")
	if err != nil {
		t.Fatalf("Failed to create reporter: %v", err)
	}
	defer os.RemoveAll("./tmp")

	// Call ReportLoadResult method
	err = reporter.ReportLoadResult(loadResult)
	if err != nil {
		t.Fatalf("ReportLoadResult failed: %v", err)
	}

	// Verify result
	verifyPipeData(t, "./tmp")
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

	// Create a Reporter instance
	reporter, err := NewReporterClient("./tmp")
	if err != nil {
		t.Fatalf("Failed to create reporter: %v", err)
	}
	defer os.RemoveAll("./tmp")
	// Call ReportCaseResult method
	err = reporter.ReportCaseResult(caseResult)
	if err != nil {
		t.Fatalf("ReportCaseResult failed: %v", err)
	}

	// Verify result
	verifyPipeData(t, "./tmp")
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
				level := int32(testCaseLog["Level"].(float64))
				if level != 2 {
					t.Errorf("Incorrect log level: %d", level)
				}
				content := testCaseLog["Content"].(string)
				if content != "step1 passed" {
					t.Errorf("Incorrect log content: %s", content)
				}
				verifyDateTimeFormat(testCaseLog, "Time")
			}
		}
	}
	resultType := int32(result["ResultType"].(float64))
	if resultType != 1 {
		t.Errorf("Incorrect result type: %d", resultType)
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
func verifyPipeData(t *testing.T, tmpDir string) {
	// Re-open the temporary file to read from the beginning
	err := filepath.WalkDir(tmpDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var result map[string]interface{}
			fmt.Printf("Raw json string: %s\n", string(content))
			err = json.Unmarshal(content, &result)
			if err != nil {
				return err
			}
			if _, ok := result["ResultType"]; ok {
				verifyCaseResultFields(t, result)
			} else {
				verifyLoadResultFields(t, result)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk temp dir: %v", err)
	}
}
