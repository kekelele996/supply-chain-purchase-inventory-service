package service

import "strings"

// isDuplicateErrMsg 判断错误是否为数据库唯一键冲突。
func isDuplicateErrMsg(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique") || strings.Contains(msg, "1062")
}
