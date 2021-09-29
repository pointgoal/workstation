// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"fmt"
	"github.com/pointgoal/workstation/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewOrganization_WithEmptyName(t *testing.T) {
	org := NewOrganization("")

	assert.NotNil(t, org)
	assert.NotEmpty(t, org.Id)
	assert.NotEmpty(t, org.Name)
	assert.Equal(t, org.Id, org.Name)
	assert.NotNil(t, org.CreateAt)
	assert.NotNil(t, org.UpdateAt)
	assert.Empty(t, org.projectMap)
}

func TestNewOrganization_HappyCase(t *testing.T) {
	org := NewOrganization("ut-org")

	assert.NotNil(t, org)
	assert.NotEmpty(t, org.Id)
	assert.Equal(t, "ut-org", org.Name)
	assert.NotNil(t, org.CreateAt)
	assert.NotNil(t, org.UpdateAt)
	assert.Empty(t, org.projectMap)
}

func TestOrganization_UpsertProject_WithNilProject(t *testing.T) {
	org := NewOrganization("ut-org")
	assert.False(t, org.UpsertProject(nil))
	assert.Empty(t, org.projectMap)
}

func TestOrganization_UpsertProject_WithInsertCase(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	assert.True(t, org.UpsertProject(proj))
	assert.True(t, org.HasProject(proj.Id))
	assert.Equal(t, proj, org.projectMap[proj.Id])
}

func TestOrganization_UpsertProject_WithUpdateCase(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	assert.True(t, org.UpsertProject(proj))
	assert.True(t, org.HasProject(proj.Id))
	assert.Equal(t, proj, org.projectMap[proj.Id])

	// Change project and update again
	proj.Name = "ut-proj-new"

	assert.True(t, org.UpsertProject(proj))
	assert.True(t, org.HasProject(proj.Id))
	assert.Equal(t, proj, org.projectMap[proj.Id])
}

func TestOrganization_RemoveProject_WithNonExist(t *testing.T) {
	org := NewOrganization("ut-org")
	assert.False(t, org.RemoveProject("fake-id"))
}

func TestOrganization_RemoveProject_HappyCase(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.True(t, org.RemoveProject(proj.Id))
}

func TestOrganization_HasProject_ExpectTrue(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.True(t, org.HasProject(proj.Id))
}

func TestOrganization_HasProject_ExpectFalse(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.False(t, org.HasProject("fake-id"))
}

func TestOrganization_ListProjects(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.Len(t, org.ListProjects(), 1)
}

func TestOrganization_GetProject(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.Equal(t, proj, org.GetProject(proj.Id))
}

func TestOrganization_String(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	assert.NotEmpty(t, org.String())

	fmt.Println(org.String())
}

func TestOrganization_MarshalJSON(t *testing.T) {
	org := NewOrganization("ut-org")
	proj := &Project{
		Id:       utils.CreateId(),
		Name:     "ut-proj",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	assert.True(t, org.UpsertProject(proj))

	bytes, err := org.MarshalJSON()
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
}

func TestOrganization_UnmarshalJSON(t *testing.T) {
	org := NewOrganization("ut-org")
	assert.Nil(t, org.UnmarshalJSON(nil))
}
