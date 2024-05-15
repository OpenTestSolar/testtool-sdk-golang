package model

// TestCase represents a test case with a name and a set of attributes.
type TestCase struct {
	Name       string            `json:"Name"`
	Attributes map[string]string `json:"Attributes"`
}
