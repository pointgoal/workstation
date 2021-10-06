// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"encoding/json"
	"github.com/pointgoal/workstation/pkg/datastore"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	rkquery "github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
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
	Controller struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
		Logger  struct {
			ZapLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"zapLogger" json:"zapLogger"`
			EventLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"eventLogger" json:"eventLogger"`
		} `yaml:"logger" json:"logger"`
	} `yaml:"controller" json:"controller"`
}

// RegisterControllerFromConfig is an implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterControllerFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.Controller.Enabled {
		controller := RegisterController()
		res[controller.GetName()] = controller
	}

	return res
}

// RegisterController will register Entry into GlobalAppCtx
func RegisterController(opts ...ControllerOption) *Controller {
	controller := &Controller{
		EntryName:        EntryName,
		EntryType:        EntryType,
		EntryDescription: EntryDescription,
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		DS:               datastore.GetDataStore(),
	}

	for i := range opts {
		opts[i](controller)
	}

	rkentry.GlobalAppCtx.AddEntry(controller)

	return controller
}

// ControllerOption will be extended in future.
type ControllerOption func(*Controller)

// Controller performs as manager of project and organizations
type Controller struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	DS               datastore.DataStore       `json:"dataStore" yaml:"dataStore"`
}

// Bootstrap entry
func (con *Controller) Bootstrap(context.Context) {
	event := con.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(con.EntryName),
		rkquery.WithEntryType(con.EntryType))

	logger := con.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// Register API
	initApi()

	// Get DB
	con.DS = datastore.GetDataStore()

	con.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping controller.", event.ListPayloads()...)
}

// Interrupt entry
func (con *Controller) Interrupt(context.Context) {
	event := con.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(con.EntryName),
		rkquery.WithEntryType(con.EntryType))
	logger := con.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// TODO: Interrupting anything related.

	con.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting controller.", event.ListPayloads()...)
}

// GetName returns entry name
func (con *Controller) GetName() string {
	return con.EntryName
}

// GetDescription returns entry description
func (con *Controller) GetDescription() string {
	return con.EntryDescription
}

// GetType returns entry type as project
func (con *Controller) GetType() string {
	return con.EntryType
}

// String returns entry as string
func (con *Controller) String() string {
	bytes, _ := json.Marshal(con)
	return string(bytes)
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// ListOrg respond to GET v1/org
func (con *Controller) ListOrg() []*datastore.Organization {
	return con.DS.ListOrg()
}

// AddOrg respond  to PUT v1/org.
func (con *Controller) AddOrg(name string) int {
	org := datastore.NewOrganization(name)
	if !con.DS.InsertOrg(org) {
		return -1
	}

	return org.Id
}

// GetOrg respond to GET v1/org/<org-id>
func (con *Controller) GetOrg(id int) *datastore.Organization {
	return con.DS.GetOrg(id)
}

// DeleteOrg respond to DELETE v1/org/<org-id>
func (con *Controller) DeleteOrg(id int) bool {
	return con.DS.RemoveOrg(id)
}

// UpdateOrg respond to POST v1/org/<org-id>
func (con *Controller) UpdateOrg(org *datastore.Organization) bool {
	return con.DS.UpdateOrg(org)
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProjects respond to GET v1/org/<org-id>/proj
func (con *Controller) ListProject(orgId int) []*datastore.Project {
	return con.DS.ListProject(orgId)
}

// GetProject respond to GET v1/org/<org-id>/proj/<proj-id>
func (con *Controller) GetProject(orgId, projId int) *datastore.Project {
	return con.DS.GetProject(orgId, projId)
}

// AddProject respond to PUT v1/org/<org-id>/proj
func (con *Controller) AddProject(proj *datastore.Project) int {
	if !con.DS.InsertProject(proj) {
		return -1
	}

	return proj.Id
}

// DeleteProject respond to DELETE v1/org/<org-id>/proj/<proj-id>
func (con *Controller) RemoveProject(orgId, projId int) bool {
	return con.DS.RemoveProject(orgId, projId)
}

// UpdateProject respond to POST v1/org/<org-id>/proj/<proj-id>
func (con *Controller) UpdateProject(proj *datastore.Project) bool {
	return con.DS.UpdateProject(proj)
}

// GetController returns ProjectEntry.
func GetController() *Controller {
	if raw := rkentry.GlobalAppCtx.GetEntry(EntryName); raw != nil {
		return raw.(*Controller)
	}

	return nil
}
