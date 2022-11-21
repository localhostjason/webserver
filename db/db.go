package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var DB *gorm.DB // 作为全局访问的db连接

func Connect() error {
	dbType, enable := getDbTypeConfig()
	if !enable {
		fmt.Println("no use db")
		return nil
	}

	switch dbType {
	case "sqlite":
		return ConnectWithSqliteConfig(getSqliteConfig())
	case "mysql":
		return ConnectWithMysqlConfig(getMysqlConfig())
	default:
		return errors.New(fmt.Sprintf("db type is not right， err:%v", dbType))
	}
}

func DBEnable() bool {
	_, enable := getDbTypeConfig()
	return enable
}

var tableModels []interface{}

// Migrate 初始化或升级表结构
func Migrate() error {
	if err := DB.AutoMigrate(tableModels...); err != nil {
		return errors.New(fmt.Sprintf("failed to migrate database:%v", err))
	}
	return nil
}

// RegTables 其他模块注册需要访问的表, 会被自动创建
func RegTables(tables ...interface{}) {
	tableModels = append(tableModels, tables...)
}

type InitDataHandler func() error

var _initHooks []InitDataHandler

// AddInitHook db连接后执行的函数， 可用于初始化数据等
func AddInitHook(h InitDataHandler) {
	_initHooks = append(_initHooks, h)
}

func InitData() error {
	for _, h := range _initHooks {
		if err := h(); err != nil {
			return err
		}
	}
	return nil
}
