package controller

import "github.com/pointgoal/workstation/pkg/datastore"

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// ListOrgResponse response of list organization
type ListOrgResponse struct {
	OrgList []*datastore.Organization `yaml:"orgList" json:"orgList"`
}

// GetOrgResponse response of get organization
type GetOrgResponse struct {
	Org        *datastore.Organization `yaml:"org" json:"org"`
	ProjectIds []string                `yaml:"projectIds" json:"projectIds"`
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

// ListProjectResponse response of list projects
type ListProjectResponse struct {
	ProjectList []*datastore.Project `yaml:"projectList" json:"projectList"`
}

// GetProjectResponse response of get project
type GetProjectResponse struct {
	Project *datastore.Project `yaml:"project" json:"project"`
}

// CreateProjectResponse response of create project
type CreateProjectResponse struct {
	ProjId int `yaml:"projId" json:"projId"`
}

// CreateProjectRequest request body
type CreateProjectRequest struct {
	Name string `yaml:"name" json:"name"`
}

// DeleteProjectResponse response of delete project
type DeleteProjectResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateProjectResponse response of update project
type UpdateProjectResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateProjectRequest request body
type UpdateProjectRequest struct {
	Name string `yaml:"name" json:"name"`
}
