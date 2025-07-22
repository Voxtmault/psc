package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/voxtmault/psc/interfaces"
	"github.com/voxtmault/psc/models"
)

type KasusController struct {
	service interfaces.KasusInterface
}

func NewKasusController(paramService interfaces.KasusInterface) *KasusController {
	return &KasusController{
		service: paramService,
	}
}

func (cc *KasusController) Get(c echo.Context) error {
	jenisMasalahID := c.FormValue("jenis_masalah_id")

	// conversi dari string ke uint
	convJMID, err := strconv.Atoi(jenisMasalahID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status_code": http.StatusBadRequest,
			"message":     "invalid jenis_masalah_id format",
		})
	}

	result, err := cc.service.Get(c.Request().Context(), &models.KasusFilter{
		JenisMasalahID: uint(convJMID),
	})
	if err != nil {
		return echo.NewHTTPError(result.StatusCode, result)
	}

	return c.JSON(result.StatusCode, result)
}

func (cc *KasusController) Get2(c echo.Context) error {
	var filter models.KasusFilter

	jenisMasalahID := c.FormValue("jenis_masalah_id")
	// conversi dari string ke uint
	convJMID, err := strconv.Atoi(jenisMasalahID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status_code": http.StatusBadRequest,
			"message":     "invalid jenis_masalah_id format",
		})
	}

	filter.JenisMasalahID = uint(convJMID)

	result, err := cc.service.Get(c.Request().Context(), &filter)
	if err != nil {
		return echo.NewHTTPError(result.StatusCode, result)
	}

	return c.JSON(result.StatusCode, result)
}

func (cc *KasusController) Get3(c echo.Context) error {
	var filter models.KasusFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
			"status_code": http.StatusUnprocessableEntity,
			"message":     "unprocessable request",
		})
	}

	if err := c.Validate(filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status_code": http.StatusBadRequest,
			"message":     "invalid jenis_masalah_id format",
		})
	}

	result, err := cc.service.Get(c.Request().Context(), &filter)
	if err != nil {
		return echo.NewHTTPError(result.StatusCode, result)
	}

	return c.JSON(result.StatusCode, result)
}

func (cc *KasusController) Create(c echo.Context) error {
	var payload models.KasusCreate
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
			"status_code": http.StatusUnprocessableEntity,
			"message":     "unprocessable request",
		})
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status_code": http.StatusBadRequest,
			"message":     "invalid jenis_masalah_id format",
		})
	}

	result, err := cc.service.Create(c.Request().Context(), &payload)
	if err != nil {
		return echo.NewHTTPError(result.StatusCode, result)
	}

	return c.JSON(result.StatusCode, result)
}
