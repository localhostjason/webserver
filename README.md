# webserver
go webserver daemon 库

主要功能：
1. tcp server
2. daemon
3. conf load
4. db for mysql sqlite
5. daemonx 支持 windows, linux
6. daemonx 支持 gin+grpc server 

用途：
1. gin server daemon
```golang
import (
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/daemonx"
	_ "github.com/localhostjason/webserver/db"
)

func getU(c *gin.Context) {
	c.JSON(200, make(map[string]string))
}

// 业务模块
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
	const defaultConfigPath = "配置目录"
	s := daemonx.NewMainServer(defaultConfigPath)

    // 可加载一些任务，比如：定时器任务
    s.LoadTasks(NewTestTask())
    //s.LoadGrpcServerApi(...) 配置文件 enable_grpc = true 开启后。s.load grpc api 才有意义
    
    s.LoadView(SetView)
    s.Run()
}
```

1. -d  获取配置
2. -x  debug跑
3. -k start 开启daemon
4. -k stop  关闭daemon

windows

1. -k install 安装服务
2. -k uninstall 卸载服务
3. -k start 开启daemon
4. -k stop  关闭daemon