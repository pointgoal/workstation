// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// LocalFsRootDirDefault default is current directory
	LocalFsRootDirDefault = "."
	// LocalFsMetaFileName as name described
	LocalFsMetaFileName = ".meta"
)

// RegisterLocalFs will register Entry into GlobalAppCtx
func RegisterLocalFs(opts ...LocalFsOption) *LocalFs {
	res := &LocalFs{
		EntryName:        EntryNameDefault,
		EntryType:        "datastore-local-fs",
		EntryDescription: "Local file system",
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		RootDir:          LocalFsRootDirDefault,
		MetaFileName:     LocalFsMetaFileName,
		lastIndex:        make(map[interface{}]int, 0),
	}

	for i := range opts {
		opts[i](res)
	}

	if !path.IsAbs(res.RootDir) {
		wd, _ := os.Getwd()
		res.RootDir = path.Join(wd, res.RootDir)
	}

	rkentry.GlobalAppCtx.AddEntry(res)

	return res
}

// LocalFsOption will be extended in future.
type LocalFsOption func(*LocalFs)

// WithRootPathLocalFs provides root directory of data store
func WithRootPathLocalFs(rootDir string) LocalFsOption {
	return func(fs *LocalFs) {
		fs.RootDir = rootDir
	}
}

// LocalFs implements interface of DataStore whose underlying storage is local file system
type LocalFs struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	RootDir          string                    `json:"rootDir" yaml:"rootDir"`
	MetaFileName     string                    `json:"metaFileName" yaml:"metaFileName"`
	lastIndex        map[interface{}]int       `json:"-" yaml:"-"`
}

// Connect to to remote/local provider
func (l *LocalFs) Connect() error {
	_, err := os.Stat(l.RootDir)

	if err != nil {
		l.ZapLoggerEntry.GetLogger().Warn("Failed to connect to local file system", zap.Error(err))
		return err
	}

	return nil
}

// IsHealthy checks healthy status remote provider
func (l *LocalFs) IsHealthy() bool {
	if err := l.Connect(); err != nil {
		return false
	}

	return true
}

// Read meta file and unmarshal to target interface
func (l *LocalFs) readMetaFile(metaFilePath string, target interface{}) error {
	bytes, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		l.ZapLoggerEntry.GetLogger().Warn(
			fmt.Sprintf("Failed to read meta file from %s", metaFilePath),
			zap.Error(err))
		return err
	}

	if err := json.Unmarshal(bytes, target); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn(
			fmt.Sprintf("Failed to unmarshal from meta file at %s", metaFilePath),
			zap.Error(err))
	}

	return nil
}

// Unmarshal to json and write to meta file
func (l *LocalFs) writeToMetaFile(metaFilePath string, source interface{}) error {
	var bytes []byte
	var err error

	// Marshal to json
	if bytes, err = json.Marshal(source); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn("Failed to marshal meta")
		return err
	}

	// Write to file system
	if err := ioutil.WriteFile(metaFilePath, bytes, os.ModePerm); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn(fmt.Sprintf("Failed to write to meta file at %s", metaFilePath))
		return err
	}

	return nil
}

// Bootstrap will bootstrap datastore
func (l *LocalFs) Bootstrap(ctx context.Context) {
	event := l.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(l.EntryName),
		rkquery.WithEntryType(l.EntryType))
	logger := l.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// Check healthy status of local file system
	if !l.IsHealthy() {
		rkcommon.ShutdownWithError(errors.New("dataStore is not healthy, shutting down"))
	}

	// List organizations, projects and load the meta into maps
	l.lastIndex[organizationKey] = l.maxOrgId()
	l.lastIndex[projectKey] = l.maxProjId()

	l.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping DataStore.", event.ListPayloads()...)
}

// Interrupt will interrupt datastore
func (l *LocalFs) Interrupt(ctx context.Context) {
	event := l.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(l.EntryName),
		rkquery.WithEntryType(l.EntryType))
	logger := l.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	l.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting DataStore.", event.ListPayloads()...)
}

// GetName returns datastore entry name
func (l *LocalFs) GetName() string {
	return l.EntryName
}

// GetType returns datastore entry type
func (l *LocalFs) GetType() string {
	return l.EntryType
}

// GetDescription returns datastore entry description
func (l *LocalFs) GetDescription() string {
	return l.EntryDescription
}

// String returns datastore as string
func (l *LocalFs) String() string {
	bytes, err := json.Marshal(l)
	if err != nil || len(bytes) < 1 {
		return "{}"
	}

	return string(bytes)
}

// ************************************************** //
// ************** Organization related ************** //
// ************************************************** //

// ListOrg as function name described
func (l *LocalFs) ListOrg() []*Organization {
	res := make([]*Organization, 0)

	// List folders
	var fsInfos []fs.FileInfo
	var err error

	if fsInfos, err = ioutil.ReadDir(path.Join(l.RootDir)); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn("Failed to list organizations", zap.Error(err))
		return res
	}

	for i := range fsInfos {
		if !fsInfos[i].IsDir() {
			continue
		}

		// Unmarshal organization meta
		org := &Organization{}
		if err := l.readMetaFile(path.Join(l.RootDir, fsInfos[i].Name(), l.MetaFileName), org); err != nil {
			continue
		}

		res = append(res, org)
	}

	return res
}

// InsertOrg as function name described
func (l *LocalFs) InsertOrg(org *Organization) bool {
	if org == nil {
		return false
	}

	l.assignRequiredFields(org)

	// 1: Create directory named with organization Id
	orgDir := path.Join(l.RootDir, strconv.Itoa(org.Id))
	if err := os.Mkdir(orgDir, os.ModePerm); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn(fmt.Sprintf("Failed to create organization folder at %s", orgDir))
		return false
	}

	// 2: Write organization meta file
	if err := l.writeToMetaFile(path.Join(l.RootDir, strconv.Itoa(org.Id), l.MetaFileName), org); err != nil {
		return false
	}

	return true
}

// GetOrg as function name described
func (l *LocalFs) GetOrg(orgId int) *Organization {
	// Read organization meta file and
	org := &Organization{}
	if err := l.readMetaFile(path.Join(l.RootDir, strconv.Itoa(orgId), l.MetaFileName), org); err != nil {
		return nil
	}

	return org
}

// RemoveOrg as function name described
func (l *LocalFs) RemoveOrg(orgId int) bool {
	if v := l.GetOrg(orgId); v == nil {
		return false
	}

	// 1: Remove organization folder
	if err := os.RemoveAll(path.Join(l.RootDir, strconv.Itoa(orgId))); err != nil {
		return false
	}

	return true
}

// UpdateOrg as function name described
func (l *LocalFs) UpdateOrg(org *Organization) bool {
	if org == nil {
		return false
	}

	if v := l.GetOrg(org.Id); v == nil {
		return false
	}

	org.UpdatedAt = time.Now()
	// 1: Read organization meta file
	if err := l.writeToMetaFile(path.Join(l.RootDir, strconv.Itoa(org.Id), l.MetaFileName), org); err != nil {
		return false
	}

	return true
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProject as function name described
func (l *LocalFs) ListProject(orgId int) []*Project {
	res := make([]*Project, 0)

	// List folders
	var fsInfos []fs.FileInfo
	var err error

	if fsInfos, err = ioutil.ReadDir(path.Join(l.RootDir, strconv.Itoa(orgId))); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn("Failed to list projects", zap.Error(err))
		return res
	}

	for i := range fsInfos {
		if !fsInfos[i].IsDir() {
			continue
		}

		// Unmarshal project meta
		proj := &Project{}
		if err := l.readMetaFile(path.Join(l.RootDir, strconv.Itoa(orgId), fsInfos[i].Name(), l.MetaFileName), proj); err != nil {
			continue
		}

		res = append(res, proj)
	}

	return res
}

// InsertProject as function name described
func (l *LocalFs) InsertProject(proj *Project) bool {
	if proj == nil {
		return false
	}

	l.assignRequiredFields(proj)

	// 1: Create directory named with project Id
	projDir := path.Join(l.RootDir, strconv.Itoa(proj.OrgId), strconv.Itoa(proj.Id))
	if err := os.Mkdir(projDir, os.ModePerm); err != nil {
		l.ZapLoggerEntry.GetLogger().Warn(fmt.Sprintf("Failed to create project folder at %s", projDir))
		return false
	}

	// 2: Write project meta file
	if err := l.writeToMetaFile(path.Join(l.RootDir, strconv.Itoa(proj.OrgId), strconv.Itoa(proj.Id), l.MetaFileName), proj); err != nil {
		return false
	}

	return true
}

// GetProject as function name described
func (l *LocalFs) GetProject(orgId, projId int) *Project {
	// Read organization meta file and
	proj := &Project{}
	if err := l.readMetaFile(path.Join(l.RootDir, strconv.Itoa(orgId), strconv.Itoa(projId), l.MetaFileName), proj); err != nil {
		return nil
	}

	return proj
}

// RemoveProject as function name described
func (l *LocalFs) RemoveProject(orgId, projId int) bool {
	// 1: Remove project folder
	if err := os.RemoveAll(path.Join(l.RootDir, strconv.Itoa(orgId), strconv.Itoa(projId))); err != nil {
		return false
	}

	return true
}

// UpdateProject as function name described
func (l *LocalFs) UpdateProject(proj *Project) bool {
	if proj == nil {
		return false
	}

	proj.UpdatedAt = time.Now()
	// 1: Read organization meta file
	if err := l.writeToMetaFile(path.Join(l.RootDir, strconv.Itoa(proj.OrgId), strconv.Itoa(proj.Id), l.MetaFileName), proj); err != nil {
		return false
	}

	return true
}

// Get max ID of Organization
func (l *LocalFs) maxOrgId() int {
	orgList := l.ListOrg()

	var res int

	for i := range orgList {
		if res < orgList[i].Id {
			res = orgList[i].Id
		}
	}

	return res
}

// Get max ID of Project
func (l *LocalFs) maxProjId() int {
	var res int

	orgList := l.ListOrg()

	for i := range orgList {
		org := orgList[i]
		projList := l.ListProject(org.Id)
		for j := range projList {
			if res < projList[j].Id {
				res = projList[j].Id
			}
		}
	}

	return res
}

// Assign required fields
func (l *LocalFs) assignRequiredFields(in interface{}) {
	switch v := in.(type) {
	case *Organization:
		id := l.lastIndex[organizationKey] + 1
		l.lastIndex[organizationKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	case *Project:
		id := l.lastIndex[projectKey] + 1
		l.lastIndex[projectKey] = id
		v.Id = id
		now := time.Now()
		v.CreatedAt = now
		v.UpdatedAt = now
	}
}
