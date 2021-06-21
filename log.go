package logrus

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-kita/log"
)

type callDepthKey struct {
}

type formatter struct {
	op        *outPutter
	formatter logrus.Formatter
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if !entry.HasCaller() {
		return f.formatter.Format(entry)
	}
	callDepth, ok := entry.Context.Value(callDepthKey{}).(int)
	if !ok {
		return f.formatter.Format(entry)
	}
	entry.Caller = f.op.callerFunc(callDepth + 8)
	return f.formatter.Format(entry)
}

type OutPutterOption func(o *outPutter)

func NowFuncOption(fn func() time.Time) OutPutterOption {
	return func(o *outPutter) {
		o.nowFunc = fn
	}
}

func LevelFuncOption(fn func(level log.Level) logrus.Level) OutPutterOption {
	return func(o *outPutter) {
		o.levelFunc = fn
	}
}

// outPutter is a log.OutPutter based on Zap logger.
type outPutter struct {
	out        *logrus.Logger
	entryPool  *sync.Pool
	nowFunc    func() time.Time
	levelFunc  func(level log.Level) logrus.Level
	callerFunc func(callDepth int) *runtime.Frame
}

// NewOutPutter produce a log.OutPutter based on logrus.Logger
func NewOutPutter(out *logrus.Logger, opts ...OutPutterOption) log.OutPutter {
	out = &logrus.Logger{
		Out:          out.Out,
		Hooks:        out.Hooks,
		Formatter:    out.Formatter,
		ReportCaller: out.ReportCaller,
		Level:        out.Level,
		ExitFunc:     out.ExitFunc,
	}
	op := &outPutter{
		out: out,
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &logrus.Entry{}
			},
		},
		nowFunc:    time.Now,
		levelFunc:  defaultLevelFunc,
		callerFunc: defaultCallerFunc,
	}
	for _, opt := range opts {
		opt(op)
	}
	out.Formatter = &formatter{op, out.Formatter}
	return op
}

func (o *outPutter) OutPut(
	ctx context.Context, name string, level log.Level, msg string, fields []log.Field, callDepth int) {
	ctx = context.WithValue(ctx, callDepthKey{}, callDepth)
	entry := o.entryPool.Get().(*logrus.Entry)
	defer func() {
		entry.Data = map[string]interface{}{}
		o.entryPool.Put(entry)
	}()
	entry.Logger = o.out
	entry.Data = make(map[string]interface{}, len(fields)+1)
	entry.Data[log.LoggerKey] = name
	for _, field := range fields {
		if len(field.Key) == 0 {
			continue
		}
		if field.Key == log.LevelKey {
			continue
		}
		entry.Data[field.Key] = log.Value(ctx, field.Value)
	}
	entry.Time = o.nowFunc()
	entry.Context = ctx
	entry.Log(o.levelFunc(level), msg)
}

func defaultLevelFunc(level log.Level) logrus.Level {
	switch {
	case level == log.InfoLevel:
		return logrus.InfoLevel
	case level == log.WarnLevel:
		return logrus.WarnLevel
	case level == log.ErrorLevel:
		return logrus.ErrorLevel
	case level == log.DebugLevel:
		return logrus.DebugLevel
	case level < log.DebugLevel:
		return logrus.TraceLevel
	default:
		return logrus.ErrorLevel
	}
}

func defaultCallerFunc(callDepth int) *runtime.Frame {
	pcs := make([]uintptr, 1)
	_ = runtime.Callers(callDepth, pcs)
	frames := runtime.CallersFrames(pcs)
	frame, _ := frames.Next()
	return &frame
}

// NewLogger create a log.Logger of name based on logrus.Logger.
func NewLogger(name string, out *logrus.Logger) log.Logger {
	return log.NewStdLogger(name, NewOutPutter(out))
}

// MakeProvider make a log.LoggerProvider function.
func MakeProvider(out *logrus.Logger) log.LoggerProvider {
	return func(name string) log.Logger {
		return NewLogger(name, out)
	}
}
