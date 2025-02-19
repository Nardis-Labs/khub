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

type GroupsHandler struct {
	provider *providers.ModuleProviders
}

// Get godoc
// @Summary Get Groups
// @Description get Groups
// @Tags Groups
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.Group
// @Router /api/groups [get]
func (c GroupsHandler) GetGroups(ctx echo.Context) error {
	groups, err := GetUserGroups(ctx, c.provider.StorageProvider, nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, groups)
}

// Get godoc
// @Summary Upsert Groups
// @Description upsert Groups
// @Tags Groups
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Group
// @Router /api/groups [put]
func (c GroupsHandler) UpsertGroup(ctx echo.Context) error {
	groups, err := GetUserGroups(ctx, c.provider.StorageProvider, nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var group types.Group
	err = json.NewDecoder(ctx.Request().Body).Decode(&group)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to decode group json body %s", err.Error()))
	}

	if group.ID != nil && *group.ID == uuid.Nil {
		gid := uuid.New()
		group.ID = &gid
	} else {
		groupFound := false
		for _, g := range groups {
			if g.ID != nil && *g.ID == *group.ID {
				groupFound = true
				break
			}
		}
		if !groupFound {
			return ctx.JSON(http.StatusUnauthorized, "unable to upsert group. You do not have permission to update this group.")
		}
	}

	g, err := c.provider.StorageProvider.UpsertGroup(group)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to upsert group %s", err.Error()))
	}
	return ctx.JSON(http.StatusOK, g)
}
