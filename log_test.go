package logrus

import (
	"bytes"
	"context"
	"encoding/json"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-kita/log"
)

func TestLogrusLog(t *testing.T) {
	buf := &bytes.Buffer{}
	ll := logrus.New()
	ll.SetLevel(logrus.TraceLevel)
	ll.SetOutput(buf)
	ll.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	ll.SetReportCaller(true)

	log.GetLevelStore().Set("", log.DebugLevel)
	log.GetLevelStore().Set("closed", log.ClosedLevel)
	provider := MakeProvider(ll)
	logger := provider("")
	for i := log.DebugLevel; i <= log.ErrorLevel; i++ {
		logger.AtLevel(context.Background(), i).
			Print("abc")
		got := buf.String()
		if !strings.Contains(got, "_test.go:") {
			t.Errorf("expect %s outPutter contains current test file name, but got: %q",
				i, got)
		}
		buf.Reset()
	}
	closed := provider("closed")
	closed.AtLevel(context.Background(), log.InfoLevel).Print("abc")
	got := buf.String()
	buf.Reset()
	if got != "" {
		t.Errorf("expect closed outPutter nothing, but got %q", got)
	}
}

func TestFormatter_Format(t *testing.T) {
	buf := &bytes.Buffer{}
	ll := logrus.New()
	ll.SetLevel(logrus.TraceLevel)
	ll.SetOutput(buf)
	ll.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	ll.SetReportCaller(false)
	op := NewOutPutter(ll,
		LevelFuncOption(func(level log.Level) logrus.Level {
			return logrus.InfoLevel
		}),
		NowFuncOption(func() time.Time {
			return time.Date(2021, time.June, 21, 22, 50, 20, 0, time.Local)
		})).(*outPutter)
	op.callerFunc = func(callDepth int) *runtime.Frame {
		return defaultCallerFunc(callDepth)
	}
	op.OutPut(context.Background(), "abc", log.ErrorLevel, "error?", []log.Field{
		{Key: "", Value: "..."},
		{Key: "xyz", Value: "XYZ"},
	}, 3)
	got := buf.String()
	t.Logf("got %s", got)
	mp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(got), &mp); err != nil {
		t.Fatal(err)
	}
	level := mp["level"].(string)
	if strings.ToLower(level) != "info" {
		t.Errorf("expect level info ,but got %q", got)
	}
	if _, ok := mp["caller"]; ok {
		t.Errorf("expect no caller in output, but contains: %q", got)
	}
	wantTime := "2021-06-21T22:50:20"
	if !strings.Contains(mp["time"].(string), wantTime) {
		t.Errorf("expect time: %q, but got: %q", wantTime, got)
	}
}
