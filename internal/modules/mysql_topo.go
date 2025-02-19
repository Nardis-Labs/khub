package modules

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

const (
	SHOW_SOURCE_QUERY         = "select HOST from performance_schema.replication_connection_configuration"
	REPLICATION_RUNNING_QUERY = "SELECT SERVICE_STATE FROM performance_schema.replication_connection_status"
	SHOW_USERS_QUERY          = "select DISTINCT USER from information_schema.processlist where USER not in ('system user', 'root', 'dbadmin', 'rdsadmin', 'event_scheduler')"
	SHOW_DMS_REPLICAS         = "select distinct USER from information_schema.processlist where COMMAND='Binlog Dump' AND USER not in ('repl')"
)

// MySQLTopo represents a MySQL database topology module SDK
type MySQLTopoSDK struct {
	MySQLDBPassword string
	Databases       []*types.MySQLDBInfo
}

// CaptureReplicationTopo captures the replication topology for the databases in the MySQLTopoSDK.
// It performs the following steps:
// 1. Initializes a map to store the short names and corresponding host:port pairs of the databases.
// 2. Iterates over the databases to populate the map.
// 3. For each database, it attempts to connect using the appropriate credentials.
// 4. Checks for the replication source and updates the source if found.
// 5. Checks if replication is running and updates the status accordingly.
// 6. Checks for DMS connected users as replicas and adds them to the database's replica list.
func (sdk *MySQLTopoSDK) CaptureReplicationTopo() {
	dbShortNameMap := make(map[string]string)
	for _, db := range sdk.Databases {
		dbShortNameMap[db.Shortname] = strings.Split(db.Host, ".")[0] + ":" + fmt.Sprintf("%d", db.Port)
	}

	for _, db := range sdk.Databases {
		// TODO: Find a nicer way to get the password for differing DBs
		log.Info().Msgf("checking replication topology for %s", db.Shortname)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", db.Username, sdk.MySQLDBPassword, db.Host, db.Port)
		connection, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Warn().Msgf("Error connecting to database, %s: %v", db.Shortname, err)
			continue
		}
		defer connection.Close()

		connection.SetConnMaxLifetime(time.Minute * 3)
		connection.SetMaxOpenConns(10)
		connection.SetMaxIdleConns(10)

		// Check for replication source if there is one
		var sourceHost sql.NullString
		err = connection.QueryRow(SHOW_SOURCE_QUERY).Scan(&sourceHost)
		if err != nil && err != sql.ErrNoRows {
			log.Warn().Msgf("Error executing SHOW_SOURCE_QUERY on %s: %v", db.Shortname, err)
			continue
		}
		if sourceHost.Valid {
			sourceHostName := strings.Split(sourceHost.String, ".")[0]
			for k := range dbShortNameMap {
				if strings.Contains(sourceHostName, k) {
					log.Info().Msgf("source found for %s: %s", db.Shortname, sourceHostName)
					db.SetSource(k)
				}
			}
		}

		// Check if replication is running
		var replState sql.NullString
		err = connection.QueryRow(REPLICATION_RUNNING_QUERY).Scan(&replState)
		if err != nil && err != sql.ErrNoRows {
			log.Warn().Msgf("Error executing REPLICATION_RUNNING_QUERY on %s: %v", db.Shortname, err)
			continue
		}
		if replState.Valid {
			log.Info().Msgf("replication running for %s: %s", db.Shortname, replState.String)
			if replState.String == "OFF" {
				db.SetReplicationRunning(false)
			} else {
				db.SetReplicationRunning(true)
			}
		}

		// Check for DMS connected users as replicas
		rows, err := connection.Query(SHOW_DMS_REPLICAS)
		if err != nil {
			log.Warn().Msgf("Error executing SHOW_DMS_REPLICAS on %s: %v", db.Shortname, err)
			continue
		}
		defer rows.Close()

		uniqueDMSReplicaNodes := make(map[string]bool)
		for rows.Next() {
			var replicaUser string
			if err := rows.Scan(&replicaUser); err != nil {
				log.Warn().Msgf("Error scanning row: %v", err)
				continue
			}
			if _, ok := uniqueDMSReplicaNodes[replicaUser]; !ok {
				sdk.Databases = append(sdk.Databases, &types.MySQLDBInfo{
					Host:      replicaUser + "-dms",
					Shortname: replicaUser,
					Source:    db.Shortname,
				})
				uniqueDMSReplicaNodes[replicaUser] = true
			}

			db.AddReplica(replicaUser)
		}
		if err := rows.Err(); err != nil {
			log.Warn().Msgf("Error iterating rows: %v", err)
		}
	}
}

// SetDBReplicas sets the replicas for each database in the MySQLTopoSDK.
// It performs the following steps:
//  1. Iterates over each database in the sdk.Databases slice.
//  2. For each database, it iterates over the other databases to find those that have the current database
//     as their source and are hosted on different hosts.
//  3. Adds the short name of the source database to the replica list of the current database.
//  4. Logs the replicas for each database.
//
// This function helps in identifying and setting up the replication relationships between the databases.
func (sdk *MySQLTopoSDK) SetDBReplicas() {
	for _, db := range sdk.Databases {
		for _, sdb := range sdk.Databases {
			if sdb.Host != db.Host && sdb.Source == db.Shortname {
				db.AddReplica(sdb.Shortname)
			}
		}
		log.Info().Msgf("Replicas for %s: %v", db.Shortname, db.Replicas)
	}
}
