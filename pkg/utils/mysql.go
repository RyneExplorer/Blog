package utils

import (
	"errors"
	"strings"

	drivermysql "github.com/go-sql-driver/mysql"
)

// IsMySQLDuplicateKey 判断是否为 MySQL 唯一键冲突
func IsMySQLDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	var me *drivermysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate entry")
}
