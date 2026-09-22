package user

import (
	"fmt"
	"strconv"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/handlers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/logger"
	pbUserService "github.com/amodemoli/microservices/exercise/user/protobuf"
)

func GetUser(app *fastic.App, c *fastic.Ctx, lg *logger.Logger, client *pbUserService.UserHTTPGoClient) {

	id, err := strconv.ParseInt(fmt.Sprint(c.UserValue("id")), 10, 64)
	if err != nil {
		c.Status(400)
		helpers.Json(c, helpers.Options{
			Status:  status.Invalid,
			Message: "invalid id",
			Code:    response.InvalidInput,
		})
		return
	}

	resp, err := client.GetUser(c.RequestCtx, &pbUserService.GetUserRequest{Id: id})
	if err != nil {
		// error from grpc
		handlers.HandleServerErrors(app, lg, c, "user", err)
		return
	}

	// user error
	if resp.Response.Error != "" {
		c.Status(response.Code(resp.Response.Code).ToHTTPStatus())
		helpers.Json(c, helpers.Options{
			Status:  status.Status(resp.Response.Status),
			Message: resp.Response.Error,
			Code:    response.Code(resp.Response.Code),
		})
		return
	}

	c.Status(response.Code(resp.Response.Code).ToHTTPStatus())
	helpers.Json(c, helpers.Options{
		Status:  status.Status(resp.Response.Status),
		Message: "done.",
		Code:    response.Code(resp.Response.Code),
		Data:    resp.User,
	})
}
