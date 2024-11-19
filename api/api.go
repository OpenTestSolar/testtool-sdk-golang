package api

import "github.com/OpenTestSolar/testtool-sdk-golang/model"

type Reporter interface {
	ReportLoadResult(loadResult *model.LoadResult) error
	ReportCaseResult(caseResult *model.TestResult) error
	ReportJunitXml(string) error
}
