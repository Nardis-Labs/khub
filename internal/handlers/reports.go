package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sullivtr/k8s_platform/internal/providers"
)

type ReportsHandler struct {
	provider *providers.ModuleProviders
}

// Get Reports godoc
// @Summary Get Reports
// @Description Get Reports
// @Tags Reports
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.Report
// @Router /api/reports [get]
func (c ReportsHandler) GetReports(ctx echo.Context) error {
	reports, err := c.provider.AWSProvider.GetS3Reports()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to fetch reports, %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, reports)
}

// Get Reports Download godoc
// @Summary Get Reports Download URL
// @Description Get Reports Download URL
// @Tags Reports
// @Accept  json
// @Produce  json
// @Param key path string true "Report Name"
// @Success 200 {object} string
// @Router /api/reports/download/{key} [get]
func (c ReportsHandler) GetReportDownloadURL(ctx echo.Context) error {
	key := ctx.Param("key")
	if key == "" {
		return ctx.JSON(http.StatusBadRequest, "key is cannot be empty")
	}

	url, err := c.provider.AWSProvider.GetReportDownloadURL(key)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to fetch report download url, %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, url)
}
