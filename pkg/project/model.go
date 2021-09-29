// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"encoding/json"
	"github.com/pointgoal/workstation/pkg/utils"
	"time"
)

// Organization defines organizations in workstation.
type Organization struct {
	Id         string              `yaml:"id" json:"id"`
	Name       string              `yaml:"name" json:"name"`
	CreateAt   time.Time           `yaml:"createAt" json:"createAt"`
	UpdateAt   time.Time           `yaml:"updateAt" json:"updateAt"`
	projectMap map[string]*Project `yaml:"-" json:"-"`
}

// NewOrganization create a new organization with name.
// A new name will be assigned as the same is random organization Id if name is empty.
func NewOrganization(name string) *Organization {
	id := utils.CreateId()

	if len(name) < 1 {
		name = id
	}

	now := time.Now()

	return &Organization{
		Id:         id,
		Name:       name,
		CreateAt:   now,
		UpdateAt:   now,
		projectMap: make(map[string]*Project, 0),
	}
}

// GetProject will returns project in map.
func (org *Organization) GetProject(id string) *Project {
	return org.projectMap[id]
}

// UpsertProject will update or add new project in organization.
func (org *Organization) UpsertProject(proj *Project) bool {
	if proj == nil {
		return false
	}

	org.projectMap[proj.Id] = proj

	return true
}

// RemoveProject will remove project in organization.
func (org *Organization) RemoveProject(id string) bool {
	if !org.HasProject(id) {
		return false
	}

	delete(org.projectMap, id)

	return true
}

// HasProject will check whether project in organization with Id.
func (org *Organization) HasProject(id string) bool {
	_, contains := org.projectMap[id]

	return contains
}

func (org *Organization) ListProjects() []*Project {
	res := make([]*Project, 0)

	for _, v := range org.projectMap {
		res = append(res, v)
	}

	return res
}

// String will marshal organization into json format.
func (org *Organization) String() string {
	bytes, _ := json.Marshal(org)
	return string(bytes)
}

// MarshalJSON will marshal organization to JSON.
func (org *Organization) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":           org.Id,
		"name":         org.Name,
		"createAt":     org.CreateAt,
		"updateAt":     org.UpdateAt,
		"projectCount": len(org.projectMap),
	}

	return json.Marshal(&m)
}

// UnmarshalJSON is not supported.
func (org *Organization) UnmarshalJSON([]byte) error {
	return nil
}

// Project defines projects in workstation.
type Project struct {
	Id       string    `yaml:"id" json:"id"`
	Name     string    `yaml:"name" json:"name"`
	CreateAt time.Time `yaml:"createAt" json:"createAt"`
	UpdateAt time.Time `yaml:"updateAt" json:"updateAt"`
}

func NewProject(name string) *Project {
	id := utils.CreateId()
	if len(name) < 1 {
		name = id
	}

	now := time.Now()

	return &Project{
		Id:       id,
		Name:     name,
		CreateAt: now,
		UpdateAt: now,
	}
}
