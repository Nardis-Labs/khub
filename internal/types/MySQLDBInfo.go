package types

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// MySQLDBInfo represents a MySQL database within the topology
type MySQLDBInfo struct {
	Host               string         `json:"host" gorm:"primaryKey"`
	Shortname          string         `json:"shortName"`
	Username           string         `json:"username"`
	Port               int            `json:"port"`
	IsPrimary          bool           `json:"isPrimary"`
	Replicas           []string       `json:"replicas" gorm:"-"`
	Source             string         `json:"source" gorm:"-"`
	ReplicationRunning bool           `json:"replication_running" gorm:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (d *MySQLDBInfo) IsValid() (bool, string) {
	errors := strings.Builder{}
	if d.Host == "" {
		errors.WriteString(fmt.Sprintln("Host is invalid. Must not be empty."))
	}

	if d.Shortname == "" {
		errors.WriteString(fmt.Sprintln("Shortname is invalid. Must not be empty"))
	}

	if d.Username == "" {
		errors.WriteString(fmt.Sprintln("Username is invalid. Must not be empty"))
	}

	errMsg := errors.String()
	if len(errMsg) > 0 {
		return false, errMsg
	}

	return true, ""
}

// SetSource sets the replication source of the database
func (db *MySQLDBInfo) SetSource(host string) {
	db.Source = host
}

// SetReplicationRunning sets the replication status of the database to running
func (db *MySQLDBInfo) SetReplicationRunning(running bool) {
	db.ReplicationRunning = running
}

// AddReplica adds a replica to the database
func (db *MySQLDBInfo) AddReplica(replica string) {
	db.Replicas = append(db.Replicas, replica)
}

type ReplTopoTreeNode struct {
	ID       string                   `json:"id"`
	Data     MySQLDBInfo              `json:"data"`
	Position ReplTopoTreeNodePosition `json:"position"`
}

type ReplTopoTreeNodePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ReplTopoTreeEdge struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	EdgeType string `json:"edgeType"`
	Animated bool   `json:"animated"`
}

type ReplTopoTree struct {
	Nodes []ReplTopoTreeNode `json:"nodes"`
	Edges []ReplTopoTreeEdge `json:"edges"`
}
