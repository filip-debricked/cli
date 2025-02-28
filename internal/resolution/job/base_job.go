package job

import (
	"errors"
	"os/exec"
)

type BaseJob struct {
	file   string
	errs   IErrors
	status chan string
}

func NewBaseJob(file string) BaseJob {
	return BaseJob{
		file:   file,
		errs:   NewErrors(file),
		status: make(chan string),
	}
}

func (j *BaseJob) GetFile() string {
	return j.file
}

func (j *BaseJob) Errors() IErrors {
	return j.errs
}

func (j *BaseJob) ReceiveStatus() chan string {
	return j.status
}

func (j *BaseJob) SendStatus(status string) {
	j.status <- status
}

func (j *BaseJob) GetExitError(err error) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	return errors.New(string(exitErr.Stderr))
}
