package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
