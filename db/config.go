package db

import "github.com/localhostjason/webserver/server/config"

const _key = "db"

type MysqlDBConfig struct {
	User            string `json:"user"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	DB              string `json:"db"`
	Charset         string `json:"charset"`
	Timeout         int    `json:"timeout"`
	MultiStatements bool   `json:"multi_statements"`
	Debug           bool   `json:"debug"`
}

type SqliteDBConfig struct {
	DbFile string `json:"db_file"`
	Debug  bool   `json:"debug"`
}

type DbConfig struct {
	DbType string         `json:"db_type"` // one of mysql . sqlite
	Enable bool           `json:"enable"`
	Mysql  MysqlDBConfig  `json:"mysql"`
	Sqlite SqliteDBConfig `json:"sqlite"`
}

func init() {
	mc := MysqlDBConfig{
		User:            "root",
		Password:        "123456",
		Host:            "127.0.0.1",
		Port:            3306,
		DB:              "test",
		Charset:         "utf8mb4",
		MultiStatements: false,
		Timeout:         5,
		Debug:           false,
	}

	sc := SqliteDBConfig{
		DbFile: "data/data.db",
		Debug:  false,
	}

	c := DbConfig{
		DbType: "sqlite",
		Enable: true,
		Mysql:  mc,
		Sqlite: sc,
	}
	_ = config.RegConfig(_key, c)
}

func getMysqlConfig() MysqlDBConfig {
	var c DbConfig
	_ = config.GetConfig(_key, &c) // todo return error
	return c.Mysql
}

func getSqliteConfig() SqliteDBConfig {
	var c DbConfig
	_ = config.GetConfig(_key, &c) // todo return error
	return c.Sqlite
}

func getDbTypeConfig() (string, bool) {
	var c DbConfig
	_ = config.GetConfig(_key, &c)
	return c.DbType, c.Enable
}
