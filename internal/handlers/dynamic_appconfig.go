package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type DynamicAppConfigHandler struct {
	provider *providers.ModuleProviders
}

// GetDynamicAppConfig godoc
// @Summary Get DynamicAppConfig
// @Description get DynamicAppConfig
// @Tags DynamicAppConfig
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.DynamicAppConfig
// @Router /api/appconfig [get]
func (c DynamicAppConfigHandler) GetDynamicAppConfig(ctx echo.Context) error {
	dac, err := c.provider.StorageProvider.GetDynamicAppConfig()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, dac)
}

// UpdateDynamicAppConfig godoc
// @Summary Update DynamicAppConfig
// @Description update DynamicAppConfig
// @Tags DynamicAppConfig
// @Accept  json
// @Produce  json
// @Success 200 {object} types.DynamicAppConfig
// @Router /api/appconfig [put]
func (c DynamicAppConfigHandler) UpdateDynamicAppConfig(ctx echo.Context) error {
	user, _, err := GetUserContext(ctx, c.provider.StorageProvider)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}

	if !user.IsAdmin {
		return ctx.JSON(http.StatusForbidden, "user must be an admin to edit dynamic app config")
	}

	var dac types.DynamicAppConfig
	err = json.NewDecoder(ctx.Request().Body).Decode(&dac)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to decode dynamic app config json body %s", err.Error()))
	}

	_, err = c.provider.StorageProvider.UpdateDynamicAppConfig(dac)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dac)
}
