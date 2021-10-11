// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegisterDataStoreFromConfig(t *testing.T) {
	// For memory
	bootConfigStr := `
repository:
  enabled: true
  provider: memory
`

	tempDir := path.Join(t.TempDir(), "boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(bootConfigStr), os.ModePerm))
	stores := RegisterRepositoryFromConfig(tempDir)

	assert.NotEmpty(t, stores)

	// For mySql
	bootConfigStr = `
repository:
  enabled: true
  provider: mySql
`

	tempDir = path.Join(t.TempDir(), "boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(bootConfigStr), os.ModePerm))
	stores = RegisterRepositoryFromConfig(tempDir)

	assert.NotEmpty(t, stores)
}

func TestGetDataStore(t *testing.T) {
	bootConfigStr := `
repository:
  enabled: true
  provider: memory
`

	tempDir := path.Join(t.TempDir(), "boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(bootConfigStr), os.ModePerm))
	RegisterRepositoryFromConfig(tempDir)

	assert.NotNil(t, GetRepository())
}
