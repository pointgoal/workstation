// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"github.com/pointgoal/workstation/pkg/utils"
	"gorm.io/gorm"
	"time"
)

var (
	orgKey  = &Org{}
	projKey = &Proj{}
)

// ************************************************ //
// ************** Base model related ************** //
// ************************************************ //

// Base defines base model of gorm model
type Base struct {
	Id        int            `yaml:"id" json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `yaml:"createdAt" json:"createdAt"`
	UpdatedAt time.Time      `yaml:"updatedAt" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `yaml:"-" json:"-" gorm:"index"`
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// Org defines organizations in workstation.
type Org struct {
	Base
	Name     string  `yaml:"name" json:"name"`
	ProjList []*Proj `yaml:"-" json:"-"`
}

// NewOrg create a new organization with name.
// A new name will be assigned as the same as random Id if name is empty.
func NewOrg(name string) *Org {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	//now := time.Now()
	return &Org{
		Name: name,
	}
}

// Equal compares organization.
// Why we need it? Because of time.Time
func (org *Org) Equal(in *Org) bool {
	if in == nil {
		return false
	}

	return org.Name == in.Name &&
		org.Id == in.Id &&
		org.CreatedAt.Equal(in.CreatedAt) &&
		org.UpdatedAt.Equal(in.UpdatedAt)
}

// String will marshal organization into json format.
func (org *Org) String() string {
	bytes, _ := json.Marshal(org)
	return string(bytes)
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// Proj defines projects in workstation.
type Proj struct {
	Base
	OrgId int    `yaml:"orgId" json:"orgId" gorm:"index"`
	Name  string `yaml:"name" json:"name" gorm:"index"`
}

// NewProject create a project with params.
// A new name will be assigned with random Id if name is empty.
func NewProj(orgId int, name string) *Proj {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	return &Proj{
		OrgId: orgId,
		Name:  name,
	}
}

// String will marshal organization into json format.
func (proj *Proj) String() string {
	bytes, _ := json.Marshal(proj)
	return string(bytes)
}

// Equal compares project.
// Why we need it? Because of time.Time
func (proj *Proj) Equal(in *Proj) bool {
	if in == nil {
		return false
	}

	return proj.OrgId == in.OrgId &&
		proj.Id == in.Id &&
		proj.Name == in.Name &&
		proj.CreatedAt.Equal(in.CreatedAt) &&
		proj.UpdatedAt.Equal(in.UpdatedAt)
}

// ******************************************** //
// ************** Source related ************** //
// ******************************************** //

type Source struct {
	Base
	OrgId  int     `yaml:"orgId" json:"orgId" gorm:"index"`
	ProjId int     `yaml:"projId" json:"projId" gorm:"index"`
	Type   string  `yaml:"type" json:"type" gorm:"index"`
	Github *Github `yaml:"github" json:"github"`
	Local  *Local  `yaml:"local" json:"local"`
}

type Github struct {
	Base
	Repository  string `yaml:"repository" json:"repository"`
	AccessToken string `yaml:"accessToken" json:"accessToken"`
}

type Local struct {
	Base
	FullPath string `yaml:"fullPath" json:"fullPath"`
}
