// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestRegisterLocalFs_HappyCase(t *testing.T) {
	store := RegisterLocalFs()
	assert.NotNil(t, store)
	assert.Equal(t, EntryNameDefault, store.EntryName)
	assert.NotEmpty(t, store.EntryType)
	assert.NotEmpty(t, store.EntryDescription)
	assert.NotNil(t, store.ZapLoggerEntry)
	assert.NotNil(t, store.EventLoggerEntry)
	assert.True(t, path.IsAbs(store.RootDir))
	assert.Equal(t, LocalFsMetaFileName, store.MetaFileName)
}

func TestWithRootPathLocalFs(t *testing.T) {
	store := RegisterLocalFs(
		WithRootPathLocalFs("ut-path"))
	assert.NotNil(t, store)

	assert.Contains(t, store.RootDir, "ut-path")
}

func TestLocalFs_Connect(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterLocalFs()
	store.Connect()
}

func TestLocalFs_IsHealthy(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterLocalFs()
	store.IsHealthy()
}

func TestLocalFs_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterLocalFs()
	store.Bootstrap(context.TODO())
}

func TestLocalFs_Interrupt(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterLocalFs()
	store.Interrupt(context.TODO())
}

func TestLocalFs_GetName(t *testing.T) {
	store := RegisterLocalFs()
	assert.NotEmpty(t, store.GetName())
}

func TestLocalFs_GetType(t *testing.T) {
	store := RegisterLocalFs()
	assert.NotEmpty(t, store.GetType())
}

func TestLocalFs_GetDescription(t *testing.T) {
	store := RegisterLocalFs()
	assert.NotEmpty(t, store.GetDescription())
}

func TestLocalFs_String(t *testing.T) {
	store := RegisterLocalFs()
	assert.NotEmpty(t, store.String())
}

func TestLocalFs_Organization_Operations(t *testing.T) {
	store := RegisterLocalFs(
		WithRootPathLocalFs(t.TempDir()))

	// Empty orgs
	assert.Empty(t, store.ListOrg())

	// Insert an org
	org := NewOrganization("ut-org")
	assert.True(t, store.InsertOrg(org))
	assert.Len(t, store.ListOrg(), 1)
	assert.True(t, org.Equal(store.GetOrg(org.Id)))

	// Update org
	org.Name = "ut-org-new"
	assert.True(t, store.UpdateOrg(org))
	assert.Len(t, store.ListOrg(), 1)
	assert.True(t, org.Equal(store.GetOrg(org.Id)))
	assert.Equal(t, org.Name, store.GetOrg(org.Id).Name)

	// Remove org
	assert.True(t, store.RemoveOrg(org.Id))
	assert.Empty(t, store.ListOrg())
}

func TestLocalFs_Project_Operations(t *testing.T) {
	store := RegisterLocalFs(
		WithRootPathLocalFs(t.TempDir()))

	// Insert an org
	org := NewOrganization("ut-org")
	assert.True(t, store.InsertOrg(org))

	// Empty projects
	assert.Empty(t, store.ListProject(org.Id))

	// Insert a project
	proj := NewProject(org.Id, "ut-proj")
	assert.True(t, store.InsertProject(proj))
	assert.Len(t, store.ListProject(proj.OrgId), 1)
	assert.True(t, proj.Equal(store.GetProject(proj.OrgId, proj.Id)))

	// Update proj
	proj.Name = "ut-proj-new"
	assert.True(t, store.UpdateProject(proj))
	assert.Len(t, store.ListProject(proj.OrgId), 1)
	assert.True(t, proj.Equal(store.GetProject(proj.OrgId, proj.Id)))
	assert.Equal(t, proj.Name, store.GetProject(proj.OrgId, proj.Id).Name)

	// Remove proj
	assert.True(t, store.RemoveProject(proj.OrgId, proj.Id))
	assert.Empty(t, store.ListProject(proj.OrgId))
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
