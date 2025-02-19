package modules

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/gorm"
)

// GetGroups will fetch all groups
func (sdk *PGSDK) GetGroups() ([]types.Group, error) {
	groups := []types.Group{}
	results := sdk.db.Model(&types.Group{}).Preload("Permissions").Preload("Users").Find(&groups)
	return groups, results.Error
}

// GetGroupsByIDs will fetch all groups in the provided list of group IDs
func (sdk *PGSDK) GetGroupsByIDs(groupIDs []uuid.UUID) ([]types.Group, error) {

	groups := []types.Group{}
	results := sdk.db.Model(&types.Group{}).Preload("Permissions").Preload("Users").Find(&groups, groupIDs)
	return groups, results.Error
}

// UpsertGroup will create or update a group
func (sdk *PGSDK) UpsertGroup(group types.Group) (*types.Group, error) {
	groupIsValid, errMsg := group.IsValid()
	if !groupIsValid {
		return nil, fmt.Errorf("validation error: %s", errMsg)
	}

	existingGroup := types.Group{}
	err := sdk.db.Model(&types.Group{}).Preload("Permissions").Preload("Users").First(&existingGroup, "name = ?", group.Name).Error

	// If the record does not exist, just save a new one
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		if err := sdk.db.Save(&group).Error; err != nil {
			return nil, err
		}
		return &group, nil
	} else if err != nil {
		// If it was a different db error, return the error
		return nil, err
	} else {
		permissionAssociationsToRemove := identifyGroupPermissionsToRemove(existingGroup.Permissions, group.Permissions)
		userAssociationsToRemove := identifyGroupUsersToRemove(existingGroup.Users, group.Users)

		for _, a := range permissionAssociationsToRemove {
			sdk.db.Unscoped().Delete(&types.Permission{}, "group_id = ? AND permission_id = ?", existingGroup.ID, a)
		}

		for _, u := range userAssociationsToRemove {
			sdk.db.Unscoped().Delete(&types.GroupUsers{}, "group_id = ? AND user_id = ?", existingGroup.ID, u)
		}

		if err := sdk.db.Save(&group).Error; err != nil {
			return nil, err
		}
		return &group, nil
	}
}

func identifyGroupPermissionsToRemove(a, b []*types.Permission) []*uuid.UUID {
	bMap := make(map[uuid.UUID]bool)
	for _, item := range b {
		bMap[*item.ID] = true
	}

	// Create a slice to store the strings not found in list
	notFound := []*uuid.UUID{}
	for _, item := range a {
		if _, found := bMap[*item.ID]; !found {
			notFound = append(notFound, item.ID)
		}
	}
	return notFound
}

func identifyGroupUsersToRemove(a, b []*types.User) []uuid.UUID {
	bMap := make(map[uuid.UUID]bool)
	for _, item := range b {
		bMap[*item.ID] = true
	}

	// Create a slice to store the strings not found in list2
	notFound := []uuid.UUID{}
	for _, item := range a {
		if _, found := bMap[*item.ID]; !found {
			notFound = append(notFound, *item.ID)
		}
	}
	return notFound
}
