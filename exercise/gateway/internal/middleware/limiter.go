package middleware

import (
	"time"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
	"github.com/amodemoli/microservices/exercise/gateway/internal/limiter"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

// custom response struct because if user need to,
// send custom json response after limit
// *Cannot Change "Status" && "Response Code" fields!
type LimiterCustomResp struct {
	// message field on json responses
	Message string
	// data field on json repsonses, can empty
	Data interface{}
}

// Limiter function is a middleware, this function maked for limit user requests
func Limiter(cnf *config.Config, period time.Duration, limit int, customResp ...LimiterCustomResp) func(next fastic.RequestHandler) fastic.RequestHandler {

	var data interface{}
	// message has default value (because is requiered)
	message := "to many requests, try again later"

	if len(customResp) != 0 {
		if customResp[0].Data != nil {
			data = customResp[0].Data
		}
		if customResp[0].Message != "" {
			message = customResp[0].Message
		}
	}

	ratePerSec := rate.Limit(float64(limit) / period.Seconds())

	store := limiter.NewStore(cnf, period, ratePerSec, limit)

	// this section run's only one time for any user
	return func(next fastic.RequestHandler) fastic.RequestHandler {
		return func(c *fastic.Ctx) {
			// this section run's per request

			// get user remote-ip and change it to string
			key := c.RemoteIP().String()

			// make new rateLimiter model
			rateLimiter := store.GetLimiter(key)

			if !rateLimiter.Allow() {
				c.Status(fasthttp.StatusTooManyRequests)
				helpers.Json(c, helpers.Options{
					Status:  status.RateLimited,
					Message: message,
					Code:    response.TooManyRequests,
					Data:    data,
				})
				return
			}

			// operation is done, call request.
			next(c)
		}
	}
}
