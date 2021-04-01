package connection

import (
	"fmt"
	"honvid/pkg/log"
	"honvid/pkg/orm/schema"
	"reflect"
	"strings"
)

func (c *Connection) Model(value interface{}) *Connection {
	// nil or different model, update refTable
	if c.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(c.refTable.Model) {
		c.refTable = schema.Parse(value, c.dialect)
	}
	return c
}

func (c *Connection) RefTable() *schema.Schema {
	if c.refTable == nil {
		log.Error("Model is not set")
	}
	return c.refTable
}

func (c *Connection) CreateTable() error {
	table := c.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := c.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (c *Connection) DropTable() error {
	_, err := c.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", c.RefTable().Name)).Exec()
	return err
}
