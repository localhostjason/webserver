package example

import (
	"fmt"
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
