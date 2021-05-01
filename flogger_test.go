package flogger_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/axkit/flogger"
	"github.com/rs/zerolog"
)

var ts = []struct {
	params         []interface{}
	expectedParams string
}{
	{
		[]interface{}{"age"},
		"[age]",
	},
	{
		[]interface{}{"age", 41},
		"[age 41]",
	},
	{
		[]interface{}{"age", 42, struct {
			Name    string
			Balance int64
		}{"John", 100}},
		"[age 42 {Name:John Balance:100}]",
	},
}

func TestLogger_Enter(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := range ts {
		flog.Enter(ts[i].params...)
		if buf.Len() == 0 {
			t.Errorf("#1 case %d failed, expected %s, got nothing", i, ts[i].expectedParams)
		}

		m := map[string]interface{}{}
		err := json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if m["params"] != ts[i].expectedParams {
			t.Errorf("#2 case %d failed, expected %s, got %s", i, ts[i].expectedParams, m["params"])
		}
		if m["repo"] != "test" {
			t.Errorf("#3 case %d failed, expected %s, got %s", i, "test", m["repo"])
		}
		if m["func"] != "TestLogger_Enter" {
			t.Errorf("#4 case %d failed, expected %s, got %s", i, "TestLogger_Enter", m["func"])
		}
		if m["message"] != "enter" {
			t.Errorf("#5 case %d failed, expected %s, got %s", i, "enter", m["message"])
		}
		buf.Reset()
	}

	flog.Enter()
	if buf.Len() == 0 {
		t.Error("#6 case failed, expected output, got nothing")
	}
}
func TestLogger_EnterSilent(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		lg.EnterSilent()
		if buf.Len() > 0 {
			t.Errorf("#1 case %d failed, expected no output, got %s", i, buf.String())
		}
	}
}

func TestLogger_EnterSilentExit(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		fc := lg.EnterSilent()
		fc.Exit(ts[i].params...)

		m := map[string]interface{}{}
		err := json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if m["params"] != ts[i].expectedParams {
			t.Errorf("#1 case %d failed, expected %s, got %s", i, ts[i].expectedParams, m["params"])
		}
		buf.Reset()
	}
}

func TestLogger_EnterSilentOnExitExit(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		fc := lg.EnterSilent()

		fc.OnExit(ts[i].params...)
		if buf.Len() > 0 {
			t.Errorf("#1 case %d failed, expected no output, got %s", i, buf.String())
		}

		fc.Exit()
		if buf.Len() == 0 {
			t.Errorf("#2 case %d failed, expected output, got nothing", i)
		}

		m := map[string]interface{}{}
		err := json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if m["params"] != ts[i].expectedParams {
			t.Errorf("#2 case %d failed, expected %s, got %s", i, ts[i].expectedParams, m["params"])
		}
		buf.Reset()
	}
}
func TestLogger_ZerologDebugLevel(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf).Level(zerolog.InfoLevel)
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		fc := lg.Enter(ts[i].params...)
		if buf.Len() > 0 {
			t.Errorf("#1 case %d failed, expected no output, got %s", i, buf.String())
		}

		fc.Exit()
		if buf.Len() > 0 {
			t.Errorf("#2 case %d failed, expected no output, got %s", i, buf.String())
		}
		buf.Reset()
	}
}
func TestLogger_SetSecondExitHandler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")
	x := 0
	flog.SetSecondExitHandler(func(fc *flogger.FuncCall) {
		x++
	})

	for i := range ts {
		fc := flog.Enter()
		if buf.Len() == 0 {
			t.Errorf("#1 case %d failed, expected output, got nothing", i)
		}
		fc.Exit()
		if x != i+1 {
			t.Errorf("#2 case %d failed, expected %d, got %d", i, i+1, x)
		}

		buf.Reset()
	}
}

func BenchmarkLogger_Enter1(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		flog.Enter(10)
		buf.Reset()
	}
}

func BenchmarkLogger_Enter3(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		fc := flog.Enter(10, "John", 42)
		_ = fc
		buf.Reset()
	}
}

func BenchmarkLogger_EnterSilentExit1(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		flog.EnterSilent().Exit(10)
		buf.Reset()
	}
}
func BenchmarkLogger_EnterSilentExit3(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		flog.EnterSilent().Exit(10, "John", 42)
		buf.Reset()
	}
}
func BenchmarkLogger_EnterExit3(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		flog.Enter(10, "John", 42).Exit()
		buf.Reset()
	}
}

func BenchmarkBuffer_WriteReset(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	_ = flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		buf.Write([]byte(`[age 42]`))
		buf.Reset()
	}
}

func BenchmarkSprintf(b *testing.B) {

	p := []interface{}{"age", 42}
	for i := 0; i < b.N; i++ {
		s := fmt.Sprintf("%+v", p...)
		_ = s
	}
}

func BenchmarkZerolog(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	at := time.Now()

	for i := 0; i < b.N; i++ {
		zl.Debug().Str("func", "TestFlogger_EnterSilent").Str("dur", time.Since(at).String()).Msg("exit")
		buf.Reset()
	}
}
