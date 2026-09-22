package limiter

import (
	"time"

	"golang.org/x/time/rate"
)

// Client struct save's Limiter struct and time of lastSeen of user
// they save it new client before send request to server
type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}
