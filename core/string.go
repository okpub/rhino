package core

import "strconv"

//int转string
func Itoa(val int) string {
	return strconv.Itoa(val)
}

//string转int
func Atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err == nil {
		return val
	}
	return 0
}

//string转int64
func Atol(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return val
	}
	return 0
}

//int64转string
func Ltoa(val int64) string {
	return strconv.FormatInt(val, 10)
}
