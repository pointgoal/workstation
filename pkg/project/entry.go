// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package project

import (
	"context"
	"encoding/json"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	rkquery "github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"sync"
)

const (
	// EntryName name of entry
	EntryName = "ws-project-name"
	// EntryType type of entry
	EntryType = "ws-project"
	// EntryDescription description of entry
	EntryDescription = "Entry for project management entry."
)

// BootConfig is a struct which is for unmarshalled YAML
type BootConfig struct {
	Project struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
		Logger  struct {
			ZapLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"zapLogger" json:"zapLogger"`
			EventLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"eventLogger" json:"eventLogger"`
		} `yaml:"logger" json:"logger"`
	} `yaml:"project" json:"project"`
}

// RegisterEntriesFromConfig is an implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterEntriesFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.Project.Enabled {
		entry := RegisterEntry()
		res[entry.GetName()] = entry
	}

	return res
}

// RegisterEntry will register Entry into GlobalAppCtx
func RegisterEntry(opts ...EntryOption) *Entry {
	entry := &Entry{
		EntryName:        EntryName,
		EntryType:        EntryType,
		EntryDescription: EntryDescription,
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		orgs:             make(map[string]*Organization, 0),
		lock:             sync.Mutex{},
	}

	for i := range opts {
		opts[i](entry)
	}

	rkentry.GlobalAppCtx.AddEntry(entry)

	return entry
}

// EntryOption will be extended in future.
type EntryOption func(*Entry)

// Entry performs as manager of project and organizations
type Entry struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	orgs             map[string]*Organization  `yaml:"-" json:"-"`
	lock             sync.Mutex                `yaml:"-" json:"-"`
}

// Bootstrap entry
func (entry *Entry) Bootstrap(context.Context) {
	event := entry.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(entry.EntryName),
		rkquery.WithEntryType(entry.EntryType))

	logger := entry.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// Register API
	initApi()

	entry.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping ProjectEntry.", event.ListPayloads()...)
}

// Interrupt entry
func (entry *Entry) Interrupt(context.Context) {
	event := entry.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(entry.EntryName),
		rkquery.WithEntryType(entry.EntryType))
	logger := entry.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// TODO: Interrupting anything related.

	entry.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting ProjectEntry.", event.ListPayloads()...)
}

// GetName returns entry name
func (entry *Entry) GetName() string {
	return entry.EntryName
}

// GetDescription returns entry description
func (entry *Entry) GetDescription() string {
	return entry.EntryDescription
}

// GetType returns entry type as project
func (entry *Entry) GetType() string {
	return entry.EntryType
}

// String returns entry as string
func (entry *Entry) String() string {
	bytes, _ := json.Marshal(entry)
	return string(bytes)
}

// ListOrgs respond to GET v1/org
func (entry *Entry) ListOrgs() []*Organization {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	res := make([]*Organization, 0)

	for _, v := range entry.orgs {
		res = append(res, v)
	}

	return res
}

// AddOrg respond  to PUT v1/org.
func (entry *Entry) AddOrg(name string) string {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := NewOrganization(name)
	entry.orgs[org.Id] = org

	return org.Id
}

// GetOrg respond to GET v1/org/<org-id>
func (entry *Entry) GetOrg(id string) *Organization {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	return entry.orgs[id]
}

// DeleteOrg respond to DELETE v1/org/<org-id>
func (entry *Entry) DeleteOrg(id string) bool {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org, ok := entry.orgs[id]
	if !ok {
		return false
	}

	projects := org.ListProjects()
	for i := range projects {
		org.RemoveProject(projects[i].Id)
	}

	delete(entry.orgs, id)

	return true
}

// UpdateOrg respond to POST v1/org/<org-id>
func (entry *Entry) UpdateOrg(org *Organization) bool {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	if org == nil {
		return false
	}

	if _, ok := entry.orgs[org.Id]; !ok {
		return false
	}

	entry.orgs[org.Id] = org

	return true
}

// ListProjects respond to GET v1/org/<org-id>/proj
func (entry *Entry) ListProjects(orgId string) []*Project {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := entry.orgs[orgId]

	if org == nil {
		return make([]*Project, 0)
	}

	return org.ListProjects()
}

// GetProject respond to GET v1/org/<org-id>/proj/<proj-id>
func (entry *Entry) GetProject(orgId, projId string) *Project {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := entry.orgs[orgId]

	if org == nil {
		return nil
	}

	return org.GetProject(projId)
}

// AddProject respond to PUT v1/org/<org-id>/proj
func (entry *Entry) AddProject(orgId string, proj *Project) string {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := entry.orgs[orgId]

	if org == nil || proj == nil {
		return ""
	}

	org.UpsertProject(proj)

	return proj.Id
}

// DeleteProject respond to DELETE v1/org/<org-id>/proj/<proj-id>
func (entry *Entry) DeleteProject(orgId, projId string) bool {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := entry.orgs[orgId]

	if org == nil {
		return false
	}

	return org.RemoveProject(projId)
}

// UpdateProject respond to POST v1/org/<org-id>/proj/<proj-id>
func (entry *Entry) UpdateProject(orgId string, proj *Project) bool {
	entry.lock.Lock()
	defer entry.lock.Unlock()

	org := entry.orgs[orgId]

	if org == nil {
		return false
	}

	return org.UpsertProject(proj)
}

// GetEntry returns ProjectEntry.
func GetEntry() *Entry {
	if raw := rkentry.GlobalAppCtx.GetEntry(EntryName); raw != nil {
		return raw.(*Entry)
	}

	return nil
}
