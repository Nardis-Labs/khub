package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
)

var (
	upgrader = websocket.Upgrader{}
)

func RegisterRoutes(e *echo.Echo, prv *providers.ModuleProviders) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	if err := registerK8sResources(e, prv); err != nil {
		return err
	}

	if err := registerAdministrationRoutes(e, prv); err != nil {
		return err
	}

	if err := registerAuthRoutes(e, prv); err != nil {
		return err
	}
	return nil
}

func registerK8sResources(e *echo.Echo, prv *providers.ModuleProviders) error {
	k8sHandler := &K8sSessionHandler{provider: prv}
	e.GET("/api/k8s/name", k8sHandler.GetClusterName)
	e.GET("/api/k8s/pods", k8sHandler.GetPods)
	e.DELETE("/api/k8s/pods", k8sHandler.DeletePod)
	e.GET("/api/k8s/deployments", k8sHandler.GetDeployments)
	e.POST("/api/k8s/deployments/scale", k8sHandler.ScaleDeployment)
	e.GET("/api/k8s/replicasets", k8sHandler.GetReplicasets)
	e.GET("/api/k8s/daemonsets", k8sHandler.GetDaemonsets)
	e.GET("/api/k8s/statefulsets", k8sHandler.GetStatefulsets)
	e.GET("/api/k8s/jobs", k8sHandler.GetJobs)
	e.GET("/api/k8s/cronjobs", k8sHandler.GetCronJobs)
	e.GET("/api/k8s/services", k8sHandler.GetServices)
	e.GET("/api/k8s/ingresses", k8sHandler.GetIngresses)
	e.GET("/api/k8s/configmaps", k8sHandler.GetConfigMaps)
	e.GET("/api/k8s/nodes", k8sHandler.GetNodes)
	e.GET("/api/k8s/clusterevents", k8sHandler.GetClusterEvents)
	e.POST("/api/k8s/rolloutrestart", k8sHandler.RolloutRestart)
	e.POST("/api/k8s/threaddump", k8sHandler.TakeTomcatThreadDump)
	return nil
}

func registerAdministrationRoutes(e *echo.Echo, prv *providers.ModuleProviders) error {
	usersHandler := &UsersHandler{provider: prv}
	e.GET("/api/users", usersHandler.GetUsers)
	e.PUT("/api/users/theme/:name", usersHandler.UpdateUserThemePreference)

	groupsHanlder := &GroupsHandler{provider: prv}
	e.GET("/api/groups", groupsHanlder.GetGroups)
	e.PUT("/api/groups", groupsHanlder.UpsertGroup)

	permissionsHandler := &PermissionsHandler{provider: prv}
	e.GET("/api/permissions", permissionsHandler.GetPermissions)
	e.PUT("/api/permissions", permissionsHandler.UpsertPermission)

	reportsHandler := &ReportsHandler{provider: prv}
	e.GET("/api/reports", reportsHandler.GetReports)
	e.GET("/api/reports/download/:key", reportsHandler.GetReportDownloadURL)

	mySQLDBInfoHandler := &MySQLDBInfoHandler{provider: prv}
	e.GET("/api/infra/mysql", mySQLDBInfoHandler.GetMySQLDBCatalog)
	e.PUT("/api/infra/mysql", mySQLDBInfoHandler.UpsertMySQLDBInfo)
	e.DELETE("/api/infra/mysql", mySQLDBInfoHandler.DeleteMySQLDBInfo)
	e.GET("/api/infra/mysql/topology", mySQLDBInfoHandler.GetReplicationTopology)

	dynamicAppConfigHandler := &DynamicAppConfigHandler{provider: prv}
	e.GET("/api/appconfig", dynamicAppConfigHandler.GetDynamicAppConfig)
	e.PUT("/api/appconfig", dynamicAppConfigHandler.UpdateDynamicAppConfig)

	return nil
}

func registerAuthRoutes(e *echo.Echo, prv *providers.ModuleProviders) error {
	authHandler := &AuthSessionHandler{
		provider: prv,
	}
	e.GET("/login", authHandler.Login)
	e.GET("/logout", authHandler.Logout)
	e.GET("/authorization-code/callback", authHandler.AuthCodeCallback)
	e.GET("/api/users/me", authHandler.UserInfo)

	return nil
}
