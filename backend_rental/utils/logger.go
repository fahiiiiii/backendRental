// utils/logger.go
package utils

import (
    "log"
    "os"
)

var (
    // Logger is the global logger instance
    Logger *log.Logger
)

// InitLogger initializes the logger with specified settings
func InitLogger() {
    // Create a log file
    logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }

    // Set up the logger to write to both standard error and the log file
    Logger = log.New(logFile, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

    // Optionally, log to standard error as well
    log.SetOutput(os.Stderr)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// LogInfo logs informational messages
func LogInfo(message string) {
    Logger.Println("INFO: " + message)
}

// LogError logs error messages
func LogError(err error) {
    Logger.Println("ERROR: " + err.Error())
}

// LogFatal logs fatal error messages and exits the application
func LogFatal(err error) {
    Logger.Fatalln("FATAL: " + err.Error())
}