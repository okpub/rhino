package core

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

/*
系统启动的时间
*/
var sys_time time.Time

func init() { sys_time = time.Now() }

func Uptime() time.Duration { return time.Since(sys_time) }

/*
*能够捕获系统挂了的原因
 */
const (
	defaultCrashPath = "./crash.log"
)

func SysErr() error {
	//path := "./crash" + time.Now().Format("2006-01-02 15:04:05") + ".log"
	file, err := os.OpenFile(defaultCrashPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err == nil {
		return syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	}
	fmt.Println("#system crash catch err:", err)
	return err
}
