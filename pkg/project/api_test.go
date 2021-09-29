// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"github.com/gin-gonic/gin"
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

func TestListOrgs(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	ListOrgs(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestGetOrgs(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	GetOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestCreateOrgs(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	CreateOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestDeleteOrgs(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	DeleteOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestUpdateOrgs(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)

	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}

	UpdateOrg(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestListProjects(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	ListProjects(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestGetProject(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	GetProject(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestCreateProject(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}

	CreateProject(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestDeleteProject(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)
	DeleteProject(ctx)

	assert.Equal(t, 200, writer.StatusCode)
}

func TestUpdateProject(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	RegisterEntry()

	writer := &httptest.TestResponseWriter{}
	ctx, _ := gin.CreateTestContext(writer)

	stringReader := strings.NewReader(`{name:"ut-name"}`)
	ctx.Request = &http.Request{
		Body: ioutil.NopCloser(stringReader),
		URL:  &url.URL{},
	}

	UpdateProject(ctx)

	assert.Equal(t, 200, writer.StatusCode)
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
