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
