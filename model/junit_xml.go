package model

import "encoding/xml"

type TestSuite struct {
	XMLName   xml.Name         `xml:"testsuite"`
	TestCases []*JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	ClassName string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure"`
}

type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}
