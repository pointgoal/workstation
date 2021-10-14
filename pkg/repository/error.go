package repository

import "fmt"

const (
	OrgNotFoundMsg             = "organization not found with orgId:%d"
	OrgFailedToGetMsg          = "failed to get organization with orgId:%d"
	ProjNotFoundMsg            = "project not found with projId:%d"
	ProjFailedToGetMsg         = "failed to get project with projId:%d"
	ProjFailedToRemove         = "failed to remove project with projId:%d"
	SourceNotFoundMsg          = "source not found with sourceId:%d"
	SourceFailedToGetMsg       = "failed to get source with sourceId:%d"
	SourceFailedToRemove       = "failed to remove source with sourceId:%d"
	OauthSourceNotFoundMsg     = "oauth source not found with source:%s"
	AccessTokenAlreadyExistMsg = "access token already exist with type:%s user:%s"
	AccessTokenNotFoundMsg     = "access token not found with type:%s user:%s"
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

type AlreadyExist struct {
	msg string `yaml:"err" json:"err"`
}

func NewAlreadyExist(msg string) *NotFound {
	return &NotFound{
		msg: msg,
	}
}

func NewAlreadyExistf(format string, a ...interface{}) *NotFound {
	return NewAlreadyExist(fmt.Sprintf(format, a...))
}

func (e *AlreadyExist) Error() string {
	return e.msg
}
