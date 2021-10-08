package controller

import "github.com/pointgoal/workstation/pkg/repository"

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
	Name string `yaml:"name" json:"name"`
}

// DeleteProjectResponse response of delete project
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