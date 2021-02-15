package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

const (
	Username = ""
	Passwd   = ""
	Host     = ""
	Port     = ""
	Dbname   = ""
)

func InitDb(username, passwd, host, port, dbname string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
			username,
			passwd,
			host,
			port,
			dbname))
	if err != nil {
		return nil, err
	}

	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(15)

	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.SnakeMapper{})

	return engine, nil
}
