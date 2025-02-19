package modules

import (
	"fmt"

	"github.com/sullivtr/k8s_platform/internal/types"
)

// GetMySQLCatalog will fetch MYSQL databases from the catalog
func (sdk *PGSDK) GetMySQLCatalog() ([]*types.MySQLDBInfo, error) {
	dbCatalog := []*types.MySQLDBInfo{}
	results := sdk.db.Model(&types.MySQLDBInfo{}).Find(&dbCatalog)
	return dbCatalog, results.Error
}

// UpsertMySQLDBInfo will create or update a MySQL database in the catalog
func (sdk *PGSDK) UpsertMySQLDBInfo(dbInfo types.MySQLDBInfo) (*types.MySQLDBInfo, error) {
	mysqlDBInfoIsValid, errMsg := dbInfo.IsValid()
	if !mysqlDBInfoIsValid {
		return nil, fmt.Errorf("validation error: %s", errMsg)
	}

	if err := sdk.db.Save(&dbInfo).Error; err != nil {
		return nil, err
	}
	return &dbInfo, nil
}

// DeleteMySQLDBInfo will delete a MySQL database from the catalog
func (sdk *PGSDK) DeleteMySQLDBInfo(dbHost string) error {
	if dbHost == "" {
		return fmt.Errorf("invalid MySQLDBInfo: %v", dbHost)
	}

	if err := sdk.db.Unscoped().Where("host = ?", dbHost).Delete(&types.MySQLDBInfo{}).Error; err != nil {
		return err
	}
	return nil
}
