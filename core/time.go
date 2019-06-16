package core

import (
	"time"
)

/*
时间格式(默认)
*/
const (
	FmtDay  = "2006-01-02"
	FmtSec  = "2006-01-02 15:04:05"
	FmtMsec = "2006-01-02 15:04:05.0"
)

/*
 * 天的时间
 */
const DayTime = time.Hour * 24

/*
秒为单位
*/
const (
	SecMin  int64 = 60
	SecHour int64 = 60 * SecMin
	SecDay  int64 = 24 * SecHour
	SecWeek int64 = 7 * SecDay
)

/*
Sleep
*/
func Sleep(msec int) {
	time.Sleep(time.Millisecond * time.Duration(msec))
}

func SleepSec(sec int) {
	time.Sleep(time.Second * time.Duration(sec))
}

/*
 * 时间戳
 */
func Nano() int64 {
	return time.Now().UnixNano()
}

func Mic() int64 {
	return Nano() / time.Microsecond.Nanoseconds()
}

func Msec() int64 {
	return Nano() / time.Millisecond.Nanoseconds()
}

func Sec() int64 {
	return Nano() / time.Second.Nanoseconds()
}

func Minute() int64 {
	return Nano() / time.Minute.Nanoseconds()
}

func Hour() int64 {
	return Nano() / time.Hour.Nanoseconds()
}

func Day() int64 {
	return Nano() / DayTime.Nanoseconds()
}

/*
 * 秒转time
 */
func SecTime(v int64) time.Time {
	return time.Unix(v, 0)
}

/*
* 纳秒转time
 */
func NanoTime(v int64) time.Time {
	return time.Unix(0, v)
}

/*
 * 格式化
 */
func Format(str string) string {
	return time.Now().Format(str)
}

/*
 * 之后的时间
 */
func Now() time.Time {
	return time.Now()
}

func AddMsec(delay int) time.Time {
	return time.Now().Add(time.Millisecond * time.Duration(delay))
}

func AddSec(delay int) time.Time {
	return time.Now().Add(time.Second * time.Duration(delay))
}

func AddMinute(delay int) time.Time {
	return time.Now().Add(time.Minute * time.Duration(delay))
}

func AddHour(delay int) time.Time {
	return time.Now().Add(time.Hour * time.Duration(delay))
}

func AddDay(delay int) time.Time {
	return time.Now().Add(DayTime * time.Duration(delay))
}

func Add(delay time.Duration) time.Time {
	return time.Now().Add(delay)
}

/*
 * 获取当天0点的时间戳(秒)
 */
func ZeroSec() int64 {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

/*
函数执行时间
*/
func Since(f func()) time.Duration {
	tm := time.Now()
	f()
	return time.Since(tm)
}
