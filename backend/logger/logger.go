package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// InitLogger initializes the logger with proper configuration
func InitLogger() {
	Log = logrus.New()

	// Set log level
	Log.SetLevel(logrus.InfoLevel)

	// Set JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Configure log rotation
	Log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})

	// Also log to stdout in development
	if gin.Mode() == gin.DebugMode {
		Log.SetOutput(os.Stdout)
	}
}

// GinLogger returns a gin middleware for logging
func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		Log.WithFields(logrus.Fields{
			"status":     param.StatusCode,
			"method":     param.Method,
			"path":       param.Path,
			"ip":         param.ClientIP,
			"user-agent": param.Request.UserAgent(),
			"latency":    param.Latency,
			"error":      param.ErrorMessage,
		}).Info("Request completed")
		return ""
	})
}

// Recovery middleware with logging
func GinRecovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(Log.Out, func(c *gin.Context, err interface{}) {
		Log.WithFields(logrus.Fields{
			"error":  err,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Error("Panic recovered")
	})
}

// LogError logs an error with context
func LogError(err error, context string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	fields["context"] = context
	Log.WithFields(fields).Error(err)
}

// LogInfo logs an info message with context
func LogInfo(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	Log.WithFields(fields).Info(message)
}

// LogWarning logs a warning message with context
func LogWarning(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	Log.WithFields(fields).Warning(message)
}

// LogDebug logs a debug message with context
func LogDebug(message string, fields logrus.Fields) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	Log.WithFields(fields).Debug(message)
}
