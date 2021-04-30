package flogger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// FuncLogger wraps zerolog logger with type/name pair.
type FuncLogger struct {
	log zerolog.Logger
}

// FuncCall holds data needed for a sigle function call tracking and logging.
type FuncCall struct {
	*FuncLogger
	name      string
	enteredAt time.Time
	isSilent  bool
	addParams bool
	paramStr  string
	p         []interface{}
}

// New builds new FuncLogger object.
func New(l *zerolog.Logger, typ, name string) FuncLogger {
	return FuncLogger{log: l.With().Str(typ, name).Logger()}
}

// Logger returns logger instance.
func (s *FuncLogger) Logger() zerolog.Logger {
	return s.log
}

// EnterSilent builds FuncCall object but does not write enter message to the log.
func (s *FuncLogger) EnterSilent(params ...interface{}) *FuncCall {
	return s.enter(true, params)
}

// Enter builds FuncCall object and writes message "enter" to the log.
// Log item gets func name and all input parameters.
func (s *FuncLogger) Enter(params ...interface{}) *FuncCall {
	return s.enter(false, params)
}

func (s *FuncLogger) enter(silent bool, params ...interface{}) *FuncCall {
	var arr [1]uintptr
	pc := arr[:]
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fc := FuncCall{isSilent: silent, enteredAt: time.Now(), FuncLogger: s}
	if dpos := strings.LastIndex(frame.Function, "."); dpos > 0 {
		fc.name = frame.Function[dpos+1:]
	}

	if silent {
		if len(params) > 0 {
			fc.p = make([]interface{}, len(params))
			copy(fc.p, params)
		}
		return &fc
	}

	if len(params) > 0 {
		fc.log.Debug().Str("func", fc.name).Str("params", fmt.Sprintf("%+v", params...)).Msg("enter")
	} else {
		fc.log.Debug().Str("func", fc.name).Msg("enter")
	}

	return &fc
}

// Report informs silent function call, that input parameters should be
// written to the log in the Exit() call.
func (fc *FuncCall) Report() *FuncCall {
	fc.addParams = true
	return fc
}

// Exit writes exit message to the log together with execution duration.
// If silent function call was asked to do flush by calling Flush() before,
// the input parameters will be written as well.
func (fc *FuncCall) Exit() {
	if !fc.isSilent {
		fc.log.Debug().Str("func", fc.name).Str("dur", time.Since(fc.enteredAt).String()).Msg("exit")
		return
	}

	// silent mode
	if !fc.addParams {
		return
	}

	if len(fc.p) > 0 {
		fc.log.Debug().Str("func", fc.name).Str("dur", time.Since(fc.enteredAt).String()).Str("params", fmt.Sprintf("%+v", fc.p...)).Msg("enter/exit")
	} else {
		fc.log.Debug().Str("func", fc.name).Str("dur", time.Since(fc.enteredAt).String()).Msg("enter/exit")
	}
}

// Error return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Error() *zerolog.Event {
	return fc.log.Error().Str("func", fc.name)
}

// Debug return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Debug() *zerolog.Event {
	return fc.log.Debug().Str("func", fc.name)
}

// Warn return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Warn() *zerolog.Event {
	return fc.log.Warn().Str("func", fc.name)
}

// Info return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Info() *zerolog.Event {
	return fc.log.Info().Str("func", fc.name)
}
