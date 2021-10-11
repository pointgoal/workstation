package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestMySql_ListOrg(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `orgs` WHERE `orgs`.`deleted_at` IS NULL")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: with organizations
	repo.sqlMock.ExpectQuery(query).
		WillReturnRows(repo.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, time.Now(), time.Now(), nil, "ut-org"))
	orgList, err := repo.ListOrg()
	assert.Len(t, orgList, 1)
	assert.Nil(t, err)

	// 3: with error
	repo.sqlMock.ExpectQuery(query).
		WillReturnError(errors.New("ut-error"))
	orgList, err = repo.ListOrg()
	assert.Empty(t, orgList)
	assert.NotNil(t, err)
}

func TestMySql_CreateOrg(t *testing.T) {
	query := regexp.QuoteMeta("INSERT INTO `orgs` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES (?,?,?,?)")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: init organization to create
	now := time.Now()
	org := &Org{
		Name: "ut-org",
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 3: happy case
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.CreateOrg(org)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 4: with nil organization
	succ, err = repo.CreateOrg(nil)
	assert.False(t, succ)
	assert.NotNil(t, err)

	// 5: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.CreateOrg(org)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_GetOrg(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `orgs` WHERE id = ? AND `orgs`.`deleted_at` IS NULL")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: happy case
	repo.sqlMock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(repo.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, time.Now(), time.Now(), nil, "ut-org"))
	org, err := repo.GetOrg(1)
	assert.NotNil(t, org)
	assert.Nil(t, err)

	// 3: with error
	repo.sqlMock.ExpectQuery(query).
		WithArgs(org.Id).
		WillReturnError(errors.New("ut-error"))
	org, err = repo.GetOrg(1)
	assert.Nil(t, org)
	assert.NotNil(t, err)
}

func TestMySql_RemoveOrg(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE `orgs` SET `deleted_at`=? WHERE `orgs`.`id` = ? AND `orgs`.`deleted_at` IS NULL")

	// 1: init now function for unit test
	now := time.Now()
	f := func() time.Time {
		return now
	}

	// 2: init repo as MySQL
	repo := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	repo.Bootstrap(context.TODO())

	// 3: without proj in org, expect success
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.RemoveOrg(1)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 4: without result
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	repo.sqlMock.ExpectCommit()
	succ, err = repo.RemoveOrg(1)
	assert.False(t, succ)
	assert.NotNil(t, err)

	// 5: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.RemoveOrg(1)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_UpdateOrg(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE `orgs` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`name`=? WHERE `id` = ?")

	// 1: init now function for unit test
	now := time.Now()
	f := func() time.Time {
		return now
	}

	// 2: init repo as MySQL
	repo := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	repo.Bootstrap(context.TODO())

	// 3: init organization for updating
	org := &Org{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 4: happy case
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name, org.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.UpdateOrg(org)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 5: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name, org.Id).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.UpdateOrg(org)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_ListProj(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `projs` WHERE (org_id = ?) AND `projs`.`deleted_at` IS NULL")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: with error
	repo.sqlMock.ExpectQuery(query).WillReturnError(errors.New("ut-error"))
	projList, err := repo.ListProj(1)
	assert.Empty(t, projList)
	assert.NotNil(t, err)

	// 3: with projects
	repo.sqlMock.ExpectQuery(query).
		WillReturnRows(repo.sqlMock.NewRows([]string{"id", "org_id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, 1, time.Now(), time.Now(), nil, "ut-org"))
	projList, err = repo.ListProj(1)
	assert.NotEmpty(t, projList)
	assert.Nil(t, err)
}

func TestMySql_CreateProj(t *testing.T) {
	query := regexp.QuoteMeta("INSERT INTO `projs` (`created_at`,`updated_at`,`deleted_at`,`org_id`,`name`) VALUES (?,?,?,?,?)")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: expect GetOrg() to be called first
	now := time.Now()

	// 3: init project to create
	proj := &Proj{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 4: happy case
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.CreateProj(proj)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 5: with nil proj
	succ, err = repo.CreateProj(nil)
	assert.False(t, succ)
	assert.NotNil(t, err)

	// 6: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.CreateProj(proj)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_GetProj(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `projs` WHERE id = ? AND `projs`.`deleted_at` IS NULL")

	// 1: init repo as MySQL
	repo := RegisterMySql(WithEnableMockDb())
	repo.Bootstrap(context.TODO())

	// 2: happy case
	repo.sqlMock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(repo.sqlMock.NewRows([]string{"id", "org_id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, 1, time.Now(), time.Now(), nil, "ut-proj"))
	proj, err := repo.GetProj(1)
	assert.NotNil(t, proj)
	assert.Nil(t, err)

	// 3: with error
	repo.sqlMock.ExpectQuery(query).
		WithArgs(proj.OrgId, proj.Id).
		WillReturnError(errors.New("ut-error"))
	proj, err = repo.GetProj(1)
	assert.Nil(t, proj)
	assert.NotNil(t, err)
}

func TestMySql_RemoveProj(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE `projs` SET `deleted_at`=? WHERE `projs`.`id` = ? AND `projs`.`deleted_at` IS NULL")

	// 1: init now function for unit test
	now := time.Now()
	f := func() time.Time {
		return now
	}

	// 2: init repo as MySQL
	repo := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	repo.Bootstrap(context.TODO())

	// 3: happy case
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.RemoveProj(1)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 4: without result
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	repo.sqlMock.ExpectCommit()
	succ, err = repo.RemoveProj(1)
	assert.False(t, succ)
	assert.NotNil(t, err)

	// 5: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(now, 1).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.RemoveProj(1)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_UpdateProj(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE `projs` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`org_id`=?,`name`=? WHERE `id` = ?")

	// 1: init now function for unit test
	now := time.Now()
	f := func() time.Time {
		return now
	}

	// 2: init repo as MySQL
	repo := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	repo.Bootstrap(context.TODO())

	proj := &Proj{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 1: happy case
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name, proj.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo.sqlMock.ExpectCommit()
	succ, err := repo.UpdateProj(proj)
	assert.True(t, succ)
	assert.Nil(t, err)

	// 2: with error
	repo.sqlMock.ExpectBegin()
	repo.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name, proj.Id).
		WillReturnError(errors.New("ut-error"))
	repo.sqlMock.ExpectRollback()
	succ, err = repo.UpdateProj(proj)
	assert.False(t, succ)
	assert.NotNil(t, err)
}

func TestMySql_Options(t *testing.T) {
	sql := RegisterMySql(
		WithUser("ut-user"),
		WithPass("ut-pass"),
		WithProtocol("ut-protocol"),
		WithAddr("ut-addr"),
		WithDatabase("ut-db"),
		WithParams([]string{"ut-param"}))

	assert.Equal(t, "ut-user", sql.user)
	assert.Equal(t, "ut-pass", sql.pass)
	assert.Equal(t, "ut-protocol", sql.protocol)
	assert.Equal(t, "ut-addr", sql.addr)
	assert.Equal(t, "ut-db", sql.database)
	assert.Contains(t, sql.params, "ut-param")
}

func TestMySql_Bootstrap(t *testing.T) {
	defer assertNotPanic(t)

	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())
}

func TestMySql_Interrupt(t *testing.T) {
	defer assertNotPanic(t)

	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())
	sql.Interrupt(context.TODO())
}

func TestMySql_EntryFunc(t *testing.T) {
	defer assertNotPanic(t)

	sql := RegisterMySql(WithEnableMockDb())
	assert.Equal(t, EntryNameDefault, sql.GetName())
	assert.Equal(t, "datastore-mysql", sql.GetType())
	assert.Equal(t, "MySQL datastore", sql.GetDescription())
	assert.NotEmpty(t, sql.String())
}

func assertNotPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, false)
	} else {
		// This should never be called in case of a bug
		assert.True(t, true)
	}
}

func assertPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, true)
	} else {
		// This should never be called in case of a bug
		assert.True(t, false)
	}
}
