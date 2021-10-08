// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		orgMap:           make(map[int]*Org, 0),
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
	orgMap           map[int]*Org              `json:"-" yaml:"-"`
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
	m.lastIndex[orgKey] = m.maxOrgId()
	m.lastIndex[projKey] = m.maxProjId()

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping repository.", event.ListPayloads()...)
}

// Interrupt will interrupt datastore
func (m *Memory) Interrupt(ctx context.Context) {
	event := m.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(m.EntryName),
		rkquery.WithEntryType(m.EntryType))
	logger := m.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting repository.", event.ListPayloads()...)
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
func (m *Memory) ListOrg() ([]*Org, error) {
	res := make([]*Org, 0)

	for _, v := range m.orgMap {
		res = append(res, v)
	}

	return res, nil
}

// CreateOrg as function name described
func (m *Memory) CreateOrg(org *Org) (bool, error) {
	if org == nil {
		return false, fmt.Errorf("nil organization")
	}

	m.assignRequiredFields(org)

	m.orgMap[org.Id] = org
	return true, nil
}

// GetOrg as function name described
func (m *Memory) GetOrg(orgId int) (*Org, error) {
	res, _ := m.orgMap[orgId]
	return res, nil
}

// RemoveOrg as function name described
func (m *Memory) RemoveOrg(orgId int) (bool, error) {
	org, contains := m.orgMap[orgId]

	if !contains || org == nil {
		return false, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	delete(m.orgMap, orgId)

	return true, nil
}

// UpdateOrg as function name described
func (m *Memory) UpdateOrg(org *Org) (bool, error) {
	if org == nil {
		return false, errors.New("nil organization")
	}

	old, ok := m.orgMap[org.Id]
	if !ok {
		return false, NewNotFoundf(OrgNotFoundMsg, org.Id)
	}

	old.Name = org.Name
	old.UpdatedAt = time.Now()

	return true, nil
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProj as function name described
func (m *Memory) ListProj(orgId int) ([]*Proj, error) {
	res := make([]*Proj, 0)

	org, ok := m.orgMap[orgId]
	if !ok {
		return res, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	for i := range org.ProjList {
		res = append(res, org.ProjList[i])
	}

	return res, nil
}

// CreateProj as function name described
func (m *Memory) CreateProj(proj *Proj) (bool, error) {
	if proj == nil {
		return false, errors.New("nil project")
	}

	m.assignRequiredFields(proj)

	org, ok := m.orgMap[proj.OrgId]
	if !ok || org == nil {
		return false, NewNotFoundf(OrgNotFoundMsg, proj.OrgId)
	}

	org.ProjList = append(org.ProjList, proj)

	return true, nil
}

// GetProj as function name described
func (m *Memory) GetProj(orgId, projId int) (*Proj, error) {
	org, ok := m.orgMap[orgId]

	if !ok || org == nil {
		return nil, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	var res *Proj
	for i := range org.ProjList {
		proj := org.ProjList[i]
		if proj.Id == projId {
			res = proj
			break
		}
	}

	if res == nil {
		return nil, NewNotFoundf(ProjNotFoundMsg, orgId, projId)
	}

	return res, nil
}

// RemoveProj as function name described
func (m *Memory) RemoveProj(orgId int, projId int) (bool, error) {
	org, ok := m.orgMap[orgId]
	if !ok || org == nil {
		return false, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	// Remove from proj list
	index := -1
	for index = range org.ProjList {
		proj := org.ProjList[index]
		if proj.Id == projId {
			break
		}
	}

	if index < 0 {
		return false, NewNotFoundf(ProjNotFoundMsg, orgId, projId)
	}
	org.ProjList = append(org.ProjList[:index], org.ProjList[index+1:]...)

	return true, nil
}

// UpdateProj as function name described
func (m *Memory) UpdateProj(proj *Proj) (bool, error) {
	if proj == nil {
		return false, fmt.Errorf("nil project")
	}

	org, ok := m.orgMap[proj.OrgId]
	if !ok || org == nil {
		return false, NewNotFoundf(OrgNotFoundMsg, proj.OrgId)
	}

	index := -1
	for index = range org.ProjList {
		e := org.ProjList[index]
		if proj.Id == e.Id {
			break
		}
	}

	if index < 0 {
		return false, NewNotFoundf(ProjNotFoundMsg, proj.OrgId, proj.Id)
	}

	org.ProjList[index].Name = proj.Name
	org.ProjList[index].UpdatedAt = time.Now()

	return true, nil
}

// Get max ID of Organization
func (m *Memory) maxOrgId() int {
	orgList, err := m.ListOrg()

	if err != nil {
		return -1
	}

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

	orgList, err := m.ListOrg()
	if err != nil {
		return -1
	}

	for i := range orgList {
		org := orgList[i]
		projList, err := m.ListProj(org.Id)
		if err != nil {
			continue
		}

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
	case *Org:
		id := m.lastIndex[orgKey] + 1
		m.lastIndex[orgKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	case *Proj:
		id := m.lastIndex[projKey] + 1
		m.lastIndex[projKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	}
}
