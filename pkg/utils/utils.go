// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/google/uuid"
	"strconv"
)

// CreateId will creates a random Id in string.
func CreateId() string {
	uuid.New()

	return uuid.New().String()
}

// ToInt convert string to uint
func ToInt(in string) int {
	if len(in) < 1 {
		return -1
	}

	res, err := strconv.Atoi(in)
	if err != nil {
		return -1
	}
	return res
}
