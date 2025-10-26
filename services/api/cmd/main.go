package main

import "lib/utils/logging"

func main() {
	logger := logging.GetLogger()
	logger.Info("Logger initialized successfully")
}
