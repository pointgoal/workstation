// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
)

const (
	// EntryNameDefault default entry name
	EntryNameDefault = "datastore"
)

func GetDataStore() DataStore {
	raw := rkentry.GlobalAppCtx.GetEntry(EntryNameDefault)
	if raw == nil {
		return nil
	}

	res, _ := raw.(DataStore)
	return res
}

// BootConfig is a struct which is for unmarshalled YAML
type BootConfig struct {
	DataStore struct {
		Enabled  bool   `yaml:"enabled" json:"enabled"`
		Provider string `yaml:"provider" json:"provider"`
		RootPath string `yaml:"rootPath" json:"rootPath"`
		Local    struct {
			DataDir string `yaml:"dataDir" json:"dataDir"`
		} `yaml:"local" json:"local"`
		MySql struct {
			User     string   `yaml:"user" json:"user"`
			Pass     string   `yaml:"pass" json:"pass"`
			Protocol string   `yaml:"protocol" json:"protocol"`
			Addr     string   `yaml:"addr" json:"addr"`
			Database string   `yaml:"database" json:"database"`
			Params   []string `yaml:"params" json:"params"`
		} `yaml:"mySql" json:"mySql"`
		Logger struct {
			ZapLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"zapLogger" json:"zapLogger"`
			EventLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"eventLogger" json:"eventLogger"`
		} `yaml:"logger" json:"logger"`
	} `yaml:"datastore" json:"datastore"`
}

// RegisterDataStoreFromConfig is an implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterDataStoreFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.DataStore.Enabled {
		switch config.DataStore.Provider {
		case "localFs":
			store := RegisterLocalFs(
				WithRootPathLocalFs(config.DataStore.RootPath))
			res[store.GetName()] = store
		case "mySql":
			store := RegisterMySql(
				WithUser(config.DataStore.MySql.User),
				WithPass(config.DataStore.MySql.Pass),
				WithProtocol(config.DataStore.MySql.Protocol),
				WithAddr(config.DataStore.MySql.Addr),
				WithDatabase(config.DataStore.MySql.Database),
				WithParams(config.DataStore.MySql.Params))
			res[store.GetName()] = store
		}
	}

	return res
}

type DataStore interface {
	rkentry.Entry

	// Connect to to remote/local provider
	Connect() error

	// IsHealthy checks healthy status remote provider
	IsHealthy() bool

	// ************************************************** //
	// ************** Organization related ************** //
	// ************************************************** //

	// ListOrg as function name described
	ListOrg() []*Organization

	// InsertOrg as function name described
	InsertOrg(org *Organization) bool

	// GetOrg as function name described
	GetOrg(int) *Organization

	// RemoveOrg as function name described
	RemoveOrg(int) bool

	// UpdateOrg as function name described
	UpdateOrg(org *Organization) bool

	// ********************************************* //
	// ************** Project related ************** //
	// ********************************************* //

	// ListProject as function name described
	ListProject(int) []*Project

	// InsertProject as function name described
	InsertProject(proj *Project) bool

	// GetProject as function name described
	GetProject(int, int) *Project

	// RemoveProject as function name described
	RemoveProject(int, int) bool

	// UpdateProject as function name described
	UpdateProject(org *Project) bool
}
