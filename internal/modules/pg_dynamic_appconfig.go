package modules

import (
	"fmt"

	"github.com/sullivtr/k8s_platform/internal/types"
)

func (sdk *PGSDK) GetDynamicAppConfig() (types.DynamicAppConfig, error) {
	config := types.DynamicAppConfig{}
	// There will only ever be one record for dynamic app config
	if err := sdk.db.Where("id = 1").First(&config).Error; err != nil {
		return types.DynamicAppConfig{}, err
	}
	return config, nil
}

// GetDynamicAppConfig will fetch the dynamic app config
func (sdk *PGSDK) UpdateDynamicAppConfig(config types.DynamicAppConfig) (types.DynamicAppConfig, error) {
	// There will only ever be one record for dynamic app config
	// If an attempt is made to update the wrong record, return an error
	if config.ID != 1 {
		return types.DynamicAppConfig{}, fmt.Errorf("invalid dynamic app config ID: %v", config.ID)
	}

	if err := sdk.db.Save(&config).Error; err != nil {
		return types.DynamicAppConfig{}, err
	}
	return config, nil
}
