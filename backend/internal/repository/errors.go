// Package repository 定义数据访问层。
package repository

import "errors"

// 仓储层哨兵错误，上层通过 errors.Is 判断。
var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicateKey = errors.New("duplicate key")
	ErrConflict     = errors.New("data conflict")
)
