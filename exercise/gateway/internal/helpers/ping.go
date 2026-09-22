package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/valyala/fasthttp"
)

// Ping function send's request to your target url,
// this function record's target service response delay and show this delay to you as (ms) - milisecond
// if your target service is offline (function get's timeout error), this function show's (-1) to you this response is milisecond
func Ping(app *fastic.App, cnf *config.Config, client *fasthttp.Client, target string, customTimeout ...time.Duration) int64 {

	// get final timeout from finalTimeout helper,
	// this function have default value for timeout, first get's timeout from customTimeout
	// and if customTimeout is empty or cannot validate it try's to get new timeout from config
	timeout := finalTimeout(app, cnf, customTimeout[0])

	// get request and response from fasthttp
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	defer func() {
		// release request and response after function work
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	// validate target host url, if is empty return -1 with
	// no real pinging (dont need to ping because they know this addres is not valid and ping result is -1 dont need to real pinging)
	if strings.TrimSpace(target) == "" {
		return -1
	}

	// validate is done, now change's target server to req variable,
	req.SetRequestURI(target)

	// changed header to "HEAD", because HEAD is lighter than GET
	req.Header.SetMethod(fasthttp.MethodHead)

	// get now time (start pining server, counting response delay)
	start := time.Now()

	// send request with timeout, used client struct for this!
	// this line is time-out of target url you can change it on config.jsonc file)
	err := client.DoTimeout(req, resp, timeout)

	// count response delay and save it to elapsed variable
	elapsed := time.Since(start)

	// handle errors, because error mean's target host does not response
	if err != nil {
		return -1
	}

	// return response delay with millisecond
	return elapsed.Milliseconds()
}

// validateTimeout function maded for show to you your timeout seconds is good or not!
// minumum value for timeout is 2 seconds and maximum is 100 seconds!
func validateTimeout(d time.Duration) bool {
	return d >= 2*time.Second && d <= 100*time.Second
}

// finalTimeout get's app && config model and customTimeout, and return's finally timeout Duration
// only maded for use on Ping
func finalTimeout(app *fastic.App, cnf *config.Config, customTimeout time.Duration) time.Duration {
	// make variable named timeout and save default value to 10 second
	var timeout time.Duration = 10 * time.Second

	if validateTimeout(customTimeout) {
		// save it user custom time-out
		timeout = customTimeout
	} else {
		// get timeout from config file and change it to time.Duration.
		cnfTimeout, err := toDuration(cnf.General["pinging_timeout"])
		// if error is empty, replace it
		if err == nil {
			if validateTimeout(cnfTimeout) {
				// save it config file timeout
				timeout = cnfTimeout
			}
		}
	}

	// return finally timeout
	return timeout
}

// anyToTime converts any type to time.Duration type
// if cannot convert result is 0
func toDuration(v any) (time.Duration, error) {

	switch val := v.(type) {
	case time.Duration:
		return val, nil
	case float64: // JSON numbers come as float64
		return time.Duration(val * float64(time.Second)), nil
	case int:
		return time.Duration(val) * time.Second, nil
	case int64:
		return time.Duration(val) * time.Second, nil
	case string:
		return time.ParseDuration(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to time.Duration", v)
	}
}

// ServicePinger helper, ping's all registered services with Ping helper function's help
// and show's response as Map
func ServicePinger(app *fastic.App, cnf *config.Config, client *fasthttp.Client, customTimeout ...time.Duration) map[string]string {

	// make, variable named timeout and setting default value to -1
	// why -1? because -1 means i dont need to use default timeout (on validate default timeout return is false)
	var timeout time.Duration = -1

	// get length of customTimeout, dont need validate here! time validates on Ping function
	if len(customTimeout) > 0 {
		// replace it to custom timeout
		timeout = customTimeout[0]
	}

	// make new variable named response for save it services ping result here
	response := make(map[string]string, len(cnf.Services))

	// get range of registered services with service url for ping
	for service, url := range cnf.Services {
		// save it service to new response struct and save service response delay (ping)
		response[service+"_service"] = fmt.Sprintf("%vms", Ping(app, cnf, client, url, timeout))
	}

	// return response after write all service's response delay
	return response
}
