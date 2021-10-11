// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pointgoal/workstation/pkg/repository"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-gin/boot"
	"github.com/stretchr/testify/assert"
	httptest "github.com/stretchr/testify/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestInitApi_WithNilGinEntry(t *testing.T) {
	defer assertPanic(t)
	initApi()
}

func TestInitApi_HappyCase(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")

	entry := rkgin.RegisterGinEntry(
		rkgin.WithNameGin("workstation"))

	initApi()

	assert.True(t, len(entry.Router.Routes()) > 0)
}

func TestListOrg(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repository.RegisterMemory()
	RegisterController()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	ListOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestGetOrg(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// with invalid org id, 404 expected
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	GetOrg(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "orgId",
		Value: "1",
	})
	GetOrg(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestCreateOrg(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repository.RegisterMemory()
	RegisterController()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	CreateOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestDeleteOrg(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	DeleteOrg(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "orgId",
		Value: "1",
	})
	DeleteOrg(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestUpdateOrg(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	UpdateOrg(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	stringReader = strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "orgId",
		Value: "1",
	})
	UpdateOrg(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestListProj(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 200
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "orgId=1;",
		},
	}
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	ListProj(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestGetProj(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404 without org
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	GetProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 404 without project
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "orgId",
		Value: "1",
	})
	GetProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	ctx.Request = &http.Request{
		URL: &url.URL{},
	}
	repo.CreateProj(&repository.Proj{
		OrgId: 1,
		Name:  "ut-org",
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "projId",
		Value: "1",
	})
	GetProj(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestCreateProj(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	CreateProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	stringReader = strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL: &url.URL{
			RawQuery: "orgId=1;",
		},
	}
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	CreateProj(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestDeleteProj(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404 without proj
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	DeleteProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	repo.CreateProj(&repository.Proj{
		Name:  "ut-proj",
		OrgId: 1,
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "projId",
		Value: "1",
	})
	DeleteProj(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestUpdateProj(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	repo := repository.RegisterMemory()
	RegisterController()

	// expect 404 without org
	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	UpdateProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 404 without proj
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	stringReader = strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	repo.CreateOrg(&repository.Org{
		Name: "ut-org",
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "projId",
		Value: "1",
	})
	UpdateProj(ctx)
	assert.Equal(t, http.StatusNotFound, writer.StatusCode)

	// expect 200
	writer = &httptest.TestResponseWriter{}
	ctx, _ = gin.CreateTestContext(writer)
	stringReader = strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}
	repo.CreateProj(&repository.Proj{
		Name:  "ut-proj",
		OrgId: 1,
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "projId",
		Value: "1",
	})
	UpdateProj(ctx)
	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func assertNotPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, false)
	} else {
		// This should never be called in case of a bug
		assert.True(t, true)
	}
}

func assertPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, true)
	} else {
		// This should never be called in case of a bug
		assert.True(t, false)
	}
}
