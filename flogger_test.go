package flogger_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		fc := lg.Enter(ts[i].params...)
		m := map[string]interface{}{}
		err := json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if m["params"] != ts[i].expectedParams {
			t.Errorf("#1 case %d failed, expected %s, got %s", i, ts[i].expectedParams, m["params"])
		}
		if m["repo"] != "test" {
			t.Errorf("#2 case %d failed, expected %s, got %s", i, "test", m["repo"])
		}
		if m["func"] != "TestLogger_Enter" {
			t.Errorf("#3 case %d failed, expected %s, got %s", i, "TestLogger_Enter", m["func"])
		}

		if m["message"] != "enter" {
			t.Errorf("#3 case %d failed, expected %s, got %s", i, "enter", m["message"])
		}

		buf.Reset()
		fc.Exit()
		if buf.Len() == 0 {
			t.Errorf("#4 case %d failed, expected log item, got nothing", i)
		}
		buf.Reset()
	}
}
func TestLogger_EnterSilent(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	lg := flogger.New(&zl, "repo", "test")

	for i := range ts {
		flog := lg.EnterSilent(ts[i].params...)
		if buf.Len() > 0 {
			t.Errorf("#1 case %d failed, expected no output, got %s", i, buf.String())
		}
		flog.Report().Exit()
		if buf.Len() == 0 {
			t.Errorf("#2 case %d failed, expected log item, got nothing", i)
		}

		m := map[string]interface{}{}
		err := json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if m["params"] != ts[i].expectedParams {
			t.Errorf("#3 case %d failed, expected %s, got %s", i, ts[i].expectedParams, m["params"])
		}
		buf.Reset()
	}

}

func TestLogger_Report(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(os.Stderr)
	flog := flogger.New(&zl, "repo", "test")

	fc := flog.EnterSilent("age", 12)
	if buf.Len() > 0 {
		t.Errorf("#1 failed, expected nothing got log item")
	}
	fc.Exit()
}

func BenchmarkLogger_EnterSilent(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		fc := flog.EnterSilent("age", 42)
		fc.Exit()
		buf.Reset()
	}
}

func BenchmarkLogger_EnterSilentEmpty(b *testing.B) {

	buf := bytes.NewBuffer(nil)
	zl := zerolog.New(buf)
	flog := flogger.New(&zl, "repo", "test")

	for i := 0; i < b.N; i++ {
		fc := flog.EnterSilent()
		fc.Exit()
		buf.Reset()
	}
}

func BenchmarkBuffer_Reset(b *testing.B) {

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
	//fc := flogger.New(&zl, "repo", "test")
	at := time.Now()

	for i := 0; i < b.N; i++ {
		zl.Debug().Str("func", "TestFlogger_EnterSilent").Str("dur", time.Since(at).String()).Msg("exit")
		buf.Reset()
	}
}
