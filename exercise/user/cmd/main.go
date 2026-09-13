package main

import (
	"context"
	"log"
	"os"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/user/internal/handler"
	"github.com/amodemoli/microservices/exercise/user/protobuf"
)

func main() {
	app := fastic.New()

	// create env file
	app.Env.Create()

	ctx := context.Background()
	if err := protobuf.RegisterUserHTTPGoServer(ctx, app.Router, &handler.Handler{}, nil); err != nil {
		log.Fatalf("cannot run user service: %v", err)
		os.Exit(1) // exit from application with code
	}
	log.Printf("server runned at localhost:%d", app.Env.Port)

	app.Run(app.Handler)
}
