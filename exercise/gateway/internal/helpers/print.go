package helpers

import (
	"time"

	"github.com/amodemoli/fastic/core/fastic"
)

// on fastic framework they cannot write message on console before run server
// im added custom print function for fix this 
func Print(app *fastic.App, color, model, message string) {
	time.AfterFunc(50*time.Millisecond, func() {
		app.Print(color, model, message)
	})
}