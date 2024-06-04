package model

type EntryParam struct {
	// TaskId is the unique identifier for the task.
	TaskId string `json:"TaskId"`
	// ProjectPath is the directory where the user code is located.
	ProjectPath string `json:"ProjectPath"`
	// Context contains the test context information.
	Context map[string]string `json:"Context"`
	// TestSelectors is a list of selectors used to choose test cases.
	TestSelectors []string `json:"TestSelectors"`
	// Collectors is a list of collectors used to gather reports.
	Collectors []string `json:"Collectors"`
	// FileReportPath is the path used for file-based reporting.
	FileReportPath string `json:"FileReportPath"`
}
