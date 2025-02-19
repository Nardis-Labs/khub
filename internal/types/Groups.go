package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group represents a group on khub application.
// Each group has a set of roles that indicate user access within the group.
//
// group_roles is a join table between groups & roles
// group_users is a join table between groups & users
type Group struct {
	ID          *uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex"`
	Permissions []*Permission  `json:"permissions" gorm:"many2many:group_permissions;"`
	Users       []*User        `json:"users" gorm:"many2many:group_users;"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (g *Group) IsValid() (bool, string) {
	errors := strings.Builder{}
	nameRegex, _ := regexp.Compile("^[a-zA-Z0-9]*$")
	if !nameRegex.MatchString(g.Name) {
		errors.WriteString(fmt.Sprintln("Group Name is invalid. Must be alphanumiric without spaces."))
	}

	errMsg := errors.String()
	if len(errMsg) > 0 {
		return false, errMsg
	}

	return true, ""
}

// GroupPermissions represents the join table between groups and permissions
type GroupPermissions struct {
	GroupID      uuid.UUID `gorm:"primaryKey;"`
	PermissionID string    `gorm:"primaryKey;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// GroupUsers represents the join table between groups and users
type GroupUsers struct {
	GroupID   uuid.UUID `gorm:"primaryKey;"`
	UserID    uuid.UUID `gorm:"primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
