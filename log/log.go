package log

import (
	"sync/atomic"
	"time"
)

type Level int32

const (
	MinLevel = Level(iota)
	DebugLevel
	InfoLevel
	ErrorLevel
	OffLevel
)

//class Logger
type Logger struct {
	level   Level
	prefix  string
	context []Field
}

func New(level Level, prefix string, context ...Field) *Logger {
	return &Logger{level: level, prefix: prefix, context: context}
}

func (this Logger) With(fields ...Field) *Logger {
	this.context = append(this.context, fields...)
	return &this
}

func (this *Logger) Level() Level {
	return Level(atomic.LoadInt32((*int32)(&this.level)))
}

func (this *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&this.level), int32(level))
}

//publish func
func (this *Logger) Debug(msg string, fields ...Field) {
	if this.Level() < InfoLevel {
		this.Publish(Event{Time: time.Now(), Level: DebugLevel, Prefix: this.prefix, Message: msg, Context: this.context, Fields: fields})
	}
}

func (this *Logger) Info(msg string, fields ...Field) {
	if this.Level() < ErrorLevel {
		this.Publish(Event{Time: time.Now(), Level: InfoLevel, Prefix: this.prefix, Message: msg, Context: this.context, Fields: fields})
	}
}

func (this *Logger) Error(msg string, fields ...Field) {
	if this.Level() < OffLevel {
		this.Publish(Event{Time: time.Now(), Level: ErrorLevel, Prefix: this.prefix, Message: msg, Context: this.context, Fields: fields})
	}
}

//override publish
func (this *Logger) Publish(evt Event) {
	stage.Publish(evt)
}
