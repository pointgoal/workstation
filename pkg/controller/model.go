package controller

import (
	"github.com/pointgoal/workstation/pkg/repository"
	"time"
)

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// Org is model for API response
type Org struct {
	Meta    *repository.Org `yaml:"meta" json:"meta"`
	ProjIds []int           `yaml:"projIds" json:"projIds"`
}

// ListOrgResponse response of list organization
type ListOrgResponse struct {
	OrgList []*Org `yaml:"orgList" json:"orgList"`
}

// GetOrgResponse response of get organization
type GetOrgResponse struct {
	Org *Org `yaml:"org" json:"org"`
}

// CreateOrgResponse response of create organization
type CreateOrgResponse struct {
	OrgId int `yaml:"orgId" json:"orgId"`
}

// DeleteOrgResponse response of delete organization
type DeleteOrgResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateOrgResponse response of update organization
type UpdateOrgResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateOrgRequest request body of update organization
type UpdateOrgRequest struct {
	Name string `yaml:"name" json:"name"`
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// Proj is model for API response
type Proj struct {
	Meta *repository.Proj `yaml:"meta" json:"meta"`
}

// ListProjResponse response of list projects
type ListProjResponse struct {
	ProjList []*Proj `yaml:"projList" json:"projList"`
}

// GetProjResponse response of get project
type GetProjResponse struct {
	Proj *Proj `yaml:"proj" json:"proj"`
}

// CreateProjResponse response of create project
type CreateProjResponse struct {
	OrgId  int `yaml:"orgId" json:"orgId"`
	ProjId int `yaml:"projId" json:"projId"`
}

// CreateProjRequest request body
type CreateProjRequest struct {
	OrgId int    `yaml:"orgId" json:"orgId"`
	Name  string `yaml:"name" json:"name"`
}

// DeleteProjResponse response of delete project
type DeleteProjResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateProjResponse response of update project
type UpdateProjResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateProjRequest request body
type UpdateProjRequest struct {
	Name string `yaml:"name" json:"name"`
}

// ******************************************** //
// ************** Source related ************** //
// ******************************************** //

// CreateSourceRequest request body
type CreateSourceRequest struct {
	Type       string `yaml:"type" json:"type"`
	Repository string `yaml:"repository" json:"repository"`
}

// CreateSourceResponse response of create source
type CreateSourceResponse struct {
	ProjId   int `yaml:"projId" json:"projId"`
	SourceId int `yaml:"sourceId" json:"sourceId"`
}

// DeleteSourceResponse response of delete source
type DeleteSourceResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// ************************************************** //
// ************ PipelineTemplate related ************ //
// ************************************************** //

// PipelineTemplate is model for API response
type PipelineTemplate struct {
	Meta *repository.PipelineTemplate `yaml:"meta" json:"meta"`
}

// ListPipelineTemplateResponse response of list organization
type ListPipelineTemplateResponse struct {
	TemplateList []*PipelineTemplate `yaml:"templateList" json:"templateList"`
}

// ListCommitsResponse response of user commits of source
type ListCommitsResponse struct {
	Commits []*Commit `yaml:"commits" json:"commits"`
}

type Commit struct {
	Id           string    `yaml:"id" json:"id"`
	Url          string    `yaml:"url" json:"url"`
	Message      string    `yaml:"message" json:"message"`
	Date         time.Time `yaml:"date" json:"date"`
	Committer    string    `yaml:"committer" json:"committer"`
	CommitterUrl string    `yaml:"committerUrl" json:"committerUrl"`
	Artifact     *Artifact `yaml:"artifact" json:"artifact"`
}

type Artifact struct {
	Id   int    `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}
