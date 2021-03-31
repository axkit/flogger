package flogger

import (
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

type FunctionCall struct {
	zerolog.Logger
	enteredAt time.Time
	isSilent  bool
}

func New(l *zerolog.Logger, typ, name string) Logger {
	return Logger{
		log: l.With().Str(typ, name).Logger()}
}

func (s *Logger) EnterSilent() *FunctionCall {
	fc := s.enter()
	fc.isSilent = true
	return fc
}

func (s *Logger) Enter(params ...interface{}) *FunctionCall {

	fc := s.enter()

	lp := len(params)
	switch {
	case lp == 1:
		fc.Logger.Debug().Interface("params", params[0]).Msg("enter")
	case lp == 0:
		fc.Logger.Debug().Msg("enter")
	case lp%2 == 0:
		e := fc.Logger.Debug()
		for i := 0; i < lp; i += 2 {
			e = e.Interface(params[i].(string), params[i+1])
		}
		e.Msg("enter")
	case lp%2 == 1:
		e := fc.Logger.Debug()
		for i := 0; i < lp-1; i += 2 {
			e = e.Interface(params[i].(string), params[i+1])
		}
		e.Msgf("%v", params[lp-1])
	default:
		// if more than
	}

	return fc
}

func (s *Logger) enter() *FunctionCall {
	var arr [15]uintptr
	pc := arr[:]
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	//fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)

	dpos := strings.LastIndex(frame.Function, ".")
	if dpos > 0 {
		frame.Function = frame.Function[dpos+1:] + "()"
	}

	return &FunctionCall{enteredAt: time.Now(), Logger: s.log.With().Str("func", frame.Function).Logger()}
}

func (fc *FunctionCall) Exit() {
	if !fc.isSilent {
		fc.Logger.Debug().Str("dur", time.Since(fc.enteredAt).String()).Msg("exit")
	}
}

func (fc *FunctionCall) Error() *zerolog.Event {
	return fc.Logger.Error()
}

func (fc *FunctionCall) Debug() *zerolog.Event {
	return fc.Logger.Debug()
}

func (fc *FunctionCall) Warn() *zerolog.Event {
	return fc.Logger.Warn()
}

func (fc *FunctionCall) Info() *zerolog.Event {
	return fc.Logger.Info()
}
