// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"github.com/pointgoal/workstation/pkg/datastore"
	rkentry "github.com/rookie-ninja/rk-entry/entry"
	rkgin "github.com/rookie-ninja/rk-gin/boot"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegisterControllerFromConfig(t *testing.T) {
	bootConfigStr := `
controller:
  enabled: true
`

	tempDir := path.Join(t.TempDir(), "boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(bootConfigStr), os.ModePerm))
	entries := RegisterControllerFromConfig(tempDir)

	assert.NotEmpty(t, entries)
}

func TestRegisterController(t *testing.T) {
	assert.NotNil(t, RegisterController())
}

func TestController_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	entry := RegisterController()

	entry.Bootstrap(context.TODO())
}

func TestController_Interrupt(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	entry := RegisterController()

	entry.Bootstrap(context.TODO())
	entry.Interrupt(context.TODO())
}

func TestController_GetName(t *testing.T) {
	entry := RegisterController()
	assert.Equal(t, EntryName, entry.GetName())
}

func TestController_GetDescription(t *testing.T) {
	entry := RegisterController()
	assert.Equal(t, EntryDescription, entry.GetDescription())
}

func TestController_GetType(t *testing.T) {
	entry := RegisterController()
	assert.Equal(t, EntryType, entry.GetType())
}

func TestController_String(t *testing.T) {
	entry := RegisterController()
	assert.NotEmpty(t, entry.String())
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

func TestController_Organization_Operations(t *testing.T) {
	store := datastore.RegisterMemory()

	// Empty orgs
	assert.Empty(t, store.ListOrg())

	// Insert an org
	org := datastore.NewOrganization("ut-org")
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

func TestController_ListOrg(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()
	assert.NotEmpty(t, entry.AddOrg("ut-org"))
	assert.Len(t, entry.ListOrg(), 1)
}

func TestController_AddOrg(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()
	entry.AddOrg("ut-org")
	assert.Len(t, entry.ListOrg(), 1)
}

func TestController_GetOrg(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()
	entry.AddOrg("ut-org")
	assert.Equal(t, "ut-org", entry.ListOrg()[0].Name)
}

func TestController_DeleteOrg(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()
	orgId := entry.AddOrg("ut-org")
	assert.True(t, entry.DeleteOrg(orgId))
	assert.False(t, entry.DeleteOrg(orgId+1))
}

func TestController_UpdateOrg(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()
	entry.AddOrg("ut-org")

	// with nil org
	assert.False(t, entry.UpdateOrg(nil))

	// with non-exist org
	assert.False(t, entry.UpdateOrg(&datastore.Organization{}))

	// expect true
	org := entry.ListOrg()[0]
	org.Name = "ut-org-2"
	assert.True(t, entry.UpdateOrg(org))
}

func TestController_ListProjects(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()

	// non-exist org
	assert.Empty(t, entry.ListProject(10000))

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := datastore.NewProject(orgId, "")
	projId := entry.AddProject(proj)
	assert.NotEmpty(t, projId)
}

func TestController_GetProject(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := datastore.NewProject(orgId, "")
	projId := entry.AddProject(proj)
	assert.NotEmpty(t, projId)

	assert.True(t, proj.Equal(entry.GetProject(orgId, projId)))
}

func TestController_DeleteProject(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := datastore.NewProject(orgId, "")
	projId := entry.AddProject(proj)
	assert.NotEmpty(t, projId)

	assert.True(t, proj.Equal(entry.GetProject(orgId, projId)))

	entry.RemoveProject(orgId, projId)
	assert.Empty(t, entry.ListProject(orgId))
}

func TestController_UpdateProject(t *testing.T) {
	datastore.RegisterMemory()
	entry := RegisterController()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := datastore.NewProject(orgId, "")
	projId := entry.AddProject(proj)
	assert.NotEmpty(t, projId)

	entry.GetProject(orgId, projId).Name = "ut-project"
	assert.True(t, entry.UpdateProject(entry.GetProject(orgId, projId)))
}
