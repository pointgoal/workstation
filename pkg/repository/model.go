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
	orgKey         = &Org{}
	projKey        = &Proj{}
	sourceKey      = &Source{}
	accessTokenKey = &AccessToken{}
)

// ************************************************ //
// ************** Base model related ************** //
// ************************************************ //

// Base defines base model of gorm model
type Base struct {
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
	Id       int     `yaml:"id" json:"id" gorm:"primaryKey"`
	Name     string  `yaml:"name" json:"name"`
	ProjList []*Proj `yaml:"-" json:"-"`
}

// NewOrg create a new organization with name.
// A new name will be assigned as the same as random Id if name is empty.
func NewOrg(name string) *Org {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	return &Org{
		Name: name,
	}
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
	Id      int     `yaml:"id" json:"id" gorm:"primaryKey"`
	OrgId   int     `yaml:"orgId" json:"orgId" gorm:"index"`
	OrgName string  `yaml:"orgName" json:"orgName" gorm:"index"`
	Name    string  `yaml:"name" json:"name" gorm:"index"`
	Source  *Source `yaml:"source" json:"source"`
}

// NewProject create a project with params.
// A new name will be assigned with random Id if name is empty.
func NewProj(name string) *Proj {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	return &Proj{
		Name: name,
	}
}

// String will marshal project into json format.
func (proj *Proj) String() string {
	bytes, _ := json.Marshal(proj)
	return string(bytes)
}

// ******************************************** //
// ************** Source related ************** //
// ******************************************** //

type Source struct {
	Base
	Id         int    `yaml:"id" json:"id" gorm:"primaryKey"`
	ProjId     int    `yaml:"projId" json:"projId" gorm:"index"`
	Type       string `yaml:"type" json:"type" gorm:"index"`
	Repository string `yaml:"repository" json:"repository"`
	User       string `yaml:"user" json:"user"`
}

// NewSource create a project with params.
// A new name will be assigned with random Id if name is empty.
func NewSource(repoType, repository string) *Source {
	return &Source{
		Type:       repoType,
		Repository: repository,
	}
}

// String will marshal source into json format.
func (src *Source) String() string {
	bytes, _ := json.Marshal(src)
	return string(bytes)
}

// ************************************************* //
// ************** AccessToken related ************** //
// ************************************************* //

type AccessToken struct {
	Base
	Id    int    `yaml:"id" json:"id" gorm:"primaryKey"`
	Type  string `yaml:"type" json:"type" gorm:"index"`
	User  string `yaml:"user" json:"user"`
	Token string `yaml:"-" json:"-"`
}

// NewAccessToken create a project with params.
// A new name will be assigned with random Id if name is empty.
func NewAccessToken(repoType, repoUser, repoToken string) *AccessToken {
	return &AccessToken{
		Type:  repoType,
		User:  repoUser,
		Token: repoToken,
	}
}

// String will marshal token into json format.
func (token *AccessToken) String() string {
	bytes, _ := json.Marshal(token)
	return string(bytes)
}

// ************************************************* //
// ************** PipelineTemplate related ************** //
// ************************************************* //

type PipelineTemplate struct {
	Base
	Id       int    `yaml:"id" json:"id" gorm:"primaryKey"`
	Name     string `yaml:"name" json:"name"`
	Language string `yaml:"language" json:"language"`
	Content  string `yaml:"content" json:"content"`
}
