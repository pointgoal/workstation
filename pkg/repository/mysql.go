// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

// RegisterEntry will register Entry into GlobalAppCtx
func RegisterMySql(opts ...MySqlOption) *MySql {
	res := &MySql{
		EntryName:        EntryNameDefault,
		EntryType:        "datastore-mysql",
		EntryDescription: "MySQL datastore",
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		user:             "root",
		pass:             "pass",
		protocol:         "tcp",
		addr:             "localhost:3306",
		database:         "workstation",
		params:           make([]string, 0),
	}

	for i := range opts {
		opts[i](res)
	}

	rkentry.GlobalAppCtx.AddEntry(res)

	return res
}

// DataStore will be extended in future.
type MySqlOption func(*MySql)

// WithUser provide user
func WithUser(user string) MySqlOption {
	return func(m *MySql) {
		if len(user) > 0 {
			m.user = user
		}
	}
}

// WithPass provide password
func WithPass(pass string) MySqlOption {
	return func(m *MySql) {
		if len(pass) > 0 {
			m.pass = pass
		}
	}
}

// WithProtocol provide protocol
func WithProtocol(protocol string) MySqlOption {
	return func(m *MySql) {
		if len(protocol) > 0 {
			m.protocol = protocol
		}
	}
}

// WithAddr provide address
func WithAddr(addr string) MySqlOption {
	return func(m *MySql) {
		if len(addr) > 0 {
			m.addr = addr
		}
	}
}

// WithDatabase provide database
func WithDatabase(database string) MySqlOption {
	return func(m *MySql) {
		if len(database) > 0 {
			m.database = database
		}
	}
}

// WithParams provide params
func WithParams(params []string) MySqlOption {
	return func(m *MySql) {
		if len(params) > 0 {
			m.params = append(m.params, params...)
		}
	}
}

// WithEnableMockDb enables mock DB
func WithEnableMockDb() MySqlOption {
	return func(m *MySql) {
		m.enableMockDb = true
	}
}

// WithNowFunc provides now functions for unit test
func WithNowFunc(f func() time.Time) MySqlOption {
	return func(m *MySql) {
		m.nowFunc = f
	}
}

// MySql implements interface of DataStore whose underlying storage is MySQL DB
type MySql struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	user             string                    `yaml:"user" json:"user"`
	pass             string                    `yaml:"pass" json:"pass"`
	protocol         string                    `yaml:"protocol" json:"protocol"`
	addr             string                    `yaml:"addr" json:"addr"`
	database         string                    `yaml:"database" json:"database"`
	params           []string                  `yaml:"params" json:"params"`
	db               *gorm.DB                  `yaml:"-" json:"-"`
	// For unit test
	enableMockDb bool             `yaml:"-" json:"-"`
	sqlMock      sqlmock.Sqlmock  `yaml:"-" json:"-"`
	nowFunc      func() time.Time `yaml:"-" json:"-"`
}

// Create database if missing
func (m *MySql) createDbIfMissing() error {
	// init gorm.DB
	sqlParams := ""
	for i := range m.params {
		sqlParams += m.params[i] + "&"
	}
	sqlParams = strings.TrimSuffix(sqlParams, "&")

	dsn := fmt.Sprintf("%s:%s@%s(%s)/?%s",
		m.user, m.pass, m.protocol, m.addr, sqlParams)

	var db *gorm.DB
	var err error

	if m.enableMockDb {
		// Mock db enabled for unit test
		var sqlDb *sql.DB
		sqlDb, m.sqlMock, _ = sqlmock.New()
		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDb,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})

		// For unit test
		m.sqlMock.ExpectExec(
			fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4;", m.database)).
			WillReturnResult(driver.RowsAffected(0))
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return err
	}

	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4;",
		m.database,
	)

	if err := db.Exec(createSQL).Error; err != nil {
		return err
	}

	return nil
}

// Connect to to remote/local provider
func (m *MySql) Connect() error {
	// init gorm.DB
	sqlParams := ""
	for i := range m.params {
		sqlParams += m.params[i] + "&"
	}
	sqlParams = strings.TrimSuffix(sqlParams, "&")

	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?%s",
		m.user, m.pass, m.protocol, m.addr, m.database, sqlParams)

	var db *gorm.DB
	var err error

	if m.enableMockDb {
		// Mock db enabled for unit test
		var sqlDb *sql.DB
		sqlDb, m.sqlMock, _ = sqlmock.New()
		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn:                      sqlDb,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{
			NowFunc: m.nowFunc,
		})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return err
	}

	m.db = db
	return nil
}

// IsHealthy checks healthy status remote provider
func (m *MySql) IsHealthy() bool {
	if d, err := m.db.DB(); err != nil {
		return false
	} else {
		if err := d.Ping(); err != nil {
			return false
		}
	}

	return true
}

// Bootstrap will bootstrap datastore
func (m *MySql) Bootstrap(ctx context.Context) {
	event := m.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(m.EntryName),
		rkquery.WithEntryType(m.EntryType))
	logger := m.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// Create db if missing
	if err := m.createDbIfMissing(); err != nil {
		m.ZapLoggerEntry.GetLogger().Error("failed to create database", zap.Error(err))
		rkcommon.ShutdownWithError(fmt.Errorf("failed to create database at %s:%s@%s(%s)/%s",
			m.user, "****", m.protocol, m.addr, m.database))
	}

	// Connect to db
	if err := m.Connect(); err != nil {
		m.ZapLoggerEntry.GetLogger().Error("failed to connect database", zap.Error(err))
		rkcommon.ShutdownWithError(fmt.Errorf("failed to open database at %s:%s@%s(%s)/%s",
			m.user, "****", m.protocol, m.addr, m.database))
	}

	m.db.AutoMigrate(&Org{})
	m.db.AutoMigrate(&Proj{})

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping repository.", event.ListPayloads()...)
}

// Interrupt will interrupt datastore
func (m *MySql) Interrupt(ctx context.Context) {
	event := m.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(m.EntryName),
		rkquery.WithEntryType(m.EntryType))
	logger := m.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	m.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting repository.", event.ListPayloads()...)
}

// GetName returns datastore entry name
func (m *MySql) GetName() string {
	return m.EntryName
}

// GetType returns datastore entry type
func (m *MySql) GetType() string {
	return m.EntryType
}

// GetDescription returns datastore entry description
func (m *MySql) GetDescription() string {
	return m.EntryDescription
}

// String returns datastore as string
func (m *MySql) String() string {
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
func (m *MySql) ListOrg() ([]*Org, error) {
	orgList := make([]*Org, 0)
	res := m.db.Find(&orgList)

	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to list organizations from DB", zap.Error(res.Error))
		return orgList, res.Error
	}

	return orgList, nil
}

// CreateOrg as function name described
func (m *MySql) CreateOrg(org *Org) (bool, error) {
	if org == nil {
		return false, errors.New("nil organization")
	}

	res := m.db.Create(org)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to create organizations in DB", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, fmt.Errorf("failed to create organization with name:%s", org.Name)
	}

	return true, nil
}

// GetOrg as function name described
func (m *MySql) GetOrg(orgId int) (*Org, error) {
	org := &Org{}
	res := m.db.Where("id = ?", orgId).Find(org)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to get organizations from DB", zap.Error(res.Error))
		return nil, res.Error
	}

	if res.RowsAffected < 1 {
		return nil, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	return org, nil
}

// RemoveOrg as function name described
func (m *MySql) RemoveOrg(orgId int) (bool, error) {
	res := m.db.Delete(&Org{}, orgId)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to delete organizations from DB", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, NewNotFoundf(OrgNotFoundMsg, orgId)
	}

	return true, nil
}

// UpdateOrg as function name described
func (m *MySql) UpdateOrg(org *Org) (bool, error) {
	if org == nil {
		return false, errors.New("nil organization")
	}

	res := m.db.Save(org)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to update organizations to DB", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, NewNotFoundf(OrgNotFoundMsg, org.Id)
	}

	return true, nil
}

// ********************************************* //
// ************** Project related ************** //
// ********************************************* //

// ListProj as function name described
func (m *MySql) ListProj(orgId int) ([]*Proj, error) {
	projList := make([]*Proj, 0)

	res := m.db.Where("org_id = ?", orgId).Find(&projList)

	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to list projects from DB", zap.Error(res.Error))
		return projList, res.Error
	}
	return projList, nil
}

// CreateProj as function name described
func (m *MySql) CreateProj(proj *Proj) (bool, error) {
	if proj == nil {
		return false, errors.New("nil project")
	}

	res := m.db.Where("org_id = ?", proj.OrgId).Create(proj)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to insert project", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, fmt.Errorf("failed to create project with orgId:%d id:%d", proj.OrgId, proj.Id)
	}

	return true, nil
}

// GetProj as function name described
func (m *MySql) GetProj(orgId, projId int) (*Proj, error) {
	proj := &Proj{}
	res := m.db.Where("org_id = ? AND id = ?", orgId, projId).Find(proj)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to get project from DB", zap.Error(res.Error))
		return nil, res.Error
	}

	if res.RowsAffected < 1 {
		return nil, NewNotFoundf(ProjNotFoundMsg, orgId, projId)
	}

	return proj, nil
}

// RemoveProj as function name described
func (m *MySql) RemoveProj(orgId, projId int) (bool, error) {
	res := m.db.Where("org_id = ?", orgId).Delete(&Proj{}, projId)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to delete project from DB", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, NewNotFoundf(ProjNotFoundMsg, orgId, projId)
	}

	return true, nil
}

// UpdateProj as function name described
func (m *MySql) UpdateProj(proj *Proj) (bool, error) {
	if proj == nil {
		return false, errors.New("nil project")
	}

	res := m.db.Save(proj)
	if res.Error != nil {
		m.ZapLoggerEntry.GetLogger().Warn("failed to update project to DB", zap.Error(res.Error))
		return false, res.Error
	}

	if res.RowsAffected < 1 {
		return false, NewNotFoundf(ProjNotFoundMsg, proj.OrgId, proj.Id)
	}

	return true, nil
}
