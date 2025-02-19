package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// DynamicConfigJSONB is a custom type for JSONB fields in the database
type DynamicConfigJSONB struct {
	DefaultReplicaScaleLimit int            `json:"defaultReplicaScaleLimit"`
	ReplicaScaleLimits       map[string]int `json:"replicaScaleLimits"`
	EnableK8sGlobalReadOnly  bool           `json:"enableK8sGlobalReadOnly"`
	K8sClusterName           string         `json:"k8sClusterName"`
	K8sClusterNamespaces     []string       `json:"k8sClusterNamespaces"`
}

// Value Marshal
func (jsonField DynamicConfigJSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *DynamicConfigJSONB) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}

type DynamicAppConfig struct {
	ID   uint               `json:"id" gorm:"primaryKey"`
	Data DynamicConfigJSONB `json:"data" gorm:"type:jsonb"`
}
