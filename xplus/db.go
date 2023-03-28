package xplus

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type Config struct {
	DataSourceName  string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogLevel        log.LogLevel

	ShowSQL      bool
	ShowExecTime bool
}

var (
	once     sync.Once
	dbEngine *xorm.Engine
)

func NewEngineWithConfig(cfg *Config) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	once.Do(func() {
		engine, _ := xorm.NewEngine("mysql", cfg.DataSourceName)
		if cfg.MaxIdleConns > 0 {
			engine.SetMaxIdleConns(cfg.MaxIdleConns)
		}

		if cfg.MaxOpenConns > 0 {
			engine.SetMaxOpenConns(cfg.MaxOpenConns)
		}

		if cfg.ConnMaxLifetime > 0 {
			engine.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		}
		if cfg.LogLevel > 0 {
			engine.SetLogLevel(cfg.LogLevel)
		}
		engine.ShowSQL(cfg.ShowSQL)
		dbEngine = engine
	})
	return nil
}

func Engine() (*xorm.Engine, error) {
	if dbEngine == nil {
		return nil, errors.New("engine is not init")
	}
	return dbEngine, nil
}

func SetLogger(logger core.ILogger) {
	dbEngine.SetLogger(logger)
}

func Close() {
	if dbEngine != nil {
		_ = dbEngine.Close()
	}
}
