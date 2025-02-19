package modules

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/gorm"
)

// GetPermissions will fetch all permissions
func (sdk *PGSDK) GetPermissions() ([]types.Permission, error) {
	permissions := []types.Permission{}
	results := sdk.db.Find(&permissions)
	return permissions, results.Error
}

// GetPermissionsByIDs will fetch all permissions in the provided list of permission IDs
func (sdk *PGSDK) GetPermissionsByIDs(permissionIDs []uuid.UUID) ([]types.Permission, error) {
	permissions := []types.Permission{}
	results := sdk.db.Model(&types.Permission{}).Find(&permissions, permissionIDs)
	return permissions, results.Error
}

// UpsertPermission will create or update a permission
func (sdk *PGSDK) UpsertPermission(permission types.Permission) (*types.Permission, error) {
	permissionIsValid, errMsg := permission.IsValid()
	if !permissionIsValid {
		return nil, fmt.Errorf("validation error: %s", errMsg)
	}

	// Perform an Upsert on the record, first by checking if the record already exists.
	if err := sdk.db.First(&types.Permission{ID: permission.ID}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := sdk.db.Create(&permission).Error; err != nil {
			return nil, err
		}
	} else {
		if err := sdk.db.Save(&permission).Error; err != nil {
			return nil, err
		}
	}

	return &permission, nil
}
