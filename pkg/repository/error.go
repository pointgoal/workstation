package repository

import "fmt"

const (
	OrgNotFoundMsg     = "organization not found with orgId:%d"
	OrgFailedToGetMsg  = "failed to get organization with orgId:%d"
	ProjNotFoundMsg    = "project not found with orgId:%d projId:%d"
	ProjFailedToGetMsg = "failed to get project with orgId:%d projId:%d"
	ProjFailedToRemove = "failed to remove project with orgId:%d projId:%d"
)

type NotFound struct {
	msg string `yaml:"err" json:"err"`
}

func NewNotFound(msg string) *NotFound {
	return &NotFound{
		msg: msg,
	}
}

func NewNotFoundf(format string, a ...interface{}) *NotFound {
	return NewNotFound(fmt.Sprintf(format, a...))
}

func (e *NotFound) Error() string {
	return e.msg
}
