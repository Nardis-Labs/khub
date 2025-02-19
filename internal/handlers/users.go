package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
)

type UsersHandler struct {
	provider *providers.ModuleProviders
}

// Get godoc
// @Summary Get Users
// @Description get Users
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.User
// @Router /api/users [get]
func (c UsersHandler) GetUsers(ctx echo.Context) error {
	users, err := c.provider.StorageProvider.GetUsers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to get users %s", err.Error()))
	}
	return ctx.JSON(http.StatusOK, users)
}

// Get godoc
// @Summary Update User Theme Preference
// @Description update User Theme Preference
// @Tags Users
// @Accept  json
// @Produce  json
// @Param name path string true "User Name"
// @Param darkMode query string false "Dark Mode"
// @Success 200 {object} []types.User
// @Router /api/users/theme/{name} [put]
func (c UsersHandler) UpdateUserThemePreference(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, "name is cannot be empty")
	}

	darkModePreference := ctx.QueryParam("darkMode")
	if darkModePreference == "" {
		return ctx.JSON(http.StatusBadRequest, "darkMode is cannot be empty")
	}

	user, err := c.provider.StorageProvider.GetUser(name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to get user with name, %s %s", name, err.Error()))
	}

	if darkModePreference == "true" {
		fmt.Println("dark mode engaged")
		user.DarkMode = true
	} else {
		user.DarkMode = false
	}
	if _, err := c.provider.StorageProvider.UpsertUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to update user theme prefence, (%s) %s", name, err.Error()))
	}

	return ctx.JSON(http.StatusOK, user)
}
