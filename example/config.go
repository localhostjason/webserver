package example

import (
	"errors"
	"fmt"
	"github.com/localhostjason/webserver/db"
	"github.com/localhostjason/webserver/server/config"
)

// 自定义的配置路径 可配置
const defaultConfigPath = "D:\\center\\console\\console.json"

func dumpDefaultConfig() {
	content, err := config.GeneDefaultConfig()
	if err != nil {
		fmt.Println("failed to generate default config")
	} else {
		fmt.Println(string(content))
	}
}

func syncDB() (err error) {
	err = db.Connect()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to migrate:%v", err))
	}
	err = db.Migrate()
	if err != nil {
		return
	}

	err = db.InitData()
	return
}
