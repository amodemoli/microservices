package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/valyala/fasthttp"
)

// Ping function send's request to your target url,
// this function record's target service response delay and show this delay to you as (ms) - milisecond
// if your target service is offline (function get's timeout error), this function show's (-1) to you this response is milisecond
func Ping(cnf *config.Config, client *fasthttp.Client, target string) int64 {

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
	// this line is time-out of target url you can change it on config.jsonc file
	err := client.DoTimeout(req, resp, cnf.PingingTimeout*time.Second)

	// count response delay and save it to elapsed variable
	elapsed := time.Since(start)

	// handle errors, because error mean's target host does not response
	if err != nil {
		return -1
	}

	// return response delay with millisecond
	return elapsed.Milliseconds()
}

// ServicePinger helper, ping's all registered services with Ping helper function's help
// and show's response as Map
func ServicePinger(cnf *config.Config, client *fasthttp.Client) map[string]string {

	// make new variable named response for save it services ping result here
	response := make(map[string]string, len(cnf.Services))

	// get range of registered services with service url for ping
	for service, url := range cnf.Services {
		// save it service to new response struct and save service response delay (ping)
		response[service+"_service"] = fmt.Sprintf("%vms", Ping(cnf, client, url))
	}

	// return response after write all service's response delay
	return response
}
