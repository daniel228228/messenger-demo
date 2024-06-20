package log

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Level int8

const (
	LDebug Level = iota
	LInfo
	LWarn
	LError
	LFatal
	LOff
	LTrace Level = -1
)

type Logger struct {
	l        zerolog.Logger
	handlers []io.Writer
	discard  bool
	service  string
}

type FormatHandlerOptions struct {
	Color      bool
	TimeFormat string
}

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func NewLogger(service string) *Logger {
	return &Logger{
		service: service,
	}
}

func (l *Logger) Add(w io.Writer) *Logger {
	if !l.discard {
		l.handlers = append(l.handlers, w)
	}

	return l
}

func (l *Logger) AddHandler(w io.Writer, level Level) *Logger {
	if !l.discard {
		l.handlers = append(l.handlers, NewHandler(w, level))
	}

	return l
}

// it's heavy handler!
func (l *Logger) AddFormatHandler(w io.Writer, level Level, opts *FormatHandlerOptions) *Logger {
	if !l.discard {
		l.handlers = append(l.handlers, NewFormatHandler(w, level, opts))
	}

	return l
}

func (l *Logger) Launch() *Logger {
	l.discard = true
	l.l = zerolog.New(zerolog.MultiLevelWriter(l.handlers...)).With().Timestamp().Str("service", l.service).Logger()

	return l
}

func (l *Logger) Trace() *Event {
	return &Event{l.l.Trace()}
}

func (l *Logger) Debug() *Event {
	return &Event{l.l.Debug()}
}

func (l *Logger) Info() *Event {
	return &Event{l.l.Info()}
}

func (l *Logger) Warn() *Event {
	return &Event{l.l.Warn()}
}

func (l *Logger) Error() *Event {
	return &Event{l.l.Error().Stack()}
}

func (l *Logger) Fatal() *Event {
	return &Event{l.l.WithLevel(zerolog.FatalLevel)}
}

type handler struct {
	io.Writer
	Level int32
}

func NewHandler(w io.Writer, level Level) *handler {
	return &handler{
		Writer: w,
		Level:  levelToZerologCompatibleLevel(level),
	}
}

// it's heavy handler!
func NewFormatHandler(w io.Writer, level Level, opts *FormatHandlerOptions) *handler {
	h := &handler{
		Writer: zerolog.ConsoleWriter{
			Out:                 w,
			NoColor:             !getFormatHandlerOpt(opts).Color,
			TimeFormat:          getFormatHandlerOpt(opts).TimeFormat,
			PartsOrder:          nil,
			PartsExclude:        nil,
			FieldsExclude:       nil,
			FormatTimestamp:     nil,
			FormatLevel:         nil,
			FormatCaller:        nil,
			FormatMessage:       nil,
			FormatFieldName:     nil,
			FormatFieldValue:    nil,
			FormatErrFieldName:  nil,
			FormatErrFieldValue: nil,
			FormatExtra:         nil,
		},
	}

	h.ChangeLevel(level)

	return h
}

func (h *handler) Write(p []byte) (n int, err error) {
	return h.Writer.Write(p)
}

func (h *handler) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l >= zerolog.Level(atomic.LoadInt32(&h.Level)) {
		return h.Writer.Write(p)
	}

	return len(p), nil
}

func (h *handler) ChangeLevel(level Level) {
	atomic.SwapInt32(&h.Level, levelToZerologCompatibleLevel(level))
}

var levelLabels = [...]string{"trace", "debug", "info", "warn", "error", "fatal", "off"}

func GetLevelFromString(level string) (Level, error) {
	for i, v := range levelLabels {
		if level == v {
			return Level(i - 1), nil
		}
	}

	return LOff, fmt.Errorf("unknown log level %q", level)
}

func (l Level) String() string {
	if l < -1 || l > 5 {
		return "unknown"
	} else {
		return levelLabels[l+1]
	}
}

func levelToZerologCompatibleLevel(l Level) int32 {
	switch l {
	case LTrace, LDebug, LInfo, LWarn, LError:
		return int32(l)
	case LFatal:
		return int32(zerolog.FatalLevel)
	default:
		return int32(zerolog.Disabled)
	}
}

func getFormatHandlerOpt(opts *FormatHandlerOptions) FormatHandlerOptions {
	defaultFormatHandlerOpt := FormatHandlerOptions{
		Color:      true,
		TimeFormat: time.RFC3339,
	}

	if opts == nil {
		return defaultFormatHandlerOpt
	}

	newOpts := *opts

	if newOpts.TimeFormat == "" {
		newOpts.TimeFormat = defaultFormatHandlerOpt.TimeFormat
	}

	return newOpts
}
