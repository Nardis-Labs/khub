package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user on the khub application.
type User struct {
	ID        *uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"uniqueIndex"`
	Email     string         `json:"email"`
	IsAdmin   bool           `json:"isAdmin"`
	LastUsed  time.Time      `json:"lastUsed"`
	DarkMode  bool           `json:"darkMode"`
	Groups    []*Group       `json:"groups" gorm:"many2many:group_users;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// UserAccessDetails represents the access a user has, including their groups and permissions
// Permission access is based on the groups the user is part of.
type UserAccessDetails struct {
	UserID        uuid.UUID   `json:"userId"`
	GroupIDs      []uuid.UUID `json:"groupIds"`
	PermissionIDs []uuid.UUID `json:"permissionTags"`
}

func (u *User) IsValid() (bool, string) {
	errors := strings.Builder{}
	emailRegex, _ := regexp.Compile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		errors.WriteString(fmt.Sprintln("User Email is invalid. Must be valid email format."))
	}

	if len(u.Name) > 50 {
		errors.WriteString(fmt.Sprintln("User Name is invalid. Name cannot exceed 50 characters."))
	}

	nameRegex, _ := regexp.Compile(`^[\w\.\-\s]+$`)
	if !nameRegex.MatchString(u.Name) {
		errors.WriteString(
			fmt.Sprintln("User Name is invalid. Name must be alphanumeric (can only contain the following special characters: - or .)"),
		)
	}

	errMsg := errors.String()
	if len(errMsg) > 0 {
		return false, errMsg
	}

	return true, ""
}
