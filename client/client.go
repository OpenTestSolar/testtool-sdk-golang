package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenTestSolar/testtool-sdk-golang/api"
	"github.com/OpenTestSolar/testtool-sdk-golang/model"
	"github.com/pkg/errors"
)

type ReporterClient struct {
	reportPath string
}

func NewReporterClient(reportPath string) (api.Reporter, error) {
	var dirName string
	if filepath.Ext(reportPath) != "" {
		dirName = filepath.Dir(reportPath)
	} else {
		dirName = reportPath
	}
	if _, err := os.Stat(dirName); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dirName, 0755)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create report path")
			}
		} else {
			return nil, errors.Wrap(err, "failed to stat path")
		}
	}
	return &ReporterClient{
		reportPath: reportPath,
	}, nil
}

func (r *ReporterClient) ReportLoadResult(loadResult *model.LoadResult) error {
	return r.sendJSON(loadResult, "")
}

func (r *ReporterClient) ReportCaseResult(caseResult *model.TestResult) error {
	return r.sendJSON(caseResult, fmt.Sprintf("%s.json", caseResult.TransferNameToHash()))
}

func (r *ReporterClient) ReportJunitXml(filePath string) error {
	return r.reportJunitXml(filePath)
}

func (r *ReporterClient) sendJSON(data interface{}, fileName string) error {
	// Marshal data to JSON with custom datetime encoding
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	// Write JSON data to the file
	if err := r.writeToFile(jsonData, fileName); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	return nil
}

func (r *ReporterClient) writeToFile(data []byte, fileName string) error {
	err := os.WriteFile(filepath.Join(r.reportPath, fileName), data, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}
	return nil
}

func (r *ReporterClient) reportJunitXml(filePath string) error {
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", filePath)
	}
	var testSuite *model.TestSuite = &model.TestSuite{}
	err = xml.Unmarshal(byteValue, testSuite)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s to junit xml", string(byteValue))
	}
	var testResults []*model.TestResult
	for _, testCase := range testSuite.TestCases {
		resultType := model.ResultTypeSucceed
		message := ""
		if testCase.Failure != nil {
			resultType = model.ResultTypeFailed
			message = testCase.Failure.Message
		}

		testResult := &model.TestResult{
			Test: &model.TestCase{
				Name:       fmt.Sprintf("%s.go?%s", strings.ReplaceAll(testCase.ClassName, ".", "/"), testCase.Name),
				Attributes: map[string]string{},
			},
			ResultType: resultType,
			Message:    message,
			Steps:      []*model.TestCaseStep{},
		}

		testResults = append(testResults, testResult)
	}
	for _, testResult := range testResults {
		if err := r.sendJSON(testResult, fmt.Sprintf("%s.json", testResult.TransferNameToHash())); err != nil {
			return errors.Wrapf(err, "failed to send JSON to %s", testResult.TransferNameToHash())
		}
	}
	return nil
}
