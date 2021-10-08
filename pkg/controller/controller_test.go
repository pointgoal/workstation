// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

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
