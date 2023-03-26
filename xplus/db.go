package xplus

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
	"time"
	"xorm.io/core"
)

type Config struct {
	DataSourceName  string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Logger          core.ILogger
	ShowSQL         bool
	ShowExecTime    bool
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
		if cfg.Logger != nil {
			engine.SetLogger(cfg.Logger)
		}
		engine.ShowSQL(cfg.ShowSQL)
		engine.ShowExecTime(cfg.ShowExecTime)
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

func Close() {
	if dbEngine != nil {
		_ = dbEngine.Close()
	}
}
