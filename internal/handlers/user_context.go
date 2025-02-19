package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/gorm"
)

// GetUserContextUpsert will fetch the user indicated by the request context
// If the user does not exist, the user is automatically created
// Otherwise, the user's lastUsed timestamp is updated
func GetUserContextUpsert(ctx echo.Context, storageProvider *providers.StorageProvider) types.User {
	username := ctx.Get("username").(string)
	userEmail := ctx.Get("email").(string)

	user := types.User{}
	var err error
	if username != "" && userEmail != "" {
		user, err = storageProvider.GetUser(username)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			user = types.User{
				Name:     username,
				Email:    userEmail,
				LastUsed: time.Now(),
				DarkMode: true,
			}
			_, err := storageProvider.UpsertUser(user)
			if err != nil {
				log.Warn().Msg("auto user upsert failed during user fetch process")
			}
		} else {
			user.LastUsed = time.Now()
			_, err := storageProvider.UpsertUser(user)
			if err != nil {
				log.Warn().Msg("auto update user last used time failed during user fetch process")
			}
		}
	}
	return user
}

// GetUserContext will fetch the user indicated by the request context
func GetUserContext(ctx echo.Context, storageProvider *providers.StorageProvider) (types.User, int, error) {
	username := ctx.Get("username").(string)
	userEmail := ctx.Get("email").(string)

	if username != "" && userEmail != "" {
		user, err := storageProvider.GetUser(username)
		return user, http.StatusOK, err
	}
	return types.User{}, http.StatusForbidden, errors.New("unable to read user context details from request (unauthenticated)")
}

func GetUserPermissions(ctx echo.Context, storageProvider *providers.StorageProvider, user *types.User, enableGlobalReadOnly bool) ([]types.Permission, error) {
	if user == nil {
		userFetched, _, err := GetUserContext(ctx, storageProvider)
		if err != nil {
			return []types.Permission{}, err
		}
		user = &userFetched
	}

	if user.IsAdmin {
		// If user is an admin, they should have access to all permissions. This will gather every permission.
		permissions, err := storageProvider.GetPermissions()
		if err != nil {
			return []types.Permission{}, fmt.Errorf("unable to get permissions for user %s", err.Error())
		}
		return permissions, nil
	} else {
		if user.ID == nil {
			return []types.Permission{}, errors.New("user ID is nil")
		}
		userAccessDetail, err := storageProvider.GetUserAccessDetails(*user.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to get user access detail %s", err.Error())
		}

		// globalReadOnlyPermission represents the permission granted to all users when enableGlobalReadOnly is true
		globalReadOnlyPermission := types.Permission{
			Name:   "global_read_only",
			AppTag: "global_read_only",
		}

		// If the user does not have any group access, and global read only is not enabled, they should not have any permissions.
		// If global read only is enabled, they should ONLY have the global read only permission.
		if len(userAccessDetail.GroupIDs) == 0 && !enableGlobalReadOnly {
			return nil, errors.New("user does not have permissions")
		} else if len(userAccessDetail.GroupIDs) == 0 && enableGlobalReadOnly {
			return []types.Permission{globalReadOnlyPermission}, nil
		}

		permissions, err := storageProvider.GetPermissionsByIDs(userAccessDetail.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("unable to get permissions for user by IDs %s", err.Error())
		}

		// Ensure that the global read only permission is included if global read only is enabled for all users
		if enableGlobalReadOnly {
			permissions = append(permissions, globalReadOnlyPermission)
		}
		return permissions, nil
	}
}

func FilterUserPermissionAccess(permissionName string, permissions []types.Permission) ([]types.Permission, error) {
	if permissionName != "" {
		filteredPermissions := []types.Permission{}
		for _, v := range permissions {
			if v.Name == permissionName {
				filteredPermissions = append(filteredPermissions, v)
			}
		}
		return filteredPermissions, nil
	} else {
		return nil, errors.New("unable to filter user permissions")
	}
}

func GetUserGroups(ctx echo.Context, storageProvider *providers.StorageProvider, user *types.User) ([]types.Group, error) {
	if user == nil {
		userFetched, _, err := GetUserContext(ctx, storageProvider)
		if err != nil {
			return nil, err
		}
		user = &userFetched
	}

	if user.IsAdmin {
		groups, err := storageProvider.GetGroups()
		if err != nil {
			return []types.Group{}, fmt.Errorf("unable to get groups %s", err.Error())
		}
		return groups, nil
	} else {
		userAccessDetail, err := storageProvider.GetUserAccessDetails(*user.ID)
		if err != nil {
			return []types.Group{}, fmt.Errorf("unable to get user group detail %s", err.Error())
		}
		if len(userAccessDetail.GroupIDs) == 0 {
			return []types.Group{}, errors.New("user does not have access to any groups")
		}
		groups, err := storageProvider.GetGroupsByIDs(userAccessDetail.GroupIDs)
		if err != nil {
			return []types.Group{}, fmt.Errorf("unable to get user's groups %s", err.Error())
		}
		return groups, nil
	}
}
