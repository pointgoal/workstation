// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-common/error"
	"github.com/rookie-ninja/rk-gin/boot"
	"net/http"
)

func initApi() {
	var ginEntry *rkgin.GinEntry

	if ginEntry = rkgin.GetGinEntry("workstation"); ginEntry == nil {
		rkcommon.ShutdownWithError(errors.New("nil GinEntry"))
	}

	// Organization
	ginEntry.Router.GET("/v1/org", ListOrgs)
	ginEntry.Router.GET("/v1/org/:org_id", GetOrg)
	ginEntry.Router.PUT("/v1/org", CreateOrg)
	ginEntry.Router.DELETE("/v1/org/:org_id", DeleteOrg)
	ginEntry.Router.POST("/v1/org/:org_id", UpdateOrg)

	// Projects
	ginEntry.Router.GET("/v1/org/:org_id/proj", ListProjects)
	ginEntry.Router.GET("/v1/org/:org_id/proj/:proj_id", GetProject)
	ginEntry.Router.PUT("/v1/org/:org_id/proj", CreateProject)
	ginEntry.Router.DELETE("v1/org/:org_id/proj/:proj_id", DeleteProject)
	ginEntry.Router.POST("v1/org/:org_id/proj/:proj_id", UpdateProject)
}

// ListOrgsResponse response of list organization
type ListOrgsResponse struct {
	Orgs []*Organization `yaml:"orgs" json:"orgs"`
}

// ListOrgs
// @Summary List organizations
// @Id 1
// @version 1.0
// @Tags organization
// @produce application/json
// @Success 200 {object} ListOrgsResponse
// @Router /v1/org [get]
func ListOrgs(ctx *gin.Context) {
	entry := GetEntry()
	ctx.JSON(http.StatusOK, &ListOrgsResponse{
		Orgs: entry.ListOrgs(),
	})
}

// GetOrgResponse response of get organization
type GetOrgResponse struct {
	Org *Organization `yaml:"orgs" json:"orgs"`
}

// GetOrg
// @Summary Get organization
// @Id 2
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org_id path string true "Organization Id"
// @Success 200 {object} GetOrgResponse
// @Router /v1/org/{org_id} [get]
func GetOrg(ctx *gin.Context) {
	entry := GetEntry()
	ctx.JSON(http.StatusOK, &GetOrgResponse{
		Org: entry.GetOrg(ctx.Param("org_id")),
	})
}

// CreateOrgResponse response of create organization
type CreateOrgResponse struct {
	OrgId string `yaml:"orgId" json:"orgId"`
}

// CreateOrg
// @Summary Create organization
// @Id 3
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org_name query string true "Organization name"
// @Success 200 {object} CreateOrgResponse
// @Router /v1/org [put]
func CreateOrg(ctx *gin.Context) {
	entry := GetEntry()
	orgId := entry.AddOrg(ctx.Query("org_name"))

	ctx.JSON(http.StatusOK, &CreateOrgResponse{
		OrgId: orgId,
	})
}

// DeleteOrgResponse response of delete organization
type DeleteOrgResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// DeleteOrg
// @Summary Delete organization
// @Id 4
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Success 200 {object} DeleteOrgResponse
// @Router /v1/org/{org_id} [delete]
func DeleteOrg(ctx *gin.Context) {
	entry := GetEntry()
	status := entry.DeleteOrg(ctx.Query("org_id"))

	ctx.JSON(http.StatusOK, &DeleteOrgResponse{
		Status: status,
	})
}

// DeleteOrgResponse response of update organization
type UpdateOrgResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateOrgRequest request body of update organization
type UpdateOrgRequest struct {
	Name string `yaml:"name" json:"name"`
}

// UpdateOrg
// @Summary Update organization
// @Id 5
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org body UpdateOrgRequest true "Organization"
// @Param org_id query string true "Organization Id"
// @Success 200 {object} UpdateOrgResponse
// @Router /v1/org/{org_id} [post]
func UpdateOrg(ctx *gin.Context) {
	entry := GetEntry()

	req := &UpdateOrgRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	orgId := ctx.Query("org_id")

	var org *Organization
	if org = entry.GetOrg(orgId); org == nil {
		ctx.JSON(http.StatusOK, &UpdateOrgResponse{
			Status: false,
		})

		return
	}

	org.Name = req.Name

	ctx.JSON(http.StatusOK, &UpdateOrgResponse{
		Status: entry.UpdateOrg(org),
	})
}

// ListProjectResponse response of list projects
type ListProjectResponse struct {
	OrgId    string     `yaml:"orgId" json:"orgId"`
	Projects []*Project `yaml:"projects" json:"projects"`
}

// ListProjects
// @Summary List projects
// @Id 6
// @version 1.0
// @Tags project
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Success 200 {object} ListProjectResponse
// @Router /v1/org/{org_id}/proj [get]
func ListProjects(ctx *gin.Context) {
	entry := GetEntry()

	orgId := ctx.Query("org_id")
	ctx.JSON(http.StatusOK, &ListProjectResponse{
		OrgId:    orgId,
		Projects: entry.ListProjects(orgId),
	})
}

// GetProjectResponse response of get project
type GetProjectResponse struct {
	OrgId   string   `yaml:"orgId" json:"orgId"`
	Project *Project `yaml:"project" json:"project"`
}

// GetProject
// @Summary Get project
// @Id 7
// @version 1.0
// @Tags project
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Param proj_id query string true "Project Id"
// @Success 200 {object} GetProjectResponse
// @Router /v1/org/{org_id}/proj/{proj_id} [get]
func GetProject(ctx *gin.Context) {
	entry := GetEntry()

	orgId := ctx.Query("org_id")
	projId := ctx.Query("proj_id")
	ctx.JSON(http.StatusOK, &GetProjectResponse{
		OrgId:   orgId,
		Project: entry.GetProject(orgId, projId),
	})
}

// CreateProjectResponse response of create project
type CreateProjectResponse struct {
	OrgId  string `yaml:"orgId" json:"orgId"`
	ProjId string `yaml:"projId" json:"projId"`
}

// CreateProjectRequest request body
type CreateProjectRequest struct {
	Name string `yaml:"name" json:"name"`
}

// CreateProject
// @Summary create project
// @Id 8
// @version 1.0
// @Tags project
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Param project body CreateProjectRequest true "Project"
// @Success 200 {object} CreateProjectResponse
// @Router /v1/org/{org_id}/proj [put]
func CreateProject(ctx *gin.Context) {
	entry := GetEntry()

	orgId := ctx.Query("org_id")
	req := &CreateProjectRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	project := NewProject(req.Name)
	projId := entry.AddProject(orgId, project)

	ctx.JSON(http.StatusOK, &CreateProjectResponse{
		OrgId:  orgId,
		ProjId: projId,
	})
}

// DeleteProjectResponse response of delete project
type DeleteProjectResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// DeleteProject
// @Summary delete project
// @Id 9
// @version 1.0
// @Tags project
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Param proj_id query string true "Project Id"
// @Success 200 {object} DeleteProjectResponse
// @Router /v1/org/{org_id}/proj/{proj_id} [delete]
func DeleteProject(ctx *gin.Context) {
	entry := GetEntry()

	orgId := ctx.Query("org_id")
	projId := ctx.Query("proj_id")

	ctx.JSON(http.StatusOK, &DeleteProjectResponse{
		Status: entry.DeleteProject(orgId, projId),
	})
}

// UpdateProjectResponse response of update project
type UpdateProjectResponse struct {
	Status bool `yaml:"status" json:"status"`
}

// UpdateProjectRequest request body
type UpdateProjectRequest struct {
	Name string `yaml:"name" json:"name"`
}

// UpdateProject
// @Summary update project
// @Id 10
// @version 1.0
// @Tags project
// @produce application/json
// @Param org_id query string true "Organization Id"
// @Param proj_id query string true "Project Id"
// @Param project body UpdateProjectRequest true "Project"
// @Success 200 {object} UpdateProjectResponse
// @Router /v1/org/{org_id}/proj/{proj_id} [post]
func UpdateProject(ctx *gin.Context) {
	entry := GetEntry()

	orgId := ctx.Query("org_id")
	projId := ctx.Query("proj_id")

	req := &UpdateProjectRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	proj := entry.GetProject(orgId, projId)
	if proj == nil {
		ctx.JSON(http.StatusOK, &UpdateProjectResponse{
			Status: false,
		})
		return
	}

	proj.Name = req.Name
	ctx.JSON(http.StatusOK, &UpdateProjectResponse{
		Status: entry.UpdateProject(orgId, proj),
	})
}
