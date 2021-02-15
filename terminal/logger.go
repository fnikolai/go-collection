package terminal

import (
	log "github.com/sirupsen/logrus"
)

func SetLogger(level string, jsonOutput bool) error {
	//  default configuration
	log.SetLevel(log.WarnLevel)
	log.SetFormatter(&log.TextFormatter{})

	// set debug log level
	switch level := level; level {
	case "debug", "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "info", "INFO":
		log.SetLevel(log.InfoLevel)
	case "warning", "WARNING":
		log.SetLevel(log.WarnLevel)
	case "error", "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "fatal", "FATAL":
		log.SetLevel(log.FatalLevel)
	case "panic", "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	// set log formatter to JSON
	if jsonOutput {
		log.SetFormatter(&log.JSONFormatter{})
	}

	return nil
}
