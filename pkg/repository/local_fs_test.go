// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

//import (
//	"context"
//	"github.com/stretchr/testify/assert"
//	"path"
//	"testing"
//)
//
//func TestRegisterLocalFs_HappyCase(t *testing.T) {
//	repo := RegisterLocalFs()
//	assert.NotNil(t, repo)
//	assert.Equal(t, EntryNameDefault, repo.EntryName)
//	assert.NotEmpty(t, repo.EntryType)
//	assert.NotEmpty(t, repo.EntryDescription)
//	assert.NotNil(t, repo.ZapLoggerEntry)
//	assert.NotNil(t, repo.EventLoggerEntry)
//	assert.True(t, path.IsAbs(repo.RootDir))
//	assert.Equal(t, LocalFsMetaFileName, repo.MetaFileName)
//}
//
//func TestWithRootPathLocalFs(t *testing.T) {
//	repo := RegisterLocalFs(
//		WithRootPathLocalFs("ut-path"))
//	assert.NotNil(t, repo)
//
//	assert.Contains(t, repo.RootDir, "ut-path")
//}
//
//func TestLocalFs_Connect(t *testing.T) {
//	defer assertNotPanic(t)
//
//	repo := RegisterLocalFs()
//	repo.Connect()
//}
//
//func TestLocalFs_IsHealthy(t *testing.T) {
//	defer assertNotPanic(t)
//
//	repo := RegisterLocalFs()
//	repo.IsHealthy()
//}
//
//func TestLocalFs_Bootstrap(t *testing.T) {
//	defer assertNotPanic(t)
//
//	repo := RegisterLocalFs()
//	repo.Bootstrap(context.TODO())
//}
//
//func TestLocalFs_Interrupt(t *testing.T) {
//	defer assertNotPanic(t)
//
//	repo := RegisterLocalFs()
//	repo.Interrupt(context.TODO())
//}
//
//func TestLocalFs_GetName(t *testing.T) {
//	repo := RegisterLocalFs()
//	assert.NotEmpty(t, repo.GetName())
//}
//
//func TestLocalFs_GetType(t *testing.T) {
//	repo := RegisterLocalFs()
//	assert.NotEmpty(t, repo.GetType())
//}
//
//func TestLocalFs_GetDescription(t *testing.T) {
//	repo := RegisterLocalFs()
//	assert.NotEmpty(t, repo.GetDescription())
//}
//
//func TestLocalFs_String(t *testing.T) {
//	repo := RegisterLocalFs()
//	assert.NotEmpty(t, repo.String())
//}
//
//func TestLocalFs_Organization_Operations(t *testing.T) {
//	repo := RegisterLocalFs(
//		WithRootPathLocalFs(t.TempDir()))
//
//	// empty orgs
//	orgList, err := repo.ListOrg()
//	assert.Nil(t, err)
//	assert.Empty(t, orgList)
//
//	// create an org
//	org := NewOrg("ut-org")
//	succ, err := repo.CreateOrg(org)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//
//	// Update org
//	org.Name = "ut-org-new"
//	succ, err = repo.UpdateOrg(org)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//
//	// Remove org
//	succ, err = repo.RemoveOrg(org.Id)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//}
//
//func TestLocalFs_Project_Operations(t *testing.T) {
//	repo := RegisterLocalFs(
//		WithRootPathLocalFs(t.TempDir()))
//
//	// create an org
//	org := NewOrg("ut-org")
//	succ, err := repo.CreateOrg(org)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//
//	// empty projects
//	projList, err := repo.ListProj(org.Id)
//	assert.Empty(t, projList)
//	assert.Nil(t, err)
//
//	// create a project
//	proj := NewProj(org.Id, "ut-proj")
//	succ, err = repo.CreateProj(proj)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//
//	// update proj
//	proj.Name = "ut-proj-new"
//	succ, err = repo.UpdateProj(proj)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//
//	// remove proj
//	succ, err = repo.RemoveProj(proj.OrgId, proj.Id)
//	assert.True(t, succ)
//	assert.Nil(t, err)
//}
//
//func assertNotPanic(t *testing.T) {
//	if r := recover(); r != nil {
//		// Expect panic to be called with non nil error
//		assert.True(t, false)
//	} else {
//		// This should never be called in case of a bug
//		assert.True(t, true)
//	}
//}
//
//func assertPanic(t *testing.T) {
//	if r := recover(); r != nil {
//		// Expect panic to be called with non nil error
//		assert.True(t, true)
//	} else {
//		// This should never be called in case of a bug
//		assert.True(t, false)
//	}
//}
