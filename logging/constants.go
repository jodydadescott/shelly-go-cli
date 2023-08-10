package logging

import (
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogLevel           = zapcore.InfoLevel
	defaultSamplingInitial    = 100
	defaultSamplingThereafter = 100
	defaultEncoding           = "console"
)
