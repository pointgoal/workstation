// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/google/uuid"
)

// CreateId will creates a random Id in string.
func CreateId() string {
	return uuid.New().String()
}
