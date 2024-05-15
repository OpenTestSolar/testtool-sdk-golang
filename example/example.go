package main

import (
	"fmt"

	"github.com/OpenTestSolar/testtool-sdk-golang/client"
	"github.com/OpenTestSolar/testtool-sdk-golang/model"
)

func main() {
	// Create a Reporter instance
	reporter, err := client.NewReporterClient()
	if err != nil {
		fmt.Printf("Failed to create reporter: %v\n", err)
		return
	}

	// Create a LoadResult object
	loadResult := &model.LoadResult{
		// ... Initialize LoadResult struct ...
	}

	// Use Reporter to report LoadResult
	err = reporter.ReportLoadResult(loadResult)
	if err != nil {
		fmt.Printf("Failed to report load result: %v\n", err)
	}

	// Create a TestResult object
	caseResult := &model.TestResult{
		// ... Initialize TestResult struct ...
	}

	// Use Reporter to report TestResult
	err = reporter.ReportCaseResult(caseResult)
	if err != nil {
		fmt.Printf("Failed to report case result: %v\n", err)
	}

	// Close Reporter
	err = reporter.Close()
	if err != nil {
		fmt.Printf("Failed to close reporter: %v\n", err)
	}
}
