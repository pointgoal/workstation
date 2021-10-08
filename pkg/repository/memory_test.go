// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterMemory_HappyCase(t *testing.T) {
	repo := RegisterMemory()
	assert.NotNil(t, repo)
	assert.Equal(t, EntryNameDefault, repo.EntryName)
	assert.NotEmpty(t, repo.EntryType)
	assert.NotEmpty(t, repo.EntryDescription)
	assert.NotNil(t, repo.ZapLoggerEntry)
	assert.NotNil(t, repo.EventLoggerEntry)
}

func TestMemory_Connect(t *testing.T) {
	defer assertNotPanic(t)

	repo := RegisterMemory()
	assert.Nil(t, repo.Connect())
}

func TestMemory_IsHealthy(t *testing.T) {
	defer assertNotPanic(t)

	repo := RegisterMemory()
	assert.True(t, repo.IsHealthy())
}

func TestMemory_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)

	repo := RegisterMemory()
	repo.Bootstrap(context.TODO())
}

func TestMemory_Interrupt(t *testing.T) {
	defer assertNotPanic(t)

	repo := RegisterMemory()
	repo.Interrupt(context.TODO())
}

func TestMemory_GetName(t *testing.T) {
	repo := RegisterMemory()
	assert.NotEmpty(t, repo.GetName())
}

func TestMemory_GetType(t *testing.T) {
	repo := RegisterMemory()
	assert.NotEmpty(t, repo.GetType())
}

func TestMemory_GetDescription(t *testing.T) {
	repo := RegisterMemory()
	assert.NotEmpty(t, repo.GetDescription())
}

func TestMemory_String(t *testing.T) {
	repo := RegisterMemory()
	assert.NotEmpty(t, repo.String())
}

func TestMemory_Organization_Operations(t *testing.T) {
	repo := RegisterMemory()

	// empty orgs
	orgList, err := repo.ListOrg()
	assert.Nil(t, err)
	assert.Empty(t, orgList)

	// create an org
	org := NewOrg("ut-org")
	succ, err := repo.CreateOrg(org)
	assert.True(t, succ)
	assert.Nil(t, err)

	// Update org
	org.Name = "ut-org-new"
	succ, err = repo.UpdateOrg(org)
	assert.True(t, succ)
	assert.Nil(t, err)

	// Remove org
	succ, err = repo.RemoveOrg(org.Id)
	assert.True(t, succ)
	assert.Nil(t, err)
}

func TestMemory_Project_Operations(t *testing.T) {
	repo := RegisterMemory()

	// create an org
	org := NewOrg("ut-org")
	succ, err := repo.CreateOrg(org)
	assert.True(t, succ)
	assert.Nil(t, err)

	// empty projects
	projList, err := repo.ListProj(org.Id)
	assert.Empty(t, projList)
	assert.Nil(t, err)

	// create a project
	proj := NewProj(org.Id, "ut-proj")
	succ, err = repo.CreateProj(proj)
	assert.True(t, succ)
	assert.Nil(t, err)

	// update proj
	proj.Name = "ut-proj-new"
	succ, err = repo.UpdateProj(proj)
	assert.True(t, succ)
	assert.Nil(t, err)

	// remove proj
	succ, err = repo.RemoveProj(proj.OrgId, proj.Id)
	assert.True(t, succ)
	assert.Nil(t, err)
}
