package providers

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/modules"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/gorm"
)

// StorageProvider is a port for the applications underlying storage/persistence layer
type StorageProvider struct {
	Session StorageSession
}

// Compile time proof of implementation
var _ IStorageProvider = (*StorageProvider)(nil)

// StorageSession represents a session with the postgres storage provider
type StorageSession struct {
	SDK modules.PGSDK
}

// InitStorageProvider will initialize the storage provider implementation.
func (p *ModuleProviders) InitStorageProvider() error {
	p.StorageProvider = &StorageProvider{
		Session: StorageSession{
			SDK: modules.InitPGDB(
				p.Config.DBAutoMigrate,
				p.Config.DBHost,
				p.Config.DBUserName,
				p.Config.DBPassword,
				p.Config.DBName,
				p.Config.Environment,
			),
		},
	}
	return nil
}

func (p *StorageProvider) GetPermissions() ([]types.Permission, error) {
	permissions, err := p.Session.SDK.GetPermissions()
	if err != nil {
		return []types.Permission{}, fmt.Errorf("unable to fetch permissions: %s", err.Error())
	}
	return permissions, nil
}

func (p *StorageProvider) GetPermissionsByIDs(permissionIDs []uuid.UUID) ([]types.Permission, error) {
	permissions, err := p.Session.SDK.GetPermissionsByIDs(permissionIDs)
	if err != nil {
		return []types.Permission{}, fmt.Errorf("unable to fetch permissions (using ID list): %s", err.Error())
	}
	return permissions, nil
}

func (p *StorageProvider) UpsertPermission(permission types.Permission) (types.Permission, error) {
	a, err := p.Session.SDK.UpsertPermission(permission)
	if err != nil {
		return types.Permission{}, fmt.Errorf("unable to upsert permission: %s", err.Error())
	}
	if a != nil {
		return *a, nil
	}
	return types.Permission{}, fmt.Errorf("permission upsert failed for unknown reason")
}

func (p *StorageProvider) GetGroups() ([]types.Group, error) {
	groups, err := p.Session.SDK.GetGroups()
	if err != nil {
		return []types.Group{}, fmt.Errorf("unable to fetch groups: %s", err.Error())
	}
	return groups, nil
}

func (p *StorageProvider) GetGroupsByIDs(groupIDs []uuid.UUID) ([]types.Group, error) {
	groups, err := p.Session.SDK.GetGroupsByIDs(groupIDs)
	if err != nil {
		return []types.Group{}, fmt.Errorf("unable to fetch groups (using ID list): %s", err.Error())
	}
	return groups, nil
}

func (p *StorageProvider) UpsertGroup(group types.Group) (types.Group, error) {
	a, err := p.Session.SDK.UpsertGroup(group)
	if err != nil {
		return types.Group{}, fmt.Errorf("unable to upsert group: %s", err.Error())
	}
	if a != nil {
		return *a, nil
	}
	return types.Group{}, fmt.Errorf("group upsert failed for unknown reason")
}

func (p *StorageProvider) GetUsers() ([]types.User, error) {
	users, err := p.Session.SDK.GetUsers()
	if err != nil {
		return []types.User{}, fmt.Errorf("unable to fetch users: %s", err.Error())
	}
	return users, nil
}

func (p *StorageProvider) GetUser(name string) (types.User, error) {
	user, err := p.Session.SDK.GetUser(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return types.User{}, err
	}

	if err != nil {
		return types.User{}, fmt.Errorf("unable to fetch user with name, %s: %s", name, err.Error())
	}
	return user, nil
}

func (p *StorageProvider) GetUserAccessDetails(userID uuid.UUID) (types.UserAccessDetails, error) {
	uad, err := p.Session.SDK.GetUserAccessDetails(userID)
	if err != nil {
		return types.UserAccessDetails{}, fmt.Errorf("unable to fetch user access details: %s", err.Error())
	}
	return uad, nil
}

func (p *StorageProvider) UpsertUser(user types.User) (types.User, error) {
	u, err := p.Session.SDK.UpsertUser(user)
	if err != nil {
		return types.User{}, fmt.Errorf("unable to upsert user: %s", err.Error())
	}
	if u != nil {
		return *u, nil
	}
	return types.User{}, fmt.Errorf("user upsert failed for unknown reason")
}

func (p *StorageProvider) GetMySQLCatalog() ([]*types.MySQLDBInfo, error) {
	dbInfo, err := p.Session.SDK.GetMySQLCatalog()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch MySQL DB catalog: %s", err.Error())
	}
	return dbInfo, nil
}

func (p *StorageProvider) UpsertMySQLDBInfo(dbInfo types.MySQLDBInfo) (*types.MySQLDBInfo, error) {
	dbi, err := p.Session.SDK.UpsertMySQLDBInfo(dbInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to upsert MySQL DB info: %s", err.Error())
	}
	if dbi != nil {
		return dbi, nil
	}
	return nil, fmt.Errorf("MySQL DB info upsert failed for unknown reason")
}

func (p *StorageProvider) DeleteMySQLDBInfo(dbHost string) error {
	err := p.Session.SDK.DeleteMySQLDBInfo(dbHost)
	if err != nil {
		return fmt.Errorf("unable to delete MySQL DB info: %s", err.Error())
	}
	return nil
}

func (p *StorageProvider) GetDynamicAppConfig() (types.DynamicAppConfig, error) {
	dac, err := p.Session.SDK.GetDynamicAppConfig()
	if err != nil {
		return types.DynamicAppConfig{}, fmt.Errorf("unable to fetch dynamic app config: %s", err.Error())
	}
	return dac, nil
}

func (p *StorageProvider) UpdateDynamicAppConfig(config types.DynamicAppConfig) (types.DynamicAppConfig, error) {
	return p.Session.SDK.UpdateDynamicAppConfig(config)
}
