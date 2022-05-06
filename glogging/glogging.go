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

var Logger *zap.Logger
var Sugared *zap.SugaredLogger

// Init customs go.uber.org/zap glogging.
//   parameter logpaths avaiable value: stdout,stderr,path/to/file;
//   parameter loglevel avaiable value: debug,warn,error,fatal.
func Init(logpaths []string, loglevel string) {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	cores := make([]zapcore.Core, 0)
	if len(logpaths) > 0 {
		for _, logto := range logpaths {
			if logto == "stdout" {
				cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stdout), loglevelStr2uint(loglevel)))
			} else if logto == "stderr" {
				cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stderr), loglevelStr2uint(loglevel)))
			} else if logto != "" {
				lumberJackLogger := &lumberjack.Logger{
					Filename:   logto,
					MaxSize:    100,
					MaxBackups: 15,
					MaxAge:     30,
					Compress:   true,
				}
				cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(lumberJackLogger), loglevelStr2uint(loglevel)))
			}
		}
	} else {
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(ioutil.Discard), loglevelStr2uint(loglevel)))
	}

	Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	defer Logger.Sync()

	Sugared = Logger.Sugar()
}

// loglevelStr2uint converts logging level from a string into zapcore.Level.
func loglevelStr2uint(loglevel string) (rs zapcore.Level) {
	loglevel = strings.ToLower(loglevel)

	switch loglevel {
	case "debug":
		{
			return zapcore.DebugLevel
		}
	case "warn":
		{
			return zapcore.WarnLevel
		}
	case "error":
		{

			return zapcore.ErrorLevel
		}
	case "fatal":
		{
			return zapcore.FatalLevel
		}
	}
	return zapcore.DebugLevel
}
