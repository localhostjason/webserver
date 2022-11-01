package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ConnectWithSqliteConfig 连接，检验配置是否正确
func ConnectWithSqliteConfig(c SqliteDBConfig) error {

	db, err := gorm.Open(sqlite.Open(c.DbFile), &gorm.Config{
		FullSaveAssociations:   true,
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return err
	}

	if c.Debug {
		db = db.Debug()
	}

	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect databse %s:%v", c.DbFile, err))
	}
	DB = db
	return nil
}
