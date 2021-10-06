// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"time"
)

// RegisterMemory will register Entry into GlobalAppCtx
func RegisterMemory() *Memory {
	res := &Memory{
		EntryName:        EntryNameDefault,
		EntryType:        "datastore-memory",
		EntryDescription: "In memory",
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		orgMap:           make(map[int]*Organization, 0),
		projectMap:       make(map[int]*Project, 0),
		lastIndex:        make(map[interface{}]int, 0),
	}

	rkentry.GlobalAppCtx.AddEntry(res)

	return res
}

// Memory implements interface of DataStore whose underlying storage is memory
type Memory struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	orgMap           map[int]*Organization     `json:"-" yaml:"-"`
	projectMap       map[int]*Project          `json:"-" yaml:"-"`
	lastIndex        map[interface{}]int       `json:"-" yaml:"-"`
}

// Connect to to remote/local provider
func (m *Memory) Connect() error {
	return nil
}

// IsHealthy checks healthy status remote provider
func (m *Memory) IsHealthy() bool {
	return true
}

// Bootstrap will bootstrap datastore
func (m *Memory) Bootstrap(ctx context.Context) {
	event := m.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(m.EntryName),
		rkquery.WithEntryType(m.EntryType))
	logger := m.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	if !m.IsHealthy() {
		rkcommon.ShutdownWithError(errors.New("dataStore is not healthy, shutting down"))
	}

	// List organizations, projects and load the meta into maps
	m.lastIndex[organizationKey] = m.maxOrgId()
	m.lastIndex[projectKey] = m.maxProjId()

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping DataStore.", event.ListPayloads()...)
}

// Interrupt will interrupt datastore
func (m *Memory) Interrupt(ctx context.Context) {
	event := m.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(m.EntryName),
		rkquery.WithEntryType(m.EntryType))
	logger := m.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting DataStore.", event.ListPayloads()...)
}

// GetName returns datastore entry name
func (m *Memory) GetName() string {
	return m.EntryName
}

// GetType returns datastore entry type
func (m *Memory) GetType() string {
	return m.EntryType
}

// GetDescription returns datastore entry description
func (m *Memory) GetDescription() string {
	return m.EntryDescription
}

// String returns datastore as string
func (m *Memory) String() string {
	bytes, err := json.Marshal(m)
	if err != nil || len(bytes) < 1 {
		return "{}"
	}

	return string(bytes)
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// ListOrg as function name described
func (m *Memory) ListOrg() []*Organization {
	res := make([]*Organization, 0)

	for _, v := range m.orgMap {
		res = append(res, v)
	}

	return res
}

// InsertOrg as function name described
func (m *Memory) InsertOrg(org *Organization) bool {
	if org == nil {
		return false
	}

	m.assignRequiredFields(org)

	m.orgMap[org.Id] = org
	return true
}

// GetOrg as function name described
func (m *Memory) GetOrg(orgId int) *Organization {
	return m.orgMap[orgId]
}

// RemoveOrg as function name described
func (m *Memory) RemoveOrg(orgId int) bool {
	_, contains := m.orgMap[orgId]

	delete(m.orgMap, orgId)

	return contains
}

// UpdateOrg as function name described
func (m *Memory) UpdateOrg(org *Organization) bool {
	if org == nil {
		return false
	}

	if _, ok := m.orgMap[org.Id]; !ok {
		return false
	}

	m.orgMap[org.Id] = org

	return true
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProject as function name described
func (m *Memory) ListProject(orgId int) []*Project {
	res := make([]*Project, 0)

	for _, v := range m.projectMap {
		if v.OrgId == orgId {
			res = append(res, v)
		}
	}

	return res
}

// InsertProject as function name described
func (m *Memory) InsertProject(proj *Project) bool {
	if proj == nil {
		return false
	}

	m.assignRequiredFields(proj)

	org := m.orgMap[proj.OrgId]
	if org == nil {
		return false
	}

	m.projectMap[proj.Id] = proj
	return true
}

// GetProject as function name described
func (m *Memory) GetProject(orgId, projId int) *Project {
	return m.projectMap[projId]
}

// RemoveProject as function name described
func (m *Memory) RemoveProject(orgId int, projId int) bool {
	org := m.orgMap[orgId]
	if org == nil {
		return false
	}

	// Remove from proj list
	delete(m.projectMap, projId)

	return true
}

// UpdateProject as function name described
func (m *Memory) UpdateProject(proj *Project) bool {
	_, ok := m.projectMap[proj.Id]
	if !ok {
		return false
	}

	m.projectMap[proj.Id] = proj

	return true
}

// Get max ID of Organization
func (m *Memory) maxOrgId() int {
	orgList := m.ListOrg()

	var res int

	for i := range orgList {
		if res < orgList[i].Id {
			res = orgList[i].Id
		}
	}

	return res
}

// Get max ID of Project
func (m *Memory) maxProjId() int {
	var res int

	orgList := m.ListOrg()

	for i := range orgList {
		org := orgList[i]
		projList := m.ListProject(org.Id)
		for j := range projList {
			if res < projList[j].Id {
				res = projList[j].Id
			}
		}
	}

	return res
}

// Assign required fields
func (m *Memory) assignRequiredFields(in interface{}) {
	switch v := in.(type) {
	case *Organization:
		id := m.lastIndex[organizationKey] + 1
		m.lastIndex[organizationKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	case *Project:
		id := m.lastIndex[projectKey] + 1
		m.lastIndex[projectKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	}
}
