package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func Connect(dsn string) (*sql.DB, error) {
	var driver string
	var connStr string
	pos := strings.Index(dsn, "://")
	if pos > -1 {
		driver = dsn[0:pos]
		connStr = dsn[pos+len("://"):]
	} else {
		panic(fmt.Errorf("db dsn checked faild:%s", dsn))
	}

	return sql.Open(driver, connStr)
}
