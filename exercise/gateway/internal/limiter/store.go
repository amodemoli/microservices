package limiter

import (
	"sync"
	"time"

	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"golang.org/x/time/rate"
)

// Store hold's all per-client limiter's in a one map
// this struct have mutex, rate, brust and more options for limiter...
type Store struct {
	// config model, for reading on every sections need it
	// e.g: on cleanup method need it to read ticker time in minute
	cnf *config.Config

	// can write and read here
	mu *sync.RWMutex

	// holds all per-client limiters here,
	// key is string (ip) and value is user client struct model!
	limiters map[string]*Client

	// token's per second
	limit rate.Limit
	// maximum packet size
	brust int
}

// NewStore function make's new store model ready to use,
// this function get's rate.Limit and brust next return's Store struct model.
func NewStore(cnf *config.Config, period time.Duration, r rate.Limit, b int) *Store {

	store := &Store{
		// config model for use it in cleanup method and get default cleanup ticker time
		// in minutes =D
		cnf: cnf,

		mu: &sync.RWMutex{},

		// create map for limiters
		limiters: make(map[string]*Client),

		// rate.Limit (token's per second)
		limit: r,
		// brust, maximum packet size
		brust: b,
	}

	// run cleanup method on goroutine
	// for run this method every minutes in background
	// without calling
	go store.cleanup(period)

	// return finally store model
	return store
}

// this method, maded for cleanup limiter's from map
func (store *Store) cleanup(period time.Duration) {

	// checking period time, because period time cannot smaller than 5 seconds
	if period < 5*time.Second {
		// 30 second is standard period time for rate limiter
		period = 30 * time.Second
	}

	// make new ticker, for clean'up limiter lists every tick
	// i got this time to config file
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		// lock with mutex
		store.mu.Lock()
		now := time.Now()

		for key, cl := range store.limiters {
			if now.Sub(cl.LastSeen) > period {
				delete(store.limiters, key)
			}
		}
		store.mu.Unlock()
	}
}

// getLimiter method update's user last seen on limiter or add's new user to map
func (store *Store) GetLimiter(key string) *rate.Limiter {

	store.mu.RLock()
	// search user into limiters map
	cl, exists := store.limiters[key]
	store.mu.RUnlock()

	if exists {
		// user exists into limiter's map
		// update lastseen and return new limiter map
		cl.LastSeen = time.Now()
		return cl.Limiter
	}

	store.mu.Lock()

	// unlock after this method
	defer store.mu.Unlock()

	cl, exists = store.limiters[key]
	if !exists {
		// make new client, client is not exists need to make new
		cl = &Client{
			Limiter:  rate.NewLimiter(store.limit, store.brust),
			LastSeen: time.Now(),
		}

		// save it client into map with key
		store.limiters[key] = cl
	}

	// reutrn client limiter
	return cl.Limiter
}
