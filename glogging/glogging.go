// logger customs go.uber.org/zap logger with level, caller context and rorate file via github.com/natefinch/lumberjack.
package glogging

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Loglevel string

const (
	DEBUG Loglevel = "debug"
	WARN  Loglevel = "warn"
	ERROR Loglevel = "error"
	FATAL Loglevel = "fatal"
)

var Logger *zap.Logger
var Sugared *zap.SugaredLogger

// Init customs go.uber.org/zap glogging.
//   parameter logpaths available value: stdout,stderr,path/to/file;
//   parameter loglevel
func Init(logpaths []string, loglevel Loglevel) {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	cores := make([]zapcore.Core, 0)
	if len(logpaths) > 0 {
		for _, logto := range logpaths {
			if logto == "stdout" {
				cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stdout), LoglevelStr2uint(loglevel)))
			} else if logto == "stderr" {
				cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stderr), LoglevelStr2uint(loglevel)))
			} else if logto != "" {
				lumberJackLogger := &lumberjack.Logger{
					Filename:   logto,
					MaxSize:    100,
					MaxBackups: 15,
					MaxAge:     30,
					Compress:   true,
				}
				cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(lumberJackLogger), LoglevelStr2uint(loglevel)))
			}
		}
	} else {
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(ioutil.Discard), LoglevelStr2uint(loglevel)))
	}

	Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	defer Logger.Sync()

	Sugared = Logger.Sugar()
}

// LoglevelStr2uint converts logging level from a string into zapcore.Level.
func LoglevelStr2uint(loglevel Loglevel) (rs zapcore.Level) {
	switch loglevel {
	case DEBUG:
		{
			return zapcore.DebugLevel
		}
	case WARN:
		{
			return zapcore.WarnLevel
		}
	case ERROR:
		{

			return zapcore.ErrorLevel
		}
	case FATAL:
		{
			return zapcore.FatalLevel
		}
	}
	return zapcore.DebugLevel
}

// String2LogLevel convert level name string into glogging.Loglevel type
func String2LogLevel(s string) Loglevel {
	level := Loglevel(strings.ToLower(s))
	if level != DEBUG &&
		level != WARN &&
		level != ERROR &&
		level != FATAL {
		level = DEBUG
	}
	return level
}
