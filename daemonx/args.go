package daemonx

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/server/config"
)

type MainWorkFunc func(r *gin.Engine) error

var SetMainWorkFunc func(r *gin.Engine) error

type MainServer struct {
	DefaultConfigPath string
	SetMainWorkFunc   MainWorkFunc
}

func NewMainServer(configPath string, setMainWorkFunc MainWorkFunc) *MainServer {
	return &MainServer{DefaultConfigPath: configPath, SetMainWorkFunc: setMainWorkFunc}
}

func (m *MainServer) Run() {
	configPath := flag.String("p", m.DefaultConfigPath, "path to config")
	initDB := flag.Bool("i", false, "int db")
	dumpConfig := flag.Bool("d", false, "dump default config")

	// for service
	singleMode := flag.Bool("x", false, "start, no daemon/service mode")
	svcCMD := flag.String("k", "", "cmds:start|stop|status, windows: install|uninstall")

	flag.Parse()

	if err := config.SetConfigFile(*configPath); err != nil {
		fmt.Println("failed to set config path", *configPath, err)
		return
	}

	// commands

	if *dumpConfig {
		dumpDefaultConfig()
		return
	}

	// DB 初始表结构和默认值
	if *initDB {
		if err := syncDB(); err != nil {
			fmt.Println("error when sync db schema", err)
			return
		}
		fmt.Println("success: sync db schema")
		return
	}

	SetMainWorkFunc = m.SetMainWorkFunc
	RunService(*singleMode, *svcCMD)
}
