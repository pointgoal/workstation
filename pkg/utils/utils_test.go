// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateId_HappyCase(t *testing.T) {
	assert.NotEmpty(t, CreateId())
}

func TestToUint(t *testing.T) {
	// For empty string
	assert.Equal(t, -1, ToInt(""))

	// For invalid number
	assert.Equal(t, -1, ToInt("invalid"))

	// For happy case
	assert.Equal(t, 99, ToInt("99"))
}
