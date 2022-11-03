package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/daemonx"
	"github.com/localhostjason/webserver/server/config"
)

type MainWorkFunc func(r *gin.Engine) error

type TMainServer struct {
	DefaultConfigPath string
	SetMainWorkFunc   MainWorkFunc
}

func NewTMainServer(configPath string, setMainWorkFunc MainWorkFunc) *TMainServer {
	return &TMainServer{DefaultConfigPath: configPath, SetMainWorkFunc: setMainWorkFunc}
}

// Run 可根据自己业务 替换扩展
func (m *TMainServer) Run() {
	configPath := flag.String("p", m.DefaultConfigPath, "path to config")
	testCmd := flag.Bool("t", false, "test cmd")

	// for service
	singleMode := flag.Bool("x", false, "start, no daemon/service mode")
	svcCMD := flag.String("k", "", "cmds:start|stop|status, windows: install|uninstall")

	flag.Parse()

	if err := config.SetConfigFile(*configPath); err != nil {
		fmt.Println("failed to set config path", *configPath, err)
		return
	}

	if *testCmd {
		fmt.Println("test cmd -t")
		return
	}

	daemonx.SetMainWorkFunc = m.SetMainWorkFunc
	daemonx.RunService(*singleMode, *svcCMD)
}

func getUx(c *gin.Context) {
	c.JSON(200, make(map[string]string))
}

func SetViewx(r *gin.Engine) error {
	api := r.Group("api")
	{
		api.GET("/u", getUx)
	}
	return nil
}

func main() {
	// 自定义的配置路径 可配置
	const defaultConfigPath = "D:\\center\\console\\console.json"
	s := NewTMainServer(defaultConfigPath, SetViewx)
	s.Run()
}
