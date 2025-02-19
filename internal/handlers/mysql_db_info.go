package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type MySQLDBInfoHandler struct {
	provider *providers.ModuleProviders
}

// GetMySQLDBCatalog godoc
// @Summary Get MySQL DB Info Catalog
// @Description get mysql db info catalog
// @Tags MySQLDBInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/infra/mysql [get]
func (c MySQLDBInfoHandler) GetMySQLDBCatalog(ctx echo.Context) error {
	userDetail, _, err := GetUserContext(ctx, c.provider.StorageProvider)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "user context unavailable")
	}

	if !userDetail.IsAdmin {
		return ctx.JSON(http.StatusForbidden, "user must be a global admin to access this mysql database catalog")
	}

	catalog, err := c.provider.StorageProvider.GetMySQLCatalog()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, catalog)
}

// UpsertMySQLDBInfo godoc
// @Summary Upsert MySQL DB Info
// @Description upsert mysql db info
// @Tags MySQLDBInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} types.MySQLDBInfo
// @Router /api/infra/mysql [put]
func (c MySQLDBInfoHandler) UpsertMySQLDBInfo(ctx echo.Context) error {
	userDetail, _, err := GetUserContext(ctx, c.provider.StorageProvider)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "user context unavailable")
	}

	if !userDetail.IsAdmin {
		return ctx.JSON(http.StatusForbidden, "user must be a global admin to update mysql database catalog items")
	}

	var mysqlDBInfo types.MySQLDBInfo
	err = json.NewDecoder(ctx.Request().Body).Decode(&mysqlDBInfo)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	dbInfo, err := c.provider.StorageProvider.UpsertMySQLDBInfo(mysqlDBInfo)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to upsert mysql db info: %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, dbInfo)
}

// DeleteMySQLDBInfo godoc
// @Summary Delete MySQL DB Info
// @Description delete mysql db info
// @Tags MySQLDBInfo
// @Accept  json
// @Produce  json
// @NoContent 204 {object} string
// @Param dbHost query string true "dbHost"
// @Router /api/infra/mysql [delete]
func (c MySQLDBInfoHandler) DeleteMySQLDBInfo(ctx echo.Context) error {
	userDetail, _, err := GetUserContext(ctx, c.provider.StorageProvider)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "user context unavailable")
	}

	if !userDetail.IsAdmin {
		return ctx.JSON(http.StatusForbidden, "user must be a global admin to delete mysql database catalog items")
	}

	dbHost := ctx.QueryParam("dbHost")
	if dbHost == "" {
		return ctx.JSON(http.StatusBadRequest, "dbHost query param is required")
	}

	err = c.provider.StorageProvider.DeleteMySQLDBInfo(dbHost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to delete mysql db info: %s", err.Error()))
	}

	return ctx.NoContent(http.StatusNoContent)
}

// GetReplicationTopology godoc
// @Summary Get MySQL Replication Topology
// @Description get mysql replication topology
// @Tags MySQLDBInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/infra/mysql/topology [get]
func (c MySQLDBInfoHandler) GetReplicationTopology(ctx echo.Context) error {
	nodes, err := c.provider.CacheProvider.Get("mysql-repl-topo-nodes")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to get mysql replication topology nodes: %s", err.Error()))
	}

	edges, err := c.provider.CacheProvider.Get("mysql-repl-topo-edges")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to get mysql replication topology edges: %s", err.Error()))
	}

	topoGraph := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	return ctx.JSON(http.StatusOK, topoGraph)
}
