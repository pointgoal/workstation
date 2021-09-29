// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"context"
	rkentry "github.com/rookie-ninja/rk-entry/entry"
	rkgin "github.com/rookie-ninja/rk-gin/boot"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegisterEntriesFromConfig(t *testing.T) {
	bootConfigStr := `
project:
  enabled: true
`

	tempDir := path.Join(t.TempDir(), "boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(bootConfigStr), os.ModePerm))
	entries := RegisterEntriesFromConfig(tempDir)

	assert.NotEmpty(t, entries)
}

func TestRegisterEntry(t *testing.T) {
	assert.NotNil(t, RegisterEntry())
}

func TestEntry_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	entry := RegisterEntry()

	entry.Bootstrap(context.TODO())
}

func TestEntry_Interrupt(t *testing.T) {
	defer assertNotPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("workstation")
	rkgin.RegisterGinEntry(rkgin.WithNameGin("workstation"))
	entry := RegisterEntry()

	entry.Bootstrap(context.TODO())
	entry.Interrupt(context.TODO())
}

func TestEntry_GetName(t *testing.T) {
	entry := RegisterEntry()
	assert.Equal(t, EntryName, entry.GetName())
}

func TestEntry_GetDescription(t *testing.T) {
	entry := RegisterEntry()
	assert.Equal(t, EntryDescription, entry.GetDescription())
}

func TestEntry_GetType(t *testing.T) {
	entry := RegisterEntry()
	assert.Equal(t, EntryType, entry.GetType())
}

func TestEntry_String(t *testing.T) {
	entry := RegisterEntry()
	assert.NotEmpty(t, entry.String())
}

func TestEntry_ListOrgs(t *testing.T) {
	entry := RegisterEntry()
	entry.AddOrg("ut-org")
	assert.Len(t, entry.ListOrgs(), 1)
}

func TestEntry_AddOrg(t *testing.T) {
	entry := RegisterEntry()
	entry.AddOrg("ut-org")
	assert.Len(t, entry.ListOrgs(), 1)
}

func TestEntry_GetOrg(t *testing.T) {
	entry := RegisterEntry()
	entry.AddOrg("ut-org")
	assert.Equal(t, "ut-org", entry.ListOrgs()[0].Name)
}

func TestDeleteOrg(t *testing.T) {
	entry := RegisterEntry()
	orgId := entry.AddOrg("ut-org")
	assert.True(t, entry.DeleteOrg(orgId))
	assert.False(t, entry.DeleteOrg("fake-org"))
}

func TestEntry_UpdateOrg(t *testing.T) {
	entry := RegisterEntry()
	entry.AddOrg("ut-org")

	// with nil org
	assert.False(t, entry.UpdateOrg(nil))

	// with non-exist org
	assert.False(t, entry.UpdateOrg(&Organization{}))

	// expect true
	org := entry.ListOrgs()[0]
	org.Name = "ut-org-2"
	assert.True(t, entry.UpdateOrg(org))
}

func TestEntry_ListProjects(t *testing.T) {
	entry := RegisterEntry()

	// non-exist org
	assert.Empty(t, entry.ListProjects("fake-org-id"))

	// happy case
	orgId := entry.AddOrg("ut-org")
	projId := entry.AddProject(orgId, NewProject(""))
	assert.NotEmpty(t, projId)
}

func TestEntry_GetProject(t *testing.T) {
	entry := RegisterEntry()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := NewProject("")
	projId := entry.AddProject(orgId, proj)
	assert.NotEmpty(t, projId)

	assert.Equal(t, proj, entry.GetProject(orgId, projId))
}

func TestEntry_DeleteProject(t *testing.T) {
	entry := RegisterEntry()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := NewProject("")
	projId := entry.AddProject(orgId, proj)
	assert.NotEmpty(t, projId)

	assert.Equal(t, proj, entry.GetProject(orgId, projId))

	entry.DeleteProject(orgId, projId)
	assert.Empty(t, entry.ListProjects(orgId))
}

func TestEntry_UpdateProject(t *testing.T) {
	entry := RegisterEntry()

	// happy case
	orgId := entry.AddOrg("ut-org")
	proj := NewProject("")
	projId := entry.AddProject(orgId, proj)
	assert.NotEmpty(t, projId)

	entry.GetProject(orgId, projId).Name = "ut-project"
	assert.True(t, entry.UpdateProject(orgId, entry.GetProject(orgId, projId)))
}
