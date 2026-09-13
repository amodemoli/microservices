package helpers

import (
	"errors"
	"time"

	"github.com/amodemoli/fastic/core/fastic"
	resp "github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/codes/status"
)

// response struct model for send it easy
type response struct {
	Status    status.Status `json:"status"`
	Message   string        `json:"message,omitempty"`
	Code      resp.Code     `json:"code,omitempty"`
	Data      interface{}   `json:"data,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// options of Json function, maded for get json request options
type Options struct {
	Status  status.Status // status of response | e.g: OK
	Message string        // content message of response | e.g: user logged in
	Code    resp.Code     // code of response | e.g: code.Ok
	Data    interface{}   // data of response | you can add other interface/map/struct or new fields for response (no limit)
}

// Json function, maded for send RestFulApi requests to target
func Json(c *fastic.Ctx, opts Options) error {
	// validate required fields, check length
	if opts.Status == "" || opts.Message == "" || opts.Code == "" {
		return errors.New("please fill out all the required fields")
	}

	// send json response after validate
	return c.JSON(response{
		Status:  opts.Status,
		Message: opts.Message,
		Code:    opts.Code,
		Data:    opts.Data,
	})
}
