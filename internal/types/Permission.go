package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a permission on the khub application.
// The permission is used to indicate which app's can be accessed by a group.
type Permission struct {
	ID        *uuid.UUID     `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex"`
	AppTag    string         `json:"appTag" gorm:"uniqueIndex"`
	Groups    []*Group       `json:"groups" gorm:"many2many:group_permissions;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *Permission) IsValid() (bool, string) {
	errors := strings.Builder{}
	nameRegex, _ := regexp.Compile("^[a-zA-Z0-9]*$")
	if !nameRegex.MatchString(r.Name) {
		errors.WriteString(fmt.Sprintln("Permission Name is invalid. Must be alphanumiric without spaces."))
	}

	errMsg := errors.String()
	if len(errMsg) > 0 {
		return false, errMsg
	}

	return true, ""
}
