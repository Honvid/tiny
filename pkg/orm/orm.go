package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"honvid/pkg/log"
	"honvid/pkg/orm/connection"
	"honvid/pkg/orm/dialect"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func New(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *connection.Connection {
	return connection.New(engine.db, engine.dialect)
}
