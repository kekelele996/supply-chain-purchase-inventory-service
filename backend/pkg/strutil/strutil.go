// Package strutil 提供与业务无关的字符串与编号生成工具。
package strutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// RandomDigits 生成长度为 n 的随机数字串（用于采购单号等）。
func RandomDigits(n int) string {
	if n <= 0 {
		n = 4
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			sb.WriteByte('0')
			continue
		}
		sb.WriteByte(byte('0' + v.Int64()))
	}
	return sb.String()
}

// OrderNo 生成采购单号，格式 PO-YYYYMMDD-XXXX。
func OrderNo(now time.Time, seq int) string {
	return fmt.Sprintf("PO-%s-%04d", now.Format("20060102"), seq)
}

// Truncate 截断字符串到指定长度（按 rune 计数），超出加省略号。
func Truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
