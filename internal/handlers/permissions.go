package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type PermissionsHandler struct {
	provider *providers.ModuleProviders
}

// Get godoc
// @Summary Get Permissions
// @Description get Permissions
// @Tags Permissions
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.Permission
// @Router /api/permissions [get]
func (c PermissionsHandler) GetPermissions(ctx echo.Context) error {
	permissions, err := c.provider.StorageProvider.GetPermissions()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, permissions)
}

// Get godoc
// @Summary Upsert Permissions
// @Description get Permissions
// @Tags Permissions
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Permission
// @Router /api/permissions [put]
func (c PermissionsHandler) UpsertPermission(ctx echo.Context) error {
	var permission types.Permission
	err := json.NewDecoder(ctx.Request().Body).Decode(&permission)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to decode permission json body %s", err.Error()))
	}

	if permission.ID != nil && *permission.ID == uuid.Nil {
		pid := uuid.New()
		permission.ID = &pid
	}

	a, err := c.provider.StorageProvider.UpsertPermission(permission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to upsert permission %s", err.Error()))
	}
	return ctx.JSON(http.StatusOK, a)
}
