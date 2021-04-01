package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	fmt.Println(typ)
	switch typ.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int, reflect.Int16, reflect.Int32:
		return "int"
	case reflect.Uint8:
		return "tinyint unsigned"

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int unsigned"

	case reflect.Int64:
		return "bigint"

	case reflect.Uint64:
		return "bigint unsigned"

	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
		//if typ.Len() > 0 && typ.Len() < 65532 {
		//	return fmt.Sprintf("varchar(%d)", typ.Len())
		//} else {
		//	return "longtext"
		//}
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %m (%m)", typ.Type().Name(), typ.Kind()))
}
