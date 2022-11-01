package main

import (
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/daemonx"
	_ "github.com/localhostjason/webserver/db"
)

func getU(c *gin.Context) {
	c.JSON(200, make(map[string]string))
}

func SetView(r *gin.Engine) error {
	api := r.Group("api")
	{
		api.GET("/u", getU)
	}
	return nil
}

// 例子
func main() {
	// 自定义的配置路径 可配置
	const defaultConfigPath = "D:\\center\\console\\console.json"
	s := daemonx.NewMainServer(defaultConfigPath, SetView)
	s.Run()
}
