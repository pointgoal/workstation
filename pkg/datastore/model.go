// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"encoding/json"
	"github.com/pointgoal/workstation/pkg/utils"
	"gorm.io/gorm"
	"time"
)

var (
	organizationKey = Organization{}
	projectKey      = Project{}
)

// ************************************************ //
// ************** Base model related ************** //
// ************************************************ //

type Base struct {
	Id        int            `yaml:"id" json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `yaml:"createdAt" json:"createdAt"`
	UpdatedAt time.Time      `yaml:"updatedAt" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `yaml:"-" json:"-" gorm:"index"`
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// Organization defines organizations in workstation.
type Organization struct {
	Base
	Name string `yaml:"name" json:"name"`
}

// NewOrganization create a new organization with name.
// A new name will be assigned as the same as random Id if name is empty.
func NewOrganization(name string) *Organization {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	//now := time.Now()
	return &Organization{
		Name: name,
	}
}

// Equal compares organization.
// Why we need it? Because of time.Time
func (org *Organization) Equal(in *Organization) bool {
	if in == nil {
		return false
	}

	return org.Name == in.Name &&
		org.Id == in.Id &&
		org.CreatedAt.Equal(in.CreatedAt) &&
		org.UpdatedAt.Equal(in.UpdatedAt)
}

// String will marshal organization into json format.
func (org *Organization) String() string {
	bytes, _ := json.Marshal(org)
	return string(bytes)
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// Project defines projects in workstation.
type Project struct {
	Base
	OrgId int    `yaml:"orgId" json:"orgId" gorm:"index"`
	Name  string `yaml:"name" json:"name" gorm:"index"`
}

// NewProject create a project with params.
// A new name will be assigned with random Id if name is empty.
func NewProject(orgId int, name string) *Project {
	if len(name) < 1 {
		name = utils.CreateId()
	}

	return &Project{
		OrgId: orgId,
		Name:  name,
	}
}

// String will marshal organization into json format.
func (proj *Project) String() string {
	bytes, _ := json.Marshal(proj)
	return string(bytes)
}

// Equal compares project.
// Why we need it? Because of time.Time
func (proj *Project) Equal(in *Project) bool {
	if in == nil {
		return false
	}

	return proj.OrgId == in.OrgId &&
		proj.Id == in.Id &&
		proj.Name == in.Name &&
		proj.CreatedAt.Equal(in.CreatedAt) &&
		proj.UpdatedAt.Equal(in.UpdatedAt)
}
