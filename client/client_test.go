package client

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/OpenTestSolar/testtool-sdk-golang/model"
	"github.com/google/go-cmp/cmp"
)

// TestReporter_ReportLoadResult tests the ReportLoadResult method
func TestReporter_ReportLoadResult(t *testing.T) {
	// Prepare test data
	loadResult := &model.LoadResult{
		Tests: []model.TestCase{
			{
				Name:       "test1",
				Attributes: map[string]string{},
			},
		},
		LoadErrors: []model.LoadError{
			{
				Name:    "test2",
				Message: "load failed",
			},
		},
	}

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
	verifyPipeData[model.LoadResult](t, tmpFile, *loadResult)
}

// TestReporter_ReportCaseResult tests the ReportCaseResult method
func TestReporter_ReportCaseResult(t *testing.T) {
	// Prepare test data
	caseResult := &model.TestResult{
		Test: model.TestCase{
			Name:       "test1",
			Attributes: map[string]string{},
		},
		StartTime:  time.Now(),
		ResultType: model.ResultTypeSucceed,
		Message:    "test passed",
		EndTime:    nil,
		Steps:      nil,
	}

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
	verifyPipeData[model.TestResult](t, tmpFile, *caseResult)
}

// verifyPipeData reads data from the pipe file and verifies if it meets expectations
func verifyPipeData[T any](t *testing.T, pipeFile *os.File, expected T) {
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

	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal data to obj: %v", err)
	}
	if !cmp.Equal(expected, result) {
		t.Errorf("Expected result %+v, but got %+v", expected, result)
	}
}
