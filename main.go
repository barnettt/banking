package main

import (
	"github.com/barnettt/banking-lib/logger"
	"github.com/barnettt/banking/app"
)

func main() {
	logger.Info("Starting App.....")
	app.StartApp()
}
