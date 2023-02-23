package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/localhostjason/webserver/db"
	"github.com/localhostjason/webserver/example"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type TestTask struct {
	Quit chan bool
	Wg   *sync.WaitGroup
}

func NewTestTask() *TestTask {
	return &TestTask{Quit: make(chan bool)}
}

func (t *TestTask) Start() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-t.Quit:
				t.Wg.Done()
				return
			case <-ticker.C:
				logrus.Info("running")

			}
		}
	}()
}

func (t *TestTask) SetWg(wg *sync.WaitGroup) {
	t.Wg = wg
}

func (t *TestTask) Stop() {
	t.Quit <- true
}

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
	//const defaultConfigPath = "./example.json"
	s := example.NewMainServer() // 默认在当前目录 自动生成配置文件和目录

	//s.SetServerConfigFile(defaultConfigPath) // 可指定 自己配置文件

	//s.LoadGrpcServerApi(...) 配置文件 enable_grpc = true 开启后。s.load grpc api 才有意义

	s.LoadView(SetView)
	s.Run()
}
