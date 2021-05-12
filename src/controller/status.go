package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

// Get is used for know Service status
func Status(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}
