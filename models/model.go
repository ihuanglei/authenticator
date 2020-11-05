package models

import (
	"fmt"
	"time"

	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/logger"

	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/simplexwork/common"
)

const (
	_tablePrefix = "at_"
)

// package 变量下划线小写字母开头, 防止变量定义重复
var (
	_Engine   *xorm.Engine
	_IDWorker *common.IDWorker
	_Tables   = []interface{}{
		new(user),
		new(userInfo),
		new(userThird),
		new(userLogin),
		new(userAddress),
		new(dict),
		new(resource),
		new(role),
		new(roleResource),
	}
)

// Init .
func Init(config *config.Config) error {
	logLevel := config.Log
	serverID := config.Server.ID
	mysql := config.Mysql

	logger.Infof("Connect Database %s@%s:%v/%s", mysql.User, mysql.Host, mysql.Port, mysql.Database)

	var err error
	_Engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4",
		mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Database))
	if err != nil {
		return err
	}

	mapper := core.NewPrefixMapper(core.GonicMapper{}, _tablePrefix)
	_Engine.SetMapper(mapper)

	if mysql.UseCache {
		cacher := caches.NewLRUCacher2(caches.NewMemoryStore(), time.Second*10, 1000)
		_Engine.SetDefaultCacher(cacher)
	}

	_Engine.SetMaxIdleConns(mysql.MaxIdleConns)
	_Engine.SetMaxOpenConns(mysql.MaxOpenConns)
	_Engine.SetLogger(logger.NewXormLogger(5-logLevel, mysql.ShowSQL))
	_Engine.SetConnMaxLifetime(time.Second * time.Duration(mysql.MaxLifeTime))

	if err := _Engine.Ping(); err != nil {
		return err
	}

	if mysql.Sync {
		err = _Engine.StoreEngine("InnoDB").Sync2(_Tables...)
		if err != nil {
			return err
		}
	}

	_IDWorker = common.NewIDWorker(serverID)

	return nil
}

// DefauleEngine .
func DefauleEngine() *xorm.Engine {
	return _Engine
}
