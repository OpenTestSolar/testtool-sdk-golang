package client

import (
	"encoding/binary"
	"encoding/json"
	"os"

	"github.com/OpenTestSolar/testtool-sdk-golang/api"
	"github.com/OpenTestSolar/testtool-sdk-golang/model"
	"github.com/nightlyone/lockfile"
	"github.com/pkg/errors"
)

const (
	MagicNumber uint32 = 0x1234ABCD
	PipeWriter  int    = 3
)

type ReporterClient struct {
	pipeIO       *os.File
	lockFile     lockfile.Lockfile
	lockFilePath string
}

func NewReporterClient() (api.Reporter, error) {
	pipeIO := os.NewFile(uintptr(PipeWriter), "pipe")
	lockFilePath := "/tmp/testsolar_reporter.lock"

	lock, err := lockfile.New(lockFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lock file")
	}

	return &ReporterClient{
		pipeIO:       pipeIO,
		lockFile:     lock,
		lockFilePath: lockFilePath,
	}, nil
}

func (r *ReporterClient) ReportLoadResult(loadResult *model.LoadResult) error {
	return r.sendJSON(loadResult)
}

func (r *ReporterClient) ReportCaseResult(caseResult *model.TestResult) error {
	return r.sendJSON(caseResult)
}

func (r *ReporterClient) Close() error {
	if r.pipeIO != nil {
		return r.pipeIO.Close()
	}
	return nil
}

func (r *ReporterClient) sendJSON(data interface{}) error {
	// Acquire lock
	err := r.lockFile.TryLock()
	if err != nil {
		return errors.Wrap(err, "failed to acquire lock")
	}
	defer r.lockFile.Unlock()

	// Marshal data to JSON with custom datetime encoding
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	// Write magic number, length, and JSON data to the pipe
	if err := r.writeToPipe(MagicNumber, jsonData); err != nil {
		return errors.Wrap(err, "failed to write to pipe")
	}

	return nil
}

func (r *ReporterClient) writeToPipe(magicNumber uint32, data []byte) error {
	length := uint32(len(data))

	// Write magic number
	if err := binary.Write(r.pipeIO, binary.LittleEndian, magicNumber); err != nil {
		return err
	}

	// Write length
	if err := binary.Write(r.pipeIO, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write data
	if _, err := r.pipeIO.Write(data); err != nil {
		return err
	}

	// Flush the pipe if possible
	if err := r.pipeIO.Sync(); err != nil {
		return err
	}

	return nil
}
