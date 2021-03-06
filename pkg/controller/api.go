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
	ginEntry.Router.GET("/v1/proj", ListProj)
	ginEntry.Router.GET("/v1/proj/:projId", GetProj)
	ginEntry.Router.PUT("/v1/proj", CreateProj)
	ginEntry.Router.DELETE("/v1/proj/:projId", DeleteProj)
	ginEntry.Router.POST("/v1/proj/:projId", UpdateProj)

	// Source
	ginEntry.Router.PUT("/v1/source", CreateSource)
	ginEntry.Router.DELETE("/v1/source/:sourceId", DeleteSource)

	// Installations
	ginEntry.Router.GET("/v1/user/installations", ListUserInstallations)
	ginEntry.Router.GET("/v1/source/:sourceId/commits", ListCommits)
	ginEntry.Router.GET("/v1/source/:sourceId/branches", ListBranchesAndTags)

	// Pipeline templates
	ginEntry.Router.GET("/v1/pipeline/template", ListPipelineTemplate)
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

func makeAlreadyExistError(ctx *gin.Context, message string, details ...interface{}) {
	ctx.JSON(http.StatusConflict, rkerror.New(
		rkerror.WithHttpCode(http.StatusConflict),
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

func convertPipelineTemplate(templateFromRepo *repository.PipelineTemplate) *PipelineTemplate {
	template := &PipelineTemplate{
		Meta: templateFromRepo,
	}

	return template
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

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

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProj
// @Summary List projects
// @Id 6
// @version 1.0
// @Tags project
// @produce application/json
// @Param orgId query int false "Organization Id"
// @Success 200 {object} ListProjResponse
// @Router /v1/proj [get]
func ListProj(ctx *gin.Context) {
	projList := make([]*Proj, 0)

	controller := GetController()
	orgId := utils.ToInt(ctx.Query("orgId"))

	// 1: list project from repo
	projListFromRepo, err := controller.Repo.ListProj(orgId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.OrgNotFoundMsg, orgId))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.OrgFailedToGetMsg, orgId), err)
		}
		return
	}

	// 2: convert to API model
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
// @Param projId path int true "Project Id"
// @Success 200 {object} GetProjResponse
// @Router /v1/proj/{projId} [get]
func GetProj(ctx *gin.Context) {
	controller := GetController()

	projId := utils.ToInt(ctx.Param("projId"))

	// 1: get project from repository
	projFromRepo, ok := isProjExist(ctx, controller, projId)
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
// @Param project body CreateProjRequest true "Project"
// @Success 200 {object} CreateProjResponse
// @Router /v1/proj [put]
func CreateProj(ctx *gin.Context) {
	controller := GetController()
	// 1: bind request
	req := &CreateProjRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	fmt.Println(req)

	// 2: get organization from repository
	if _, ok := isOrgExist(ctx, controller, req.OrgId); !ok {
		return
	}

	// 3: create project
	proj := repository.NewProj(req.Name)
	proj.OrgId = req.OrgId
	proj.OrgName = req.OrgName
	_, err := controller.Repo.CreateProj(proj)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to create project with orgId:%d", req.OrgId), err)
		return
	}

	ctx.JSON(http.StatusOK, &CreateProjResponse{
		OrgId:  req.OrgId,
		ProjId: proj.Id,
	})
}

// DeleteProj
// @Summary delete project
// @Id 9
// @version 1.0
// @Tags project
// @produce application/json
// @Param projId path int true "Project Id"
// @Success 200 {object} DeleteProjResponse
// @Router /v1/proj/{projId} [delete]
func DeleteProj(ctx *gin.Context) {
	controller := GetController()

	projId := utils.ToInt(ctx.Param("projId"))

	// 1: remove project
	succ, err := controller.Repo.RemoveProj(projId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, projId), err)
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.ProjFailedToRemove, projId), err)
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
// @Param projId path int true "Project Id"
// @Param project body UpdateProjRequest true "Project"
// @Success 200 {object} UpdateProjResponse
// @Router /v1/proj/{projId} [post]
func UpdateProj(ctx *gin.Context) {
	controller := GetController()

	projId := utils.ToInt(ctx.Param("projId"))

	// 1: bind request
	req := &UpdateProjRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	// 3: get project from repository
	projFromRepo, ok := isProjExist(ctx, controller, projId)
	if !ok || projFromRepo == nil {
		return
	}

	// 4: update values in project
	projFromRepo.Name = req.Name

	// 5: update project to repository
	succ, err := controller.Repo.UpdateProj(projFromRepo)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to update project with projId:%d", projId), err)
		return
	}

	ctx.JSON(http.StatusOK, &UpdateProjResponse{
		Status: succ,
	})
}

// ******************************************** //
// ************** Source related ************** //
// ******************************************** //

// CreateSource
// @Summary create source
// @Id 11
// @version 1.0
// @Tags source
// @produce application/json
// @Param projId query int true "Project Id"
// @Param source body CreateSourceRequest true "Source"
// @Success 200 {object} CreateSourceResponse
// @Router /v1/source [put]
func CreateSource(ctx *gin.Context) {
	controller := GetController()
	projId := utils.ToInt(ctx.Query("projId"))

	// 1: bind request
	req := &CreateSourceRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(err.Error())))
		return
	}

	// 2: get project from repository
	if proj, ok := isProjExist(ctx, controller, projId); !ok {
		return
	} else if proj.Source != nil {
		makeAlreadyExistError(ctx, fmt.Sprintf("source already exist in project with id:%d", proj.Source.Id))
		return
	}

	// 3: create source
	src := repository.NewSource(req.Type, req.Repository)
	src.ProjId = projId
	_, err := controller.Repo.CreateSource(src)
	if err != nil {
		makeInternalError(ctx, fmt.Sprintf("failed to create source with projId:%d", projId), err)
		return
	}

	ctx.JSON(http.StatusOK, &CreateSourceResponse{
		ProjId:   src.ProjId,
		SourceId: src.Id,
	})
}

// DeleteSource
// @Summary delete source
// @Id 12
// @version 1.0
// @Tags source
// @produce application/json
// @Param sourceId path int true "Source Id"
// @Success 200 {object} DeleteProjResponse
// @Router /v1/source/{sourceId} [delete]
func DeleteSource(ctx *gin.Context) {
	controller := GetController()

	sourceId := utils.ToInt(ctx.Param("sourceId"))

	// 1: remove project
	succ, err := controller.Repo.RemoveSource(sourceId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.SourceNotFoundMsg, sourceId), err)
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.SourceFailedToRemove, sourceId), err)
		}
		return
	}

	ctx.JSON(http.StatusOK, &DeleteSourceResponse{
		Status: succ,
	})
}

// ListUserInstallations
// @Summary List user installations
// @Id 13
// @version 1.0
// @Tags installation
// @produce application/json
// @Param source query string true "Source"
// @Param user query string true "user"
// @Success 200
// @Router /v1/user/installations [get]
func ListUserInstallations(ctx *gin.Context) {
	user := ctx.Query("user")
	source := ctx.Query("source")

	var res []*Installation
	var err error

	switch source {
	case "github":
		res, err = ListUserInstallationsFromGithub(user)
	default:
		err = fmt.Errorf("unrecognized source:%s with user:%s", source, user)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage(fmt.Sprintf("failed to list user installations from %s", source)),
			rkerror.WithDetails(err)))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ListPipelineTemplate
// @Summary List pipeline templates
// @Id 14
// @version 1.0
// @Tags pipeline
// @produce application/json
// @Success 200 {object} ListPipelineTemplateResponse
// @Router /v1/pipeline/template [get]
func ListPipelineTemplate(ctx *gin.Context) {
	controller := GetController()
	templateList := make([]*PipelineTemplate, 0)

	// 1: list organization
	templateListFromRepo, err := controller.Repo.ListPipelineTemplate()
	if err != nil {
		makeInternalError(ctx, "failed to list pipeline templates", err)
		return
	}

	for i := range templateListFromRepo {
		// convert to API model
		templateList = append(templateList, convertPipelineTemplate(templateListFromRepo[i]))
	}

	ctx.JSON(http.StatusOK, &ListPipelineTemplateResponse{
		TemplateList: templateList,
	})
}

// ListCommits
// @Summary List user installation commits
// @Id 15
// @version 1.0
// @Tags installation
// @produce application/json
// @Param sourceId path int true "Source Id"
// @Param branch query string true "Branch"
// @Param perPage query int false "Number of commits per page"
// @Param page query int false "Page number to fetch"
// @Success 200
// @Router /v1/source/{sourceId}/commits [get]
func ListCommits(ctx *gin.Context) {
	controller := GetController()

	sourceId := utils.ToInt(ctx.Param("sourceId"))
	branch := ctx.Query("branch")
	perPage, page := normalizePage(utils.ToInt(ctx.Query("perPage")), utils.ToInt(ctx.Query("page")))

	// 1: get source from repository
	sourceFromRepo, ok := isSourceExist(ctx, controller, sourceId)
	if !ok || sourceFromRepo == nil {
		return
	}

	// 2: get access token from repository
	tokenFromRepo, ok := isAccessTokenExist(ctx, controller, sourceFromRepo.Type, sourceFromRepo.User)
	if !ok || sourceFromRepo == nil {
		return
	}

	// 3: List commits
	commits, err := ListCommitsFromGithub(sourceFromRepo, branch, tokenFromRepo.Token, perPage, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rkerror.New(
			rkerror.WithHttpCode(http.StatusInternalServerError),
			rkerror.WithMessage("Failed to list commits from github"),
			rkerror.WithDetails(err)))
		return
	}

	ctx.JSON(http.StatusOK, &ListCommitsResponse{
		Commits: commits,
	})
}

// ListBranchesAndTags
// @Summary List branches and tags
// @Id 16
// @version 1.0
// @Tags installation
// @produce application/json
// @Param sourceId path int true "Source Id"
// @Param perPage query int false "Number of commits per page"
// @Param page query int false "Page number to fetch"
// @Success 200
// @Router /v1/source/{sourceId}/branches [get]
func ListBranchesAndTags(ctx *gin.Context) {
	controller := GetController()

	sourceId := utils.ToInt(ctx.Param("sourceId"))
	perPage, page := normalizePage(utils.ToInt(ctx.Query("perPage")), utils.ToInt(ctx.Query("page")))

	// 1: get source from repository
	sourceFromRepo, ok := isSourceExist(ctx, controller, sourceId)
	if !ok || sourceFromRepo == nil {
		return
	}

	// 2: get access token from repository
	tokenFromRepo, ok := isAccessTokenExist(ctx, controller, sourceFromRepo.Type, sourceFromRepo.User)
	if !ok || sourceFromRepo == nil {
		return
	}

	// 3: List commits
	branches, tags, err := ListBranchesAndTagsFromGithub(sourceFromRepo, tokenFromRepo.Token, perPage, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rkerror.New(
			rkerror.WithHttpCode(http.StatusInternalServerError),
			rkerror.WithMessage("Failed to list commits from github"),
			rkerror.WithDetails(err)))
		return
	}

	ctx.JSON(http.StatusOK, &ListBranchesAndTagsResponse{
		Branches: branches,
		Tags:     tags,
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

func isProjExist(ctx *gin.Context, controller *Controller, projId int) (*repository.Proj, bool) {
	proj, err := controller.Repo.GetProj(projId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, projId))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.ProjFailedToGetMsg, projId), err)
		}
		return nil, false
	}
	if proj == nil {
		makeNotFoundError(ctx, fmt.Sprintf(repository.ProjNotFoundMsg, projId))
		return nil, false
	}

	return proj, true
}

func isSourceExist(ctx *gin.Context, controller *Controller, sourceId int) (*repository.Source, bool) {
	src, err := controller.Repo.GetSource(sourceId)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.SourceNotFoundMsg, sourceId))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.SourceFailedToGetMsg, sourceId), err)
		}
		return nil, false
	}
	if src == nil {
		makeNotFoundError(ctx, fmt.Sprintf(repository.SourceNotFoundMsg, sourceId))
		return nil, false
	}

	return src, true
}

func isAccessTokenExist(ctx *gin.Context, controller *Controller, repoType, repoUser string) (*repository.AccessToken, bool) {
	token, err := controller.Repo.GetAccessToken(repoType, repoUser)
	if err != nil {
		switch err.(type) {
		case *repository.NotFound:
			makeNotFoundError(ctx, fmt.Sprintf(repository.AccessTokenNotFoundMsg, repoType, repoUser))
		default:
			makeInternalError(ctx, fmt.Sprintf(repository.AccessTokenFailedToGetMsg, repoType, repoUser), err)
		}
		return nil, false
	}
	if token == nil {
		makeNotFoundError(ctx, fmt.Sprintf(repository.AccessTokenNotFoundMsg, repoType, repoUser))
		return nil, false
	}

	return token, true
}
