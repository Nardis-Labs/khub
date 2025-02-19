package modules

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sullivtr/k8s_platform/internal/types"
)

func (s *PGSuite) TestGetMySQLDBCatalog() {
	sdk := PGSDK{db: s.DB}
	host := "test-host.cloud.com"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "my_sqldb_infos"`)).
		WillReturnRows(sqlmock.NewRows([]string{"host", "shortname", "username", "port"}).
			AddRow(host, "test-host", "khub", 3306))

	resp, err := sdk.GetMySQLCatalog()
	s.NoError(err, "unexpected error while fetching mysql db info")

	s.Equal(resp[0].Host, host)
	s.Equal(resp[0].Shortname, "test-host")
	s.Equal(resp[0].Username, "khub")
	s.Equal(resp[0].Port, 3306)
}

func (s *PGSuite) TestUpsertMySQLDBInfo() {
	sdk := PGSDK{db: s.DB}
	dbInfo := &types.MySQLDBInfo{
		Host:      "test-host.cloud.com",
		Shortname: "test-host",
		Username:  "khub",
		Port:      3306,
		IsPrimary: false,
	}
	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "my_sqldb_infos" SET "shortname"=$1,"username"=$2,"port"=$3,"is_primary"=$4,"created_at"=$5,"updated_at"=$6,"deleted_at"=$7 WHERE "my_sqldb_infos"."deleted_at" IS NULL AND "host" = $8`)).
		WithArgs(dbInfo.Shortname, dbInfo.Username, dbInfo.Port, dbInfo.IsPrimary, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), dbInfo.Host).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	_, err := sdk.UpsertMySQLDBInfo(*dbInfo)
	s.NoError(err, "unexpected error while upserting mysql db info")
}

func (s *PGSuite) TestDeleteMySQLDBInfo() {
	sdk := PGSDK{db: s.DB}
	dbInfo := &types.MySQLDBInfo{
		Host: "test-host.cloud.com",
	}
	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "my_sqldb_infos" WHERE host = $1`)).
		WithArgs(dbInfo.Host).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := sdk.DeleteMySQLDBInfo(dbInfo.Host)
	s.NoError(err, "unexpected error while deleting mysql db info")
}
