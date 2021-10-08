// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pointgoal/workstation/pkg/repository"
	"github.com/pointgoal/workstation/pkg/utils"
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
	ginEntry.Router.GET("/v1/org", ListOrg)
	ginEntry.Router.GET("/v1/org/:orgId", GetOrg)
	ginEntry.Router.PUT("/v1/org", CreateOrg)
	ginEntry.Router.DELETE("/v1/org/:orgId", DeleteOrg)
	ginEntry.Router.POST("/v1/org/:orgId", UpdateOrg)

	// Project
	ginEntry.Router.GET("/v1/org/:orgId/proj", ListProj)
	ginEntry.Router.GET("/v1/org/:orgId/proj/:projId", GetProj)
	ginEntry.Router.PUT("/v1/org/:orgId/proj", CreateProj)
	ginEntry.Router.DELETE("v1/org/:orgId/proj/:projId", DeleteProj)
	ginEntry.Router.POST("v1/org/:orgId/proj/:projId", UpdateProj)
}

func makeInternalError(ctx *gin.Context, message string, details ...interface{}) {
	ctx.JSON(http.StatusInternalServerError, rkerror.New(
		rkerror.WithHttpCode(http.StatusInternalServerError),
		rkerror.WithMessage(message),
		rkerror.WithDetails(details...)))
}

func makeNotFoundError(ctx *gin.Context, message string, details ...interface{}) {
	ctx.JSON(http.StatusNotFound, rkerror.New(
		rkerror.WithHttpCode(http.StatusNotFound),
		rkerror.WithMessage(message),
		rkerror.WithDetails(details...)))
}

func convertOrg(orgFromRepo *repository.Org, projFromRepo []*repository.Proj) *Org {
	org := &Org{
		Meta:    orgFromRepo,
		ProjIds: make([]int, 0),
	}

	for i := range projFromRepo {
		org.ProjIds = append(org.ProjIds, projFromRepo[i].Id)
	}

	return org
}

func convertProj(projFromRepo *repository.Proj) *Proj {
	proj := &Proj{
		Meta: projFromRepo,
	}

	return proj
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
	orgList := make([]*Org, 0)

	// 1: list organization
	orgListFromRepo, err := controller.Repo.ListOrg()
	if err != nil {
		makeInternalError(ctx, "failed to list organizations", err)
		return
	}

	// 2: list projects
	for i := range orgListFromRepo {
		projListFromRepo, err := controller.Repo.ListProj(orgListFromRepo[i].Id)
		if err != nil {
			makeInternalError(ctx, fmt.Sprintf("failed to list projects with orgId:%d.", orgListFromRepo[i].Id), err)
			return
		}
		// convert to API model
		orgList = append(orgList, convertOrg(orgListFromRepo[i], projListFromRepo))
	}

	ctx.JSON(http.StatusOK, &ListOrgResponse{
		OrgList: orgList,
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

	// 1: get organization from repo
	orgFromRepo, ok := isOrgExist(ctx, controller, orgId)
	if !ok || orgFromRepo == nil {
		return
	}

	// 2: list projects from repo
	projListFromRepo, err := controller.Repo.ListProj(orgFromRepo.Id)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to list projects from repository with orgId:%d.", orgFromRepo.Id), err)
		return
	}

	ctx.JSON(http.StatusOK, &GetOrgResponse{
		Org: convertOrg(orgFromRepo, projListFromRepo),
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
	name := ctx.Query("orgName")
	orgForRepo := repository.NewOrg(name)

	_, err := controller.Repo.CreateOrg(orgForRepo)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to create organization with name:%s", name), err)
		return
	}

	ctx.JSON(http.StatusOK, &CreateOrgResponse{
		OrgId: orgForRepo.Id,
	})
}

// DeleteOrg
// @Summary Delete organization
// @Id 4
// @version 1.0
// @Tags organization
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Success 200 {object} DeleteOrgResponse
// @Router /v1/org/{orgId} [delete]
func DeleteOrg(ctx *gin.Context) {
	controller := GetController()
	orgId := utils.ToInt(ctx.Param("orgId"))

	// 1: get organization first
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 2: list projects from org
	projListFromRepo, err := controller.Repo.ListProj(orgId)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to list porojects with orgId:%d", orgId), err)
		return
	}
	if len(projListFromRepo) > 0 {
		ctx.JSON(http.StatusForbidden, rkerror.New(
			rkerror.WithHttpCode(http.StatusForbidden),
			rkerror.WithMessage("organization is not empty, please delete or migrate projects to another organization first.")))
		return
	}

	// 2: remove organization
	succ, err := controller.Repo.RemoveOrg(orgId)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("Failed to delete organization with id:%d", orgId), err)
		return
	}

	ctx.JSON(http.StatusOK, &DeleteOrgResponse{
		Status: succ,
	})
}

// UpdateOrg
// @Summary Update organization
// @Id 5
// @version 1.0
// @Tags organization
// @produce application/json
// @Param org body UpdateOrgRequest true "Organization"
// @Param orgId path int true "Organization Id"
// @Success 200 {object} UpdateOrgResponse
// @Router /v1/org/{orgId} [post]
func UpdateOrg(ctx *gin.Context) {
	controller := GetController()

	// 1: bind request
	req := &UpdateOrgRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage("Invalid request"),
			rkerror.WithDetails(err)))
		return
	}

	// 2: get organization from repo
	orgId := utils.ToInt(ctx.Param("orgId"))
	org, ok := isOrgExist(ctx, controller, orgId)
	if !ok {
		return
	}

	// 3: replace fields
	org.Name = req.Name

	// 4: update in repo
	succ, err := controller.Repo.UpdateOrg(org)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to update organization with orgId:%d", orgId), err)
		return
	}

	ctx.JSON(http.StatusOK, &UpdateOrgResponse{
		Status: succ,
	})
}

// ListProj
// @Summary List projects
// @Id 6
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Success 200 {object} ListProjResponse
// @Router /v1/org/{orgId}/proj [get]
func ListProj(ctx *gin.Context) {
	projList := make([]*Proj, 0)

	controller := GetController()
	orgId := utils.ToInt(ctx.Param("orgId"))

	// 1: get org from repo
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 2: list project from repo
	projListFromRepo, err := controller.Repo.ListProj(orgId)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to list project from repository with orgId:%d", orgId))
		return
	}

	// 3: convert to API model
	for i := range projListFromRepo {
		projList = append(projList, convertProj(projListFromRepo[i]))
	}

	ctx.JSON(http.StatusOK, &ListProjResponse{
		ProjList: projList,
	})
}

// GetProj
// @Summary Get project
// @Id 7
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Param projId path int true "Project Id"
// @Success 200 {object} GetProjResponse
// @Router /v1/org/{orgId}/proj/{projId} [get]
func GetProj(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Param("orgId"))
	projId := utils.ToInt(ctx.Param("projId"))

	// 1: get organization from repository
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 2: get project from repository
	projFromRepo, ok := isProjExist(ctx, controller, orgId, projId)
	if !ok || projFromRepo == nil {
		return
	}

	ctx.JSON(http.StatusOK, &GetProjResponse{
		Proj: convertProj(projFromRepo),
	})
}

// CreateProj
// @Summary create project
// @Id 8
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Param project body CreateProjRequest true "Project"
// @Success 200 {object} CreateProjResponse
// @Router /v1/org/{orgId}/proj [put]
func CreateProj(ctx *gin.Context) {
	controller := GetController()
	orgId := utils.ToInt(ctx.Param("orgId"))

	// 1: bind request
	req := &CreateProjRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	// 2: get organization from repository
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 3: create project
	proj := repository.NewProj(orgId, req.Name)
	_, err := controller.Repo.CreateProj(proj)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to create project with orgId:%d", orgId), err)
		return
	}

	ctx.JSON(http.StatusOK, &CreateProjResponse{
		OrgId:  orgId,
		ProjId: proj.Id,
	})
}

// DeleteProj
// @Summary delete project
// @Id 9
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Param projId path int true "Project Id"
// @Success 200 {object} DeleteProjResponse
// @Router /v1/org/{orgId}/proj/{projId} [delete]
func DeleteProj(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Param("orgId"))
	projId := utils.ToInt(ctx.Param("projId"))

	// 1: get organization from repository
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 2: remove project
	succ, err := controller.Repo.RemoveProj(orgId, projId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, orgId, projId), err)
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.ProjFailedToRemove, orgId, projId), err)
		}
		return
	}

	ctx.JSON(http.StatusOK, &DeleteProjResponse{
		Status: succ,
	})
}

// UpdateProj
// @Summary update project
// @Id 10
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId path int true "Organization Id"
// @Param projId path int true "Project Id"
// @Param project body UpdateProjRequest true "Project"
// @Success 200 {object} UpdateProjResponse
// @Router /v1/org/{orgId}/proj/{projId} [post]
func UpdateProj(ctx *gin.Context) {
	controller := GetController()

	orgId := utils.ToInt(ctx.Param("orgId"))
	projId := utils.ToInt(ctx.Param("projId"))

	// 1: bind request
	req := &UpdateProjRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	// 2: get organization from repository
	if _, ok := isOrgExist(ctx, controller, orgId); !ok {
		return
	}

	// 3: get project from repository
	projFromRepo, ok := isProjExist(ctx, controller, orgId, projId)
	if !ok || projFromRepo == nil {
		return
	}

	// 4: update values in project
	projFromRepo.Name = req.Name

	// 5: update project to repository
	succ, err := controller.Repo.UpdateProj(projFromRepo)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to update project with orgId:%d projId:%d", orgId, projId), err)
		return
	}

	ctx.JSON(http.StatusOK, &UpdateProjResponse{
		Status: succ,
	})
}

func isOrgExist(ctx *gin.Context, controller *Controller, orgId int) (*repository.Org, bool) {
	org, err := controller.Repo.GetOrg(orgId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.OrgNotFoundMsg, orgId))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.OrgFailedToGetMsg, orgId), err)
		}
		return nil, false
	}
	if org == nil {
		makeNotFoundError(ctx, fmt.Sprintf(repository.OrgNotFoundMsg, orgId))
		return nil, false
	}

	return org, true
}

func isProjExist(ctx *gin.Context, controller *Controller, orgId, projId int) (*repository.Proj, bool) {
	proj, err := controller.Repo.GetProj(orgId, projId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, orgId, projId))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.ProjFailedToGetMsg, orgId, projId), err)
		}
		return nil, false
	}
	if proj == nil {
		makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, orgId, projId))
		return nil, false
	}

	return proj, true
}
