// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrganization_WithEmptyName(t *testing.T) {
	org := NewOrganization("")

	assert.NotNil(t, org)
	assert.NotEmpty(t, org.Name)
}

func TestNewOrganization_HappyCase(t *testing.T) {
	org := NewOrganization("ut-org")

	assert.NotNil(t, org)
	assert.Equal(t, "ut-org", org.Name)
}

func TestOrganization_String(t *testing.T) {
	org := NewOrganization("ut-org")
	assert.NotEmpty(t, org.String())
}

func TestNewProject_WithEmptyName(t *testing.T) {
	proj := NewProject(1, "")

	assert.NotNil(t, proj)
	assert.NotEmpty(t, proj.Name)
}

func TestNewProject_HappyCase(t *testing.T) {
	proj := NewProject(1, "ut-proj")

	assert.NotNil(t, proj)
	assert.Equal(t, "ut-proj", proj.Name)
}

func TestProject_String(t *testing.T) {
	proj := NewProject(1, "")
	assert.NotEmpty(t, proj.String())
}
