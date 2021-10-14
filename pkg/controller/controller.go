// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"encoding/json"
	"github.com/pointgoal/workstation/pkg/repository"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	rkquery "github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
)

const (
	// EntryName name of entry
	EntryName = "ws-controller"
	// EntryType type of entry
	EntryType = "ws-controller"
	// EntryDescription description of entry
	EntryDescription = "Entry for controller entry."
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
		Repo:             repository.GetRepository(),
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
	Repo             repository.Repository     `json:"repository" yaml:"repository"`
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
	con.Repo = repository.GetRepository()

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

// GetController returns ProjectEntry.
func GetController() *Controller {
	if raw := rkentry.GlobalAppCtx.GetEntry(EntryName); raw != nil {
		return raw.(*Controller)
	}

	return nil
}
