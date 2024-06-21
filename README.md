# testtool-sdk-golang
TestTool Golang SDK  for TestSolar

## Installation

Use the `go get` command to install this SDK:

```bash
go get github.com/OpenTestSolar/testtool-sdk-golang
```

## Usage Example

Here is a simple usage example:

```go
package main

import (
	"fmt"

	"github.com/OpenTestSolar/testtool-sdk-golang/client"
	"github.com/OpenTestSolar/testtool-sdk-golang/model"
)

func main() {
	// Create a Reporter instance
	reporter, err := client.NewReporterClient("/tmp")
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
}
```

