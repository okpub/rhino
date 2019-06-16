package log

import (
	"testing"
	"time"
)

var (
	plog = New(DebugLevel, "[MAILBOX]")
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			plog.Debug("[ACTOR] Recovering", Object("error", r), Stack())
		}
	}()
	panic("test panic")
}

func TestLogger_With(t *testing.T) {
	base := New(DebugLevel, "Base")
	//
	l := base.With(Int("my", 12))
	l.Debug("我只是一个测试的", Int("bar", 11))
	time.Sleep(time.Millisecond)
}

func Benchmark_OffLevel_TwoFields(b *testing.B) {
	l := New(MinLevel, "")
	for i := 0; i < b.N; i++ {
		l.Debug("foo", Int("bar", 32), Bool("fum", false))
	}
}

func Benchmark_OffLevel_OnlyContext(b *testing.B) {
	l := New(MinLevel, "", Int("bar", 32), Bool("fum", false))
	for i := 0; i < b.N; i++ {
		l.Debug("foo")
	}
}

func Benchmark_DebugLevel_OnlyContext_OneSubscriber(b *testing.B) {
	s1 := Subscribe(func(Event) {})

	l := New(DebugLevel, "", Int("bar", 32), Bool("fum", false))
	for i := 0; i < b.N; i++ {
		l.Debug("foo")
	}
	Unsubscribe(s1)
}

func Benchmark_DebugLevel_OnlyContext_MultipleSubscribers(b *testing.B) {
	s1 := Subscribe(func(Event) {})
	s2 := Subscribe(func(Event) {})

	l := New(DebugLevel, "", Int("bar", 32), Bool("fum", false))
	for i := 0; i < b.N; i++ {
		l.Debug("foo")
	}

	s1.Unsubscribe()
	s2.Unsubscribe()
}
