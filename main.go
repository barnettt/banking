package main

import (
	"github.com/barnettt/banking/app"
	"github.com/barnettt/banking/logger"
)

func main() {
	logger.Info("Starting App.....")
	app.StartApp()
}
