package connection

import (
	"database/sql"
	"fmt"
	"honvid/pkg/log"
	"honvid/pkg/orm/dialect"
	"honvid/pkg/orm/schema"
	"strings"
)

type Connection struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVars  []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect) *Connection {
	return &Connection{
		db:      db,
		dialect: dialect,
	}
}

func (c *Connection) Clear() {
	c.sql.Reset()
	c.sqlVars = nil
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) Raw(sql string, values ...interface{}) *Connection {
	c.sql.WriteString(sql)
	c.sql.WriteString(" ")
	c.sqlVars = append(c.sqlVars, values...)
	return c
}

// Exec raw sql with sqlVars
func (c *Connection) Exec() (result sql.Result, err error) {
	defer c.Clear()
	log.Info(c.sql.String(), c.sqlVars)
	if result, err = c.DB().Exec(c.sql.String(), c.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db
func (c *Connection) QueryRow() *sql.Row {
	defer c.Clear()
	log.Info(c.sql.String(), c.sqlVars)
	return c.DB().QueryRow(c.sql.String(), c.sqlVars...)
}

// QueryRows gets a list of records from db
func (c *Connection) QueryRows() (rows *sql.Rows, err error) {
	defer c.Clear()
	log.Info(c.sql.String(), c.sqlVars)
	if rows, err = c.DB().Query(c.sql.String(), c.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (c *Connection) HasTable(tableName string) bool {
	var name string
	// allow mysql database name with '-' character
	if err := c.db.QueryRow(fmt.Sprintf("SHOW TABLES FROM `%s` WHERE `Tables_in_%s` = ?", "demo", "demo"), tableName).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	} else {
		return true
	}
}
