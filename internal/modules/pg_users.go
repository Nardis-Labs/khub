package modules

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
)

// GetUsers will fetch all users
func (sdk *PGSDK) GetUsers() ([]types.User, error) {
	users := []types.User{}
	results := sdk.db.Find(&users)
	return users, results.Error
}

// GetUser will fetch a single user
func (sdk *PGSDK) GetUser(username string) (types.User, error) {
	user := types.User{}
	results := sdk.db.Preload("Groups").Where("name = ?", username).First(&user)
	return user, results.Error
}

// GetUserAccessDetails will fetch a list of accounts and groups that a user can access (for authorization)
func (sdk *PGSDK) GetUserAccessDetails(userID uuid.UUID) (types.UserAccessDetails, error) {
	groupIDs := []uuid.UUID{}
	groupIDResults := sdk.db.Raw("SELECT group_id FROM group_users WHERE user_id = ?", userID).Scan(&groupIDs)
	if groupIDResults.Error != nil {
		return types.UserAccessDetails{}, groupIDResults.Error
	}

	// short circuit the user does not have access to any groups
	if len(groupIDs) == 0 {
		return types.UserAccessDetails{
			UserID:        userID,
			GroupIDs:      groupIDs,
			PermissionIDs: []uuid.UUID{},
		}, nil
	}

	groups := []types.Group{}
	if groupResultsError := sdk.db.Model(&types.Group{}).Preload("Permissions").Find(&groups, groupIDs).Error; groupResultsError != nil {
		return types.UserAccessDetails{}, groupIDResults.Error
	}

	permissionIDs := []uuid.UUID{}
	for _, g := range groups {
		for _, p := range g.Permissions {
			permissionIDs = append(permissionIDs, *p.ID)
		}
	}

	return types.UserAccessDetails{
		UserID:        userID,
		GroupIDs:      groupIDs,
		PermissionIDs: permissionIDs,
	}, nil
}

// UpsertUser will create or update a user
func (sdk *PGSDK) UpsertUser(user types.User) (*types.User, error) {
	userIsValid, errMsg := user.IsValid()
	if !userIsValid {
		return nil, fmt.Errorf("validation error: %s", errMsg)
	}

	if err := sdk.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
