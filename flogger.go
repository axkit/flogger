// The flogger package provides a simple function call logging by wrapping github.com/rs/zerolog package.
package flogger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// FuncLogger wraps zerolog logger with type/name pair.
type FuncLogger struct {
	log               zerolog.Logger
	customExitHandler func(*FuncCall)
}

// G holds global default logger. Use flooger.G to use features without custom FuncLogger.
var G = FuncLogger{log: log.Logger}

// FuncCall holds data needed for a sigle function call tracking and logging.
type FuncCall struct {
	*FuncLogger
	// Name holds function name.
	Name string
	// EnteredAt time captured by Enter or EnterSilent.
	EnteredAt time.Time

	// IsSilentEnter holds true if the FunCall build by EnterSilent method.
	IsSilentEnter bool

	// Params hold values to be printed in Exit().
	Params []interface{}

	// Ignore holds true if zerolog.Level not DebugLevel or TraceLevel.
	Ignore bool
}

// New builds new FuncLogger object. Parameters typ and name
// will be placed to the log item.
func New(l *zerolog.Logger, typ, name string) FuncLogger {
	return FuncLogger{log: l.With().Str(typ, name).Logger()}
}

// SetSecondExitHandler assigns secondary handler to be called in Exit().
func (fl *FuncLogger) SetSecondExitHandler(f func(*FuncCall)) *FuncLogger {
	fl.customExitHandler = f
	return fl
}

// Logger returns logger instance.
func (fl *FuncLogger) Logger() zerolog.Logger {
	return fl.log
}

// EnterSilent builds FuncCall object but does not write enter message to the log.
func (fl *FuncLogger) EnterSilent() *FuncCall {
	return fl.enter(true, nil)
}

// Enter builds FuncCall object and writes message "enter" to the log.
// Log item gets func name and all input parameters.
func (fl *FuncLogger) Enter(params ...interface{}) *FuncCall {
	if len(params) == 0 {
		return fl.enter(false)
	}
	return fl.enter(false, params)
}

func (fl *FuncLogger) enter(silent bool, params ...interface{}) *FuncCall {
	if fl.log.GetLevel() > zerolog.DebugLevel {
		return &FuncCall{Ignore: true, FuncLogger: fl}
	}

	var arr [1]uintptr
	pc := arr[:]
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fc := FuncCall{IsSilentEnter: silent, EnteredAt: time.Now(), FuncLogger: fl}
	if dpos := strings.LastIndex(frame.Function, "."); dpos > 0 {
		fc.Name = frame.Function[dpos+1:]
	}

	if silent {
		return &fc
	}

	if len(params) > 0 {
		fc.log.Debug().Str("func", fc.Name).Str("params", fmt.Sprintf("%+v", params...)).Msg("enter")
	} else {
		fc.log.Debug().Str("func", fc.Name).Msg("enter")
	}

	return &fc
}

// OnExit informs silent function call, that parameters should be
// written to the log in the Exit() call.
func (fc *FuncCall) OnExit(params ...interface{}) *FuncCall {
	if len(fc.Params) == 0 {
		fc.Params = make([]interface{}, len(params))
		copy(fc.Params, params)
	}
	return fc
}

// ExitHandler implements default behavior on exit.
//
// If the function started with EnterSilent(), Exit() will be silent as well.
// If was called OnExit() between, then Exit() lost silence.
var ExitHandler = func(fc *FuncCall) {

	if fc.IsSilentEnter && len(fc.Params) == 0 {
		return
	}

	e := fc.log.Debug().Str("func", fc.Name).
		Int64("dur", int64(time.Since(fc.EnteredAt)))

	if len(fc.Params) > 0 {
		e = e.Str("params", fmt.Sprintf("%+v", fc.Params))
	}

	if fc.IsSilentEnter {
		e.Msg("enter/exit")
	} else {
		e.Msg("exit")
	}
	return
}

// Exit writes to the log function execution duration.
// If you dont want have 2 log lines in the log (for enter and exit) you can
// have 1 log line by calling defer flog.EnterSilent().Exit(id, name, age).
func (fc *FuncCall) Exit(params ...interface{}) {
	if fc.FuncLogger.customExitHandler != nil {
		fc.FuncLogger.customExitHandler(fc)
	}

	if fc.Ignore {
		return
	}

	if len(params) > 0 && len(fc.Params) == 0 {
		fc.Params = make([]interface{}, len(params))
		copy(fc.Params, params)
	}

	ExitHandler(fc)
}

// Error return zerolog.Event for writing intermediate log items
// between Enter() and Exit().
func (fc *FuncCall) Error() *zerolog.Event {
	return fc.log.Error().Str("func", fc.Name)
}

// Debug return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Debug() *zerolog.Event {
	return fc.log.Debug().Str("func", fc.Name)
}

// Warn return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Warn() *zerolog.Event {
	return fc.log.Warn().Str("func", fc.Name)
}

// Info return zerolog.Event for writing intermetiate log items
// between Enter() and Exit().
func (fc *FuncCall) Info() *zerolog.Event {
	return fc.log.Info().Str("func", fc.Name)
}
