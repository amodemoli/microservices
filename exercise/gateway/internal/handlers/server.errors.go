package handlers

import (
	"strings"

	"github.com/amodemoli/fastic/core/color"
	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
)

func HandleServerErrors(app *fastic.App, c *fastic.Ctx, err error) {
	if errIsTimeout(err) {
		c.Status(504)
		helpers.Json(c, helpers.Options{
			Status:  status.Timeout,
			Message: "service temporarily unavailable",
			Code:    response.ServiceUnavailable,
		})
		return
	}

	c.Status(500)
	helpers.Json(c, helpers.Options{
		Status:  status.Error,
		Message: "unknown error, check terminal for see more...",
		Code:    response.InternalServerError,
	})
	helpers.Print(app, color.Red, "USER-SERVICE-ERROR", err.Error())

}

func errIsTimeout(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
