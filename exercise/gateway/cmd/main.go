package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amodemoli/fastic/core/color"
	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/amodemoli/microservices/exercise/gateway/internal/handlers/user"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/logger"
	"github.com/amodemoli/microservices/exercise/gateway/internal/middleware"
	"github.com/valyala/fasthttp"
)

// make httpClient with fasthttp
var httpClient = &fasthttp.Client{
	MaxIdleConnDuration: 10 * time.Second,
	ReadTimeout:         30 * time.Second,
	WriteTimeout:        30 * time.Second,
}

func main() {
	app := fastic.New()
	// create .env file for proxy server settings
	app.Env.Create()

	// load config file with configLoader helper
	cnf := configLoader(app)

	lg := logger.New(app, cnf)

	// path for get proxy and other services status
	// response json with restful-api
	app.Get("/health", func(c *fastic.Ctx) {
		helpers.Json(c, helpers.Options{
			Status:  status.Ready,
			Message: "proxy is ready to use",
			Code:    response.Healthly,
			Data:    helpers.ServicePinger(app, cnf, httpClient),
		})

	}, middleware.Limiter(cnf, 1*time.Minute, 2, middleware.LimiterCustomResp{
		Message: "to many requests, two requests per minute",
	}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // cancel after ending server

	// create user-service client and register all paths of user-service,
	// with error handeling and print erorrs on log/terminal (dont need handle error here.)
	user.Register(app, lg, ctx, cnf.Services["user"], httpClient, nil) // used nil for middlewares, because i dont need add some middlewares now.

	// run proxy
	app.Run(app.Handler)
}

// helpers and more...

// config loader helper, maked for load config file and handle errors
func configLoader(app *fastic.App) *config.Config {
	cnfPath := "../config.jsonc"          // path of config file
	err, cnf := config.Load(app, cnfPath) // load request
	if err != nil {                       // handle errors
		// have error
		if !app.Env.DevelopemtMode {
			// write as log and exit from server
			log.Fatalf("cannot load config file, exiting from app.\n▸ turn on development mode for ignore this")
		}
		// write log on console
		helpers.Print(app, color.Red, "CNF-ERROR", fmt.Sprintf("cannot load config file: %v", err))

		// return empty config variable.
		return cnf
	}
	// return default config variable
	return cnf
}

// do it:
// 1) fix fake pinger on /health path [DONE]
// 2) create mini proxy handler for handle services requests [DONE]
// 3) fix config file problem (with config file value i cannot connect to user service) [DONE]
// 4) define user service (at line 55) and other services in custom function and... [DONE]
// 5) if target service is offline show custom message for this [DONE]
// 6) add logger (with log file) [DONE]
// 7) add auto-backuper for logging files [DONE]
// 8) create auth middlewares
// 9) create custom middleware for limit requests for see ping of services (10 request per minute) [DONE]
// 10) adding alert function for send notification to admin/user (send email or system notification) for send emergency error's to user (after adding check ping.go:43)
// 11) adding test's for logger [DONE]
