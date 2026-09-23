package handlers

import (
	"fmt"
	"strings"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/logger"
)

func HandleServerErrors(app *fastic.App, lg *logger.Logger, c *fastic.Ctx, service string, err error) {
	prefix := "service-errors"

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
	lg.Error(prefix, fmt.Sprintf("%s service error: %v [handlers/server.errors.go]", service, err))

}

func errIsTimeout(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
