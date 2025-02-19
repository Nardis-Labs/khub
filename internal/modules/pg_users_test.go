package modules

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
)

func (s *PGSuite) TestGetUsers() {
	sdk := PGSDK{db: s.DB}
	uid := uuid.New()
	lastUsed := time.Now()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "last_used"}).
			AddRow(uid, "test user", "test.user@gmail.com", lastUsed))

	resp, err := sdk.GetUsers()
	s.NoError(err, "unexpected error while fetching users")

	s.Equal(*resp[0].ID, uid)
	s.Equal(resp[0].Name, "test user")
	s.Equal(resp[0].Email, "test.user@gmail.com")
	s.Equal(resp[0].LastUsed, lastUsed)
}

func (s *PGSuite) TestGetUser() {
	sdk := PGSDK{db: s.DB}
	uid := uuid.New()
	groupID := uuid.New()
	lastUsed := time.Now()

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_users" WHERE "group_users"."user_id" = $1`)).
		WithArgs(uid).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "user_id"}).
			AddRow(groupID, uid))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE name = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("test user", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "last_used"}).
			AddRow(uid, "test user", "test.user@gmail.com", lastUsed))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE "groups"."id" = $1 AND "groups"."deleted_at" IS NULL`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(groupID, "Admin"))

	resp, err := sdk.GetUser("test user")
	s.NoError(err, "unexpected error while fetching accounts")

	s.Equal(*resp.ID, uid)
	s.Equal(resp.Name, "test user")
	s.Equal(resp.Email, "test.user@gmail.com")
	s.Equal(resp.LastUsed, lastUsed)
}

func (s *PGSuite) TestGetUserPermissionAccess() {
	sdk := PGSDK{db: s.DB}
	uid := uuid.New()
	groupID := uuid.New()
	permissionID := uuid.New()
	permissionTag := "grid_batch"

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT group_id FROM group_users WHERE user_id = $1`)).
		WithArgs(uid).
		WillReturnRows(sqlmock.NewRows([]string{"group_id"}).
			AddRow(groupID))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_permissions" WHERE "group_permissions"."group_id" = $1`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "permission_id"}).
			AddRow(groupID, permissionID))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "permissions" WHERE "permissions"."id" = $1 AND "permissions"."deleted_at" IS NULL`)).
		WithArgs(permissionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "app_tag"}).
			AddRow(permissionID, "GridBatch", permissionTag))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE "groups"."id" = $1 AND "groups"."deleted_at" IS NULL`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(groupID, "Admin"))

	resp, err := sdk.GetUserAccessDetails(uid)
	s.NoError(err, "unexpected error while fetching accounts")

	s.Equal(resp.UserID, uid)
	s.Equal(resp.GroupIDs[0], groupID)
	s.Equal(resp.PermissionIDs[0], permissionID)
}

func (s *PGSuite) TestUpsertUserCreate() {
	sdk := PGSDK{db: s.DB}
	uid := uuid.New()
	user := types.User{
		Name:     "test.user",
		Email:    "test.user@gmail.com",
		LastUsed: time.Now(),
		IsAdmin:  false,
	}

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	// Expecting a create query.
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","is_admin","last_used","dark_mode","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(user.Name, user.Email, user.IsAdmin, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "is_admin", "last_used", "dark_mode", "created_at", "updated_at", "deleted_at"}).
				AddRow(&uid, user.Name, user.Email, user.IsAdmin, user.LastUsed, true, time.Now(), time.Now(), sql.NullTime{}))

	s.mock.ExpectCommit()

	// Call the UpsertUser function
	result, err := sdk.UpsertUser(user)
	if err != nil {
		s.Errorf(err, "error was not expected while creating user: %s")
	}

	// Ensure the user result is not nil
	s.NotNil(result, "user result should not be nil")
	s.Equal(*result.ID, uid)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}

func (s *PGSuite) TestUpsertUserUpdate() {
	sdk := PGSDK{db: s.DB}
	uid := uuid.New()
	user := types.User{
		ID:       &uid,
		Name:     "test.user",
		Email:    "test.user@gmail.com",
		LastUsed: time.Now(),
		IsAdmin:  false,
	}

	userUpdated := types.User{
		ID:       &uid,
		Name:     "test.user",
		Email:    "test.user@gmail.com",
		LastUsed: time.Now().Add(48 * time.Hour),
		IsAdmin:  false,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "is_admin", "last_used", "dark_mode", "created_at", "updated_at", "deleted_at"}).
		AddRow(user.ID, user.Name, user.Email, user.IsAdmin, user.LastUsed, true, time.Now(), time.Now(), sql.NullTime{})

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	// Expect the initial select query to see if the record exists yet.
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(user.ID).
		WillReturnRows(rows)

	// Expecting a update query.
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"is_admin"=$3,"last_used"=$4,"dark_mode"=$5,"created_at"=$6,"updated_at"=$7,"deleted_at"=$8 WHERE "users"."deleted_at" IS NULL AND "id" = $9`)).
		WithArgs(userUpdated.Name, userUpdated.Email, userUpdated.IsAdmin, userUpdated.LastUsed, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uid).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	// Call the UpsertUser function
	result, err := sdk.UpsertUser(userUpdated)
	if err != nil {
		s.Errorf(err, "error was not expected while creating user: %s")
	}

	// Ensure the user result is not nil
	s.NotNil(result, "user result should not be nil")
	s.Equal(*result.ID, uid)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}
