// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
)

const (
	// EntryNameDefault default entry name
	EntryNameDefault = "datastore"
)

func GetRepository() Repository {
	raw := rkentry.GlobalAppCtx.GetEntry(EntryNameDefault)
	if raw == nil {
		return nil
	}

	res, _ := raw.(Repository)
	return res
}

// BootConfig is a struct which is for unmarshalled YAML
type BootConfig struct {
	Repository struct {
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

// RegisterRepositoryFromConfig is an implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterRepositoryFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.Repository.Enabled {
		switch config.Repository.Provider {
		case "localFs":
			repo := RegisterLocalFs(
				WithRootPathLocalFs(config.Repository.RootPath))
			res[repo.GetName()] = repo
		case "mySql":
			repo := RegisterMySql(
				WithUser(config.Repository.MySql.User),
				WithPass(config.Repository.MySql.Pass),
				WithProtocol(config.Repository.MySql.Protocol),
				WithAddr(config.Repository.MySql.Addr),
				WithDatabase(config.Repository.MySql.Database),
				WithParams(config.Repository.MySql.Params))
			res[repo.GetName()] = repo
		default:
			repo := RegisterMemory()
			res[repo.GetName()] = repo
		}
	}

	return res
}

type Repository interface {
	rkentry.Entry

	// Connect to to remote/local provider
	Connect() error

	// IsHealthy checks healthy status remote provider
	IsHealthy() bool

	// ************************************************** //
	// ************** Organization related ************** //
	// ************************************************** //

	// ListOrg as function name described
	ListOrg() ([]*Org, error)

	// CreateOrg as function name described
	CreateOrg(org *Org) (bool, error)

	// GetOrg as function name described
	GetOrg(int) (*Org, error)

	// RemoveOrg as function name described
	RemoveOrg(int) (bool, error)

	// UpdateOrg as function name described
	UpdateOrg(org *Org) (bool, error)

	// ********************************************* //
	// ************** Project related ************** //
	// ********************************************* //

	// ListProj as function name described
	ListProj(int) ([]*Proj, error)

	// CreateProj as function name described
	CreateProj(proj *Proj) (bool, error)

	// GetProj as function name described
	GetProj(int, int) (*Proj, error)

	// RemoveProj as function name described
	RemoveProj(int, int) (bool, error)

	// UpdateProj as function name described
	UpdateProj(org *Proj) (bool, error)
}
