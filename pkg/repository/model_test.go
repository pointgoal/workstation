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
	proj := NewProj(1, "")

	assert.NotNil(t, proj)
	assert.NotEmpty(t, proj.Name)
}

func TestNewProj_HappyCase(t *testing.T) {
	proj := NewProj(1, "ut-proj")

	assert.NotNil(t, proj)
	assert.Equal(t, "ut-proj", proj.Name)
}

func TestProject_String(t *testing.T) {
	proj := NewProj(1, "")
	assert.NotEmpty(t, proj.String())
}
