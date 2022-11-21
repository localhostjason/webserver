package main

import (
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/daemonx"
	_ "github.com/localhostjason/webserver/db"
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
	const defaultConfigPath = "D:\\center\\console\\console.json"
	s := daemonx.NewMainServer(defaultConfigPath)

	// 可加载一些任务，比如：定时器任务
	//s.LoadTasks(NewTestTask())
	//s.LoadGrpcServerApi(...) 配置文件 enable_grpc = true 开启后。s.load grpc api 才有意义

	s.LoadView(SetView)
	s.Run()
}
