// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pointgoal/workstation/pkg/datastore"
	"github.com/pointgoal/workstation/pkg/utils"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-common/error"
	"github.com/rookie-ninja/rk-gin/boot"
	"net/http"
	"strconv"
)

func initApi() {
	var ginEntry *rkgin.GinEntry

	if ginEntry = rkgin.GetGinEntry("workstation"); ginEntry == nil {
		rkcommon.ShutdownWithError(errors.New("nil GinEntry"))
	}

	// Organization
	ginEntry.Router.GET("/v1/org", ListOrg)
	ginEntry.Router.GET("/v1/org/:orgId", GetOrg)
	ginEntry.Router.PUT("/v1/org", CreateOrg)
	ginEntry.Router.DELETE("/v1/org/:orgId", DeleteOrg)
	ginEntry.Router.POST("/v1/org/:orgId", UpdateOrg)

	// Projects
	ginEntry.Router.GET("/v1/org/:orgId/proj", ListProject)
	ginEntry.Router.GET("/v1/org/:orgId/proj/:projId", GetProject)
	ginEntry.Router.PUT("/v1/org/:orgId/proj", CreateProject)
	ginEntry.Router.DELETE("v1/org/:orgId/proj/:projId", DeleteProject)
	ginEntry.Router.POST("v1/org/:orgId/proj/:projId", UpdateProject)
}

// ListOrg
// @Summary List organizations
// @Id 1
// @version 1.0
// @Tags organization
// @produce application/json
// @Success 200 {object} ListOrgResponse
// @Router /v1/org [get]
func ListOrg(ctx *gin.Context) {
	controller := GetController()
	ctx.JSON(http.StatusOK, &ListOrgResponse{
		OrgList: controller.ListOrg(),
	})
}

// GetOrg
// @Summary Get organization
// @Id 2
// @version 1.0
// @Tags organization
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Success 200 {object} GetOrgResponse
// @Router /v1/org/{orgId} [get]
func GetOrg(ctx *gin.Context) {
	controller := GetController()
	orgId := utils.ToInt(ctx.Param("orgId"))

	projIds := make([]string, 0)
	projects := controller.ListProject(orgId)
	for i := range projects {
		projIds = append(projIds, strconv.Itoa(projects[i].Id))
	}

	ctx.JSON(http.StatusOK, &GetOrgResponse{
		Org:        controller.GetOrg(orgId),
		ProjectIds: projIds,
	})
}

// CreateOrg
// @Summary Create organization
// @Id 3
// @version 1.0
// @Tags organization
// @produce application/json
// @Param orgName query string true "Organization name"
// @Success 200 {object} CreateOrgResponse
// @Router /v1/org [put]
func CreateOrg(ctx *gin.Context) {
	controller := GetController()
	orgId := controller.AddOrg(ctx.Query("orgName"))

	ctx.JSON(http.StatusOK, &CreateOrgResponse{
		OrgId: orgId,
	})
}

// DeleteOrg
// @Summary Delete organization
// @Id 4
// @version 1.0
// @Tags organization
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Success 200 {object} DeleteOrgResponse
// @Router /v1/org/{orgId} [delete]
func DeleteOrg(ctx *gin.Context) {
	controller := GetController()
	status := controller.DeleteOrg(utils.ToInt(ctx.Query("orgId")))

	ctx.JSON(http.StatusOK, &DeleteOrgResponse{
		Status: status,
	})
}

// UpdateOrg
// @Summary Update organization
// @Id 5
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org body UpdateOrgRequest true "Organization"
// @Param orgId query int true "Organization Id"
// @Success 200 {object} UpdateOrgResponse
// @Router /v1/org/{orgId} [post]
func UpdateOrg(ctx *gin.Context) {
	controller := GetController()

	req := &UpdateOrgRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	orgId := utils.ToInt(ctx.Query("orgId"))

	var org *datastore.Organization
	if org = controller.GetOrg(orgId); org == nil {
		ctx.JSON(http.StatusOK, &UpdateOrgResponse{
			Status: false,
		})

		return
	}

	org.Name = req.Name

	ctx.JSON(http.StatusOK, &UpdateOrgResponse{
		Status: controller.UpdateOrg(org),
	})
}

// ListProject
// @Summary List projects
// @Id 6
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Success 200 {object} ListProjectResponse
// @Router /v1/org/{orgId}/proj [get]
func ListProject(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Query("orgId"))
	ctx.JSON(http.StatusOK, &ListProjectResponse{
		ProjectList: controller.ListProject(orgId),
	})
}

// GetProject
// @Summary Get project
// @Id 7
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Param projId query int true "Project Id"
// @Success 200 {object} GetProjectResponse
// @Router /v1/org/{orgId}/proj/{projId} [get]
func GetProject(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Query("orgId"))
	projId := utils.ToInt(ctx.Query("projId"))
	ctx.JSON(http.StatusOK, &GetProjectResponse{
		Project: controller.GetProject(orgId, projId),
	})
}

// CreateProject
// @Summary create project
// @Id 8
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Param project body CreateProjectRequest true "Project"
// @Success 200 {object} CreateProjectResponse
// @Router /v1/org/{orgId}/proj [put]
func CreateProject(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Query("orgId"))
	req := &CreateProjectRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	project := datastore.NewProject(orgId, req.Name)
	projId := controller.AddProject(project)

	ctx.JSON(http.StatusOK, &CreateProjectResponse{
		ProjId: projId,
	})
}

// DeleteProject
// @Summary delete project
// @Id 9
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Param projId query int true "Project Id"
// @Success 200 {object} DeleteProjectResponse
// @Router /v1/org/{orgId}/proj/{projId} [delete]
func DeleteProject(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Query("orgId"))
	projId := utils.ToInt(ctx.Query("projId"))

	ctx.JSON(http.StatusOK, &DeleteProjectResponse{
		Status: controller.RemoveProject(orgId, projId),
	})
}

// UpdateProject
// @Summary update project
// @Id 10
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int true "Organization Id"
// @Param projId query int true "Project Id"
// @Param project body UpdateProjectRequest true "Project"
// @Success 200 {object} UpdateProjectResponse
// @Router /v1/org/{orgId}/proj/{projId} [post]
func UpdateProject(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Query("orgId"))
	projId := utils.ToInt(ctx.Query("projId"))

	req := &UpdateProjectRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	proj := controller.GetProject(orgId, projId)
	if proj == nil {
		ctx.JSON(http.StatusOK, &UpdateProjectResponse{
			Status: false,
		})
		return
	}

	proj.Name = req.Name
	ctx.JSON(http.StatusOK, &UpdateProjectResponse{
		Status: controller.UpdateProject(proj),
	})
}
