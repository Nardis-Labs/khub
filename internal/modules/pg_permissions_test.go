package modules

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/gorm"
)

func (s *PGSuite) TestGetPermission() {
	sdk := PGSDK{db: s.DB}
	pid := uuid.New()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "permissions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "app_tag"}).
			AddRow(pid, "App", "app"))

	resp, err := sdk.GetPermissions()
	s.NoError(err, "unexpected error while fetching permission")

	s.Equal(*resp[0].ID, pid)
	s.Equal(resp[0].Name, "App")
	s.Equal(resp[0].AppTag, "app")
}

func (s *PGSuite) TestGetPermissionsByIDs() {
	sdk := PGSDK{db: s.DB}
	permissionID1 := uuid.New()
	permissionID2 := uuid.New()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "permissions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "app_tag"}).
			AddRow(permissionID1, "API", "api").
			AddRow(permissionID2, "APP", "app"))

	resp, err := sdk.GetPermissionsByIDs([]uuid.UUID{permissionID1, permissionID2})
	s.NoError(err, "unexpected error while fetching permissions with IDs")

	s.Equal(*resp[0].ID, permissionID1)
	s.Equal(*resp[1].ID, permissionID2)
	s.Equal(resp[0].Name, "API")
	s.Equal(resp[1].Name, "APP")
	s.Equal(resp[0].AppTag, "api")
	s.Equal(resp[1].AppTag, "app")
}

func (s *PGSuite) TestUpsertPermissionCreate() {
	sdk := PGSDK{db: s.DB}
	pid := uuid.New()
	permission := types.Permission{
		ID:     &pid,
		Name:   "App",
		AppTag: "app",
	}

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	// Expect the initial select query to see if the record exists yet.
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."deleted_at" IS NULL AND "permissions"."id" = $1 ORDER BY "permissions"."id" LIMIT $2`)).
		WithArgs(permission.ID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Expecting a create query.
	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "permissions" ("id","name","app_tag","created_at","updated_at","deleted_at") 
	VALUES ($1,$2,$3,$4,$5,$6)`)).
		WithArgs(permission.ID, permission.Name, permission.AppTag, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	// Call the UpsertPermission function
	result, err := sdk.UpsertPermission(permission)
	if err != nil {
		s.Errorf(err, "error was not expected while creating permission: %s")
	}

	// Ensure the permission result is not nil
	s.NotNil(result, "permission result should not be nil")
	s.Equal(*result.ID, *permission.ID)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}

func (s *PGSuite) TestUpsertPermissionUpdate() {
	sdk := PGSDK{db: s.DB}
	pid := uuid.New()
	permission := types.Permission{
		ID:     &pid,
		Name:   "App",
		AppTag: "app",
	}

	permissionNew := types.Permission{
		ID:     &pid,
		Name:   "AppUpdated",
		AppTag: "app",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "app_tag", "created_at", "updated_at", "deleted_at"}).
		AddRow(permission.ID, permission.Name, permission.AppTag, time.Now(), time.Now(), sql.NullTime{})

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	// Expect the initial select query to see if the record exists yet.
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."deleted_at" IS NULL AND "permissions"."id" = $1 ORDER BY "permissions"."id" LIMIT $2`)).
		WithArgs(permission.ID, 1).
		WillReturnRows(rows)

	// Expecting a create query.
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "permissions" SET "name"=$1,"app_tag"=$2,"created_at"=$3,"updated_at"=$4,"deleted_at"=$5 WHERE "permissions"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs(permissionNew.Name, permissionNew.AppTag, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), permission.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	// Call the UpsertPermission function
	result, err := sdk.UpsertPermission(permissionNew)
	if err != nil {
		s.Errorf(err, "error was not expected while creating permission: %s")
	}

	// Ensure the permission result is not nil
	s.NotNil(result, "permission result should not be nil")
	s.Equal(result.Name, permissionNew.Name)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}
