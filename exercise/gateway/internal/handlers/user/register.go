package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/amodemoli/fastic/core/color"
	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	pbUserService "github.com/amodemoli/microservices/exercise/user/protobuf"
	"github.com/valyala/fasthttp"
)

type Client struct {
	// target service gRCP server url
	Target string

	// Context of request
	Ctx context.Context
	// Maded fasthttp client
	FasthttpClient *fasthttp.Client
}

// Register function, create's client for connect to user service,
// and handle user request path's here
func Register(
	// application struct, for print errors on terminal
	app *fastic.App,
	// context (can be with timeout)
	ctx context.Context,
	// user-service target host url
	target string,
	// get maded fasthttp client from main.go
	fasthttpClient *fasthttp.Client,
	// middlewares, can be empty use <nil>
	middlewares []func(ctx context.Context, req *fasthttp.Request, handler func(ctx context.Context, req *fasthttp.Request) (resp *fasthttp.Response, err error)) (resp *fasthttp.Response, err error)) error {

	// validate target service url length
	if target == "" {
		// create error message for use
		message := "user service target-url cannot be empty"
		// check server development mode status
		// for security, im check if development_mode is off i ignore this error else: exit from application with 1 error code.!
		if app.Env.DevelopemtMode {
			// development mode is true, only write erorr on terminal
			helpers.Print(app, color.Yellow, "WARNING", message)
			return errors.New(message) // return because other section dont works now
		}
		// development mode is false, server on production mode.! exit from server and show error
		// this is only for security =D
		log.Fatalf("%s, change it on main.go first!\nenable development_mode for ignore this error", message)
	}

	// get user client from protobuf of user service
	userClient, err := pbUserService.GetUserHTTPGoClient(ctx, fasthttpClient, target, nil)
	if err != nil {
		helpers.Print(app, color.Red, "CONN-ERROR", fmt.Sprintf("user-service connection failed: %v", err))
		// return from register function (dont need to connect.)
		// i used return because other next step's of code don't works. exiting now without registering
		return err
	}

	// connection is succsess with no errors, register user paths.

	app.Get("/api/user/{id}", func(c *fastic.Ctx) {
		GetUser(app, c, userClient)
	})

	// return from function after registering all paths
	return nil
}
