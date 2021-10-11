// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrg_WithEmptyName(t *testing.T) {
	org := NewOrg("")

	assert.NotNil(t, org)
	assert.NotEmpty(t, org.Name)
}

func TestNewOrg_HappyCase(t *testing.T) {
	org := NewOrg("ut-org")

	assert.NotNil(t, org)
	assert.Equal(t, "ut-org", org.Name)
}

func TestNewOrg_String(t *testing.T) {
	org := NewOrg("ut-org")
	assert.NotEmpty(t, org.String())
}

func TestNewProj_WithEmptyName(t *testing.T) {
	proj := NewProj("")

	assert.NotNil(t, proj)
	assert.NotEmpty(t, proj.Name)
}

func TestNewProj_HappyCase(t *testing.T) {
	proj := NewProj("ut-proj")

	assert.NotNil(t, proj)
	assert.Equal(t, "ut-proj", proj.Name)
}

func TestProj_String(t *testing.T) {
	proj := NewProj("")
	assert.NotEmpty(t, proj.String())
}

func TestNewSource_HappyCase(t *testing.T) {
	src := NewSource("github", "ut/ut-repo")
	src.ProjId = 1
	assert.NotNil(t, src)
	assert.Equal(t, 1, src.ProjId)
	assert.Equal(t, "github", src.Type)
	assert.Equal(t, "ut/ut-repo", src.Repository)
}

func TestSource_String(t *testing.T) {
	src := NewSource("github", "ut/ut-repo")
	assert.NotEmpty(t, src.String())
}
