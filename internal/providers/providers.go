package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/types"
	v1 "k8s.io/api/core/v1"
)

// Compile-time proof of interface implementation.
var _ IModuleProviders = (*ModuleProviders)(nil)

// IModuleProviders is the interface representation of a module providers library.
//
// A module provider represents a (mostly) technology agnostic entry point for underlying
// systems & services (both internal and external) such as cloud platforms, datastores, caches etc.
type IModuleProviders interface {
	InitK8sProvider()
	InitCacheProvider() error
	InitStorageProvider() error
	InitAWSProvider()
	InitMySQLTopoProvider()
	StartDataSink(ctx context.Context, intervalSeconds int)
}

// IAWSProvider is an interface representing functionality for an AWS provider
type IAWSProvider interface {
	GetS3Reports() ([]types.Report, error)
	GetReportDownloadURL(reportName string) (string, error)
}

// IK8sProvider is an interface representing functionality for a kubernetes provider
type IK8sProvider interface {
	GetPods(namespaces []string) (any, error)
	GetPod(namespace, podName string) (*v1.Pod, error)
	DeletePod(namespace, podName string) error
	GetDeployments(namespaces []string) (any, error)
	ScaleDeployment(namespace, deploymentName string, replicas int32) error
	GetDaemonsets(namespaces []string) (any, error)
	GetReplicasets(namespaces []string) (any, error)
	GetStatefulsets(namespaces []string) (any, error)
	GetJobs(namespaces []string) (any, error)
	GetCronJobs(namespaces []string) (any, error)
	GetServices(namespaces []string) (any, error)
	GetIngresses(namespaces []string) (any, error)
	GetConfigMaps(namespaces []string) (any, error)
	GetNodes() (any, error)
	GetClusterEvents() (any, error)
	RolloutRestartDeployment(string, string) error
	RolloutRestartDaemonSet(string, string) error
	RolloutRestartStatefulSet(string, string) error
}

// IStorageProvider is an interface representing functionality for a storage/persistence provider
type IStorageProvider interface {
	GetPermissions() ([]types.Permission, error)
	GetPermissionsByIDs(accountIDs []uuid.UUID) ([]types.Permission, error)
	UpsertPermission(account types.Permission) (types.Permission, error)
	GetGroups() ([]types.Group, error)
	GetGroupsByIDs(groupIDs []uuid.UUID) ([]types.Group, error)
	UpsertGroup(group types.Group) (types.Group, error)
	GetUsers() ([]types.User, error)
	GetUser(name string) (types.User, error)
	GetUserAccessDetails(userID uuid.UUID) (types.UserAccessDetails, error)
	UpsertUser(user types.User) (types.User, error)
	GetMySQLCatalog() ([]*types.MySQLDBInfo, error)
	UpsertMySQLDBInfo(dbInfo types.MySQLDBInfo) (*types.MySQLDBInfo, error)
	DeleteMySQLDBInfo(dbHost string) error
	GetDynamicAppConfig() (types.DynamicAppConfig, error)
	UpdateDynamicAppConfig(config types.DynamicAppConfig) (types.DynamicAppConfig, error)
}

// ICacheProvider is an interface representing functionality for a storage/persistence provider
type ICacheProvider interface {
	Get(key string) (any, error)
	Put(key string, value any) error
	InitAuthSessionStore() *redisstore.RedisStore
}

// IMySQLTopoProvider is an interface representing functionality for a MySQL topology provider
type IMySQLTopoProvider interface {
	CaptureReplicationTopology() ([]types.ReplTopoTreeNode, []types.ReplTopoTreeEdge, error)
}

// ModuleProviders is a struct containing the collection of known providers to
// be made available for use at runtime.
type ModuleProviders struct {
	Config            *config.Config
	CacheProvider     *CacheProvider
	K8sProvider       *K8sApiProvider
	StorageProvider   *StorageProvider
	MySQLTopoProvider *MySQLTopoProvider
	AWSProvider       *AWSProvider
}
