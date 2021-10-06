package datastore

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
	query := regexp.QuoteMeta("SELECT * FROM `organizations`")

	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	// With empty org result
	sql.sqlMock.ExpectQuery(query)
	assert.Empty(t, sql.ListOrg())

	// With orgs
	sql.sqlMock.ExpectQuery(query).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, time.Now(), time.Now(), nil, "ut-org"))

	assert.Len(t, sql.ListOrg(), 1)
}

func TestMySql_InsertOrg(t *testing.T) {
	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	// Happy case
	now := time.Now()
	org := &Organization{
		Name: "ut-org",
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	query := regexp.QuoteMeta(
		"INSERT INTO `organizations` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES (?,?,?,?)")

	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.InsertOrg(org))

	// With nil org
	assert.False(t, sql.InsertOrg(nil))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()

	assert.False(t, sql.InsertOrg(org))
}

func TestMySql_GetOrg(t *testing.T) {
	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	now := time.Now()
	org := &Organization{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	query := regexp.QuoteMeta(
		"SELECT * FROM `organizations` WHERE id = ? AND `organizations`.`deleted_at` IS NULL")

	// Happy case
	sql.sqlMock.ExpectQuery(query).
		WithArgs(org.Id).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(org.Id, org.CreatedAt, org.UpdatedAt, nil, org.Name))
	assert.NotNil(t, sql.GetOrg(org.Id))

	// With error
	sql.sqlMock.ExpectQuery(query).
		WithArgs(org.Id).
		WillReturnError(errors.New("ut-error"))
	assert.Nil(t, sql.GetOrg(org.Id))
}

func TestMySql_RemoveOrg(t *testing.T) {
	now := time.Now()

	f := func() time.Time {
		return now
	}

	sql := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	sql.Bootstrap(context.TODO())

	query := regexp.QuoteMeta(
		"UPDATE `organizations` SET `deleted_at`=? WHERE `organizations`.`id` = ? AND `organizations`.`deleted_at` IS NULL")

	org := &Organization{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Happy case
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, org.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.RemoveOrg(org.Id))

	// Without result
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, org.Id).
		WillReturnResult(sqlmock.NewResult(1, 0))
	sql.sqlMock.ExpectCommit()
	assert.False(t, sql.RemoveOrg(org.Id))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, org.Id).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()
	assert.False(t, sql.RemoveOrg(org.Id))
}

func TestMySql_UpdateOrg(t *testing.T) {
	now := time.Now()

	f := func() time.Time {
		return now
	}

	sql := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	sql.Bootstrap(context.TODO())

	query := regexp.QuoteMeta(
		"UPDATE `organizations` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`name`=? WHERE `id` = ?")

	// Happy case
	org := &Organization{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name, org.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.UpdateOrg(org))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(org.CreatedAt, org.UpdatedAt, nil, org.Name, org.Id).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()
	assert.False(t, sql.UpdateOrg(org))
}

func TestMySql_ListProject(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `projects` WHERE (org_id = ?) AND `projects`.`deleted_at` IS NULL")

	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	// With error
	sql.sqlMock.ExpectQuery(query).WillReturnError(errors.New("ut-error"))
	assert.Empty(t, sql.ListProject(1))

	// With projects
	sql.sqlMock.ExpectQuery(query).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "org_id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, 1, time.Now(), time.Now(), nil, "ut-org"))

	assert.Len(t, sql.ListProject(1), 1)
}

func TestMySql_InsertProject(t *testing.T) {
	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	now := time.Now()
	proj := &Project{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	query := regexp.QuoteMeta(
		"INSERT INTO `projects` (`created_at`,`updated_at`,`deleted_at`,`org_id`,`name`) VALUES (?,?,?,?,?)")

	// Happy case
	// add organization
	org := &Organization{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Happy case
	sql.sqlMock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `organizations` WHERE id = ? AND `organizations`.`deleted_at` IS NULL")).
		WithArgs(org.Id).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(org.Id, org.CreatedAt, org.UpdatedAt, nil, org.Name))

	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.InsertProject(proj))

	// With nil proj
	assert.False(t, sql.InsertProject(nil))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()
	assert.False(t, sql.InsertProject(proj))
}

func TestMySql_GetProject(t *testing.T) {
	sql := RegisterMySql(WithEnableMockDb())
	sql.Bootstrap(context.TODO())

	now := time.Now()
	proj := &Project{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	query := regexp.QuoteMeta(
		"SELECT * FROM `projects` WHERE (org_id = ? AND id = ?) AND `projects`.`deleted_at` IS NULL")

	// Happy case
	sql.sqlMock.ExpectQuery(query).
		WithArgs(proj.OrgId, proj.Id).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "org_id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(proj.Id, proj.OrgId, proj.CreatedAt, proj.UpdatedAt, nil, proj.Name))
	assert.NotNil(t, sql.GetProject(proj.OrgId, proj.Id))

	// With error
	sql.sqlMock.ExpectQuery(query).
		WithArgs(proj.OrgId, proj.Id).
		WillReturnError(errors.New("ut-error"))
	assert.Nil(t, sql.GetProject(proj.OrgId, proj.Id))
}

func TestMySql_RemoveProject(t *testing.T) {
	now := time.Now()

	f := func() time.Time {
		return now
	}

	sql := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	sql.Bootstrap(context.TODO())

	query := regexp.QuoteMeta(
		"UPDATE `projects` SET `deleted_at`=? WHERE (org_id = ?) AND `projects`.`id` = ? AND `projects`.`deleted_at` IS NULL")

	proj := &Project{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Happy case
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, proj.OrgId, proj.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.RemoveProject(proj.OrgId, proj.Id))

	// Without result
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, proj.OrgId, proj.Id).
		WillReturnResult(sqlmock.NewResult(1, 0))
	sql.sqlMock.ExpectCommit()
	assert.False(t, sql.RemoveProject(proj.OrgId, proj.Id))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(now, proj.OrgId, proj.Id).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()
	assert.False(t, sql.RemoveProject(proj.OrgId, proj.Id))
}

func TestMySql_UpdateProject(t *testing.T) {
	now := time.Now()

	f := func() time.Time {
		return now
	}

	sql := RegisterMySql(
		WithEnableMockDb(),
		WithNowFunc(f))
	sql.Bootstrap(context.TODO())

	query := regexp.QuoteMeta(
		"UPDATE `projects` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`org_id`=?,`name`=? WHERE `id` = ?")

	// Happy case
	proj := &Project{
		Name:  "ut-proj",
		OrgId: 1,
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// add organization
	org := &Organization{
		Name: "ut-org",
		Base: Base{
			Id:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Happy case
	sql.sqlMock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `organizations` WHERE id = ? AND `organizations`.`deleted_at` IS NULL")).
		WithArgs(org.Id).
		WillReturnRows(sql.sqlMock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(org.Id, org.CreatedAt, org.UpdatedAt, nil, org.Name))

	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name, proj.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.sqlMock.ExpectCommit()
	assert.True(t, sql.UpdateProject(proj))

	// With error
	sql.sqlMock.ExpectBegin()
	sql.sqlMock.ExpectExec(query).
		WithArgs(proj.CreatedAt, proj.UpdatedAt, nil, proj.OrgId, proj.Name, proj.Id).
		WillReturnError(errors.New("ut-error"))
	sql.sqlMock.ExpectRollback()
	assert.False(t, sql.UpdateProject(proj))
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
