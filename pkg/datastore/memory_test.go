// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterMemory_HappyCase(t *testing.T) {
	store := RegisterMemory()
	assert.NotNil(t, store)
	assert.Equal(t, EntryNameDefault, store.EntryName)
	assert.NotEmpty(t, store.EntryType)
	assert.NotEmpty(t, store.EntryDescription)
	assert.NotNil(t, store.ZapLoggerEntry)
	assert.NotNil(t, store.EventLoggerEntry)
}

func TestMemory_Connect(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterMemory()
	assert.Nil(t, store.Connect())
}

func TestMemory_IsHealthy(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterMemory()
	assert.True(t, store.IsHealthy())
}

func TestMemory_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterMemory()
	store.Bootstrap(context.TODO())
}

func TestMemory_Interrupt(t *testing.T) {
	defer assertNotPanic(t)

	store := RegisterMemory()
	store.Interrupt(context.TODO())
}

func TestMemory_GetName(t *testing.T) {
	store := RegisterMemory()
	assert.NotEmpty(t, store.GetName())
}

func TestMemory_GetType(t *testing.T) {
	store := RegisterMemory()
	assert.NotEmpty(t, store.GetType())
}

func TestMemory_GetDescription(t *testing.T) {
	store := RegisterMemory()
	assert.NotEmpty(t, store.GetDescription())
}

func TestMemory_String(t *testing.T) {
	store := RegisterMemory()
	assert.NotEmpty(t, store.String())
}

func TestMemory_Organization_Operations(t *testing.T) {
	store := RegisterMemory()

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

func TestMemory_Project_Operations(t *testing.T) {
	store := RegisterMemory()

	// Insert an org
	org := NewOrganization("ut-org")
	assert.True(t, store.InsertOrg(org))

	// Empty projects
	//assert.Empty(t, store.ListProjects(org.Id))

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
