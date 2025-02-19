package modules

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
)

func (s *PGSuite) TestGetGroups() {
	sdk := PGSDK{db: s.DB}
	groupID := uuid.New()
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_users" WHERE "group_users"."group_id" = $1`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "user_id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_permissions" WHERE "group_permissions"."group_id" = $1`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "permission_id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE "groups"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(groupID, "Admin"))

	resp, err := sdk.GetGroups()
	s.NoError(err, "unexpected error while fetching groups")

	s.Equal(*resp[0].ID, groupID)
	s.Equal(resp[0].Name, "Admin")
}

func (s *PGSuite) TestGetGroupsByIDs() {
	sdk := PGSDK{db: s.DB}
	groupID1 := uuid.New()
	groupID2 := uuid.New()
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_users" WHERE "group_users"."group_id" IN ($1,$2)`)).
		WithArgs(groupID1, groupID2).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "user_id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_permissions" WHERE "group_permissions"."group_id" IN ($1,$2)`)).
		WithArgs(groupID1, groupID2).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "permission_id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE "groups"."id" IN ($1,$2) AND "groups"."deleted_at" IS NULL`)).
		WithArgs(groupID1, groupID2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(groupID1, "Admin").
			AddRow(groupID2, "Mondo"))

	resp, err := sdk.GetGroupsByIDs([]uuid.UUID{groupID1, groupID2})
	s.NoError(err, "unexpected error while fetching groups")

	s.Equal(*resp[0].ID, groupID1)
	s.Equal(resp[0].Name, "Admin")
	s.Equal(*resp[1].ID, groupID2)
	s.Equal(resp[1].Name, "Mondo")
}

func (s *PGSuite) TestUpsertGroupCreate() {
	sdk := PGSDK{db: s.DB}
	gid := uuid.New()
	group := types.Group{
		Name: "Admin",
	}

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE name = $1 AND "groups"."deleted_at" IS NULL ORDER BY "groups"."id" LIMIT $2`)).
		WithArgs(group.Name, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Expecting a create query.
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "groups" ("name","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(group.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
				AddRow(&gid, group.Name, time.Now(), time.Now(), sql.NullTime{}))

	s.mock.ExpectCommit()

	// Call the UpsertGroup function
	result, err := sdk.UpsertGroup(group)
	if err != nil {
		s.NoErrorf(err, "error was not expected while creating group: %s")
	}

	// Ensure the group result is not nil
	s.NotNil(result, "group result should not be nil")
	s.Equal(*result.ID, gid)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}

func (s *PGSuite) TestUpsertGroupUpdate() {
	sdk := PGSDK{db: s.DB}
	gid := uuid.New()
	pid1 := uuid.New()
	pid2 := uuid.New()
	groupPermission1 := types.Permission{
		ID: &pid1,
	}

	groupPermission2 := types.Permission{
		ID: &pid2,
	}

	group := types.Group{
		ID:   &gid,
		Name: "Admin",
	}

	groupUpdated := types.Group{
		ID:   group.ID,
		Name: "Admin",
		Permissions: []*types.Permission{
			&groupPermission1,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(group.ID, group.Name, time.Now(), time.Now(), sql.NullTime{})

	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "groups" WHERE name = $1 AND "groups"."deleted_at" IS NULL ORDER BY "groups"."id" LIMIT $2`)).
		WithArgs(group.Name, 1).
		WillReturnRows(rows)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_permissions" WHERE "group_permissions"."group_id" = $1`)).
		WithArgs(gid).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "permission_id"}).
			AddRow(gid, groupPermission1.ID).
			AddRow(gid, groupPermission2.ID))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "group_users" WHERE "group_users"."group_id" = $1`)).
		WithArgs(gid).
		WillReturnRows(sqlmock.NewRows([]string{"group_id", "user_id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "permissions" WHERE "permissions"."id" IN ($1,$2) AND "permissions`)).
		WithArgs(groupPermission1.ID, groupPermission2.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupPermission1.ID).AddRow(groupPermission2.ID))

	s.mock.ExpectCommit()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "permissions" ("id","name","app_tag","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`)).
		WithArgs(groupPermission1.ID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "permissions" ("id","name","app_tag","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`)).
		WithArgs(groupPermission2.ID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "group_permissions" WHERE group_id = $1 AND permission_id = $2`)).
		WithArgs(gid, groupPermission2.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Expecting a update query.
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "groups" SET "name"=$1,"created_at"=$2,"updated_at"=$3,"deleted_at"=$4 WHERE "groups"."deleted_at" IS NULL AND "id" = $5`)).
		WithArgs(groupUpdated.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), gid).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "group_permissions" ("permission_id","group_id") VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING "group_id"`)).
		WithArgs(groupPermission1.ID, gid).
		WillReturnRows(sqlmock.NewRows([]string{"group_id"}).AddRow(gid))

	s.mock.ExpectCommit()

	// Call the UpsertGroup function
	result, err := sdk.UpsertGroup(groupUpdated)
	if err != nil {
		s.NoErrorf(err, "error was not expected while updating group: %s", err.Error())
	}

	// Ensure the group result is not nil
	s.NotNil(result, "group result should not be nil")
	s.Equal(result.Name, groupUpdated.Name)

	// Assert that the expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s")
	}

}
