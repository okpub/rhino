package core

import (
	"fmt"
	"os"
	"syscall"
)

/*
*能够捕获系统挂了的原因
 */
func SysErr() error {
	//path := "./crash_" + DateOf(FmtSec) + ".log"
	logFile, err := os.OpenFile("./crash.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err == nil {
		return syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
	}
	fmt.Println("#crash err:", err)
	return err
}
