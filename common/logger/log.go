package logger

import (
	"github.com/inconshreveable/log15"
	"github.com/inconshreveable/log15/ext"
	"github.com/spf13/viper"
	"strings"
)

const (
	CallerStack = "WithCallerStack"
)

func InitLogHandle(l log15.Handler) log15.Handler {
	level := convertLogLevel(viper.GetString("log_level"))
	h := log15.MultiHandler(
		log15.FilterHandler(func(r *log15.Record) bool {
			return r.Lvl <= log15.LvlError || containCallerStackParam(r.Ctx)
		}, log15.CallerStackHandler("%+v", l)),
		log15.FilterHandler(func(r *log15.Record) bool {
			return r.Lvl >= log15.LvlWarn
		}, l),
	)
	h = ext.FatalHandler(h)
	h = log15.LvlFilterHandler(level, h)
	return h
}

func WithCallerStack(l log15.Logger) log15.Logger {
	return l.New(log15.Ctx{CallerStack: true})
}

func containCallerStackParam(ctx []interface{}) bool {
	length := len(ctx)
	if length == 0 || length%2 != 0 {
		return false
	}
	for i := 0; i < length; i += 2 {
		key, ok := ctx[i].(string)
		if !ok {
			continue
		}
		if strings.EqualFold(key, CallerStack) && ctx[i+1].(bool) {
			return true
		}
	}
	return false
}

func convertLogLevel(lvl string) log15.Lvl {
	switch lvl {
	case "DEBUG":
		return log15.LvlDebug
	case "INFO":
		return log15.LvlInfo
	case "WARN":
		return log15.LvlWarn
	case "ERROR":
		return log15.LvlError
	}
	return log15.LvlDebug
}
