package webserver

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"webserver/config"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func createLogDir(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log dir %s:%v", logDir, err)
	}
	return nil
}

func SetLogConfig(toConsole bool) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	if err := createLogDir(cfg.LogPath); err != nil {
		return err
	}
	accessLog := filepath.Join(cfg.LogPath, cfg.AccessLog)
	errLog := filepath.Join(cfg.LogPath, cfg.ErrorLog)
	syslog := filepath.Join(cfg.LogPath, cfg.SysLog)

	setGinLog(accessLog, errLog, toConsole)
	setSysLog(syslog, cfg.LogLevel, toConsole)
	return nil
}

func setGinLog(accessLog, errLog string, toConsole bool) {
	accessLogRotate, err := rotatelogs.New(accessLog)
	if err != nil {
		log.Fatalln("failed to create access log file ", accessLog, err)
	}

	if toConsole {
		gin.DefaultWriter = io.MultiWriter(accessLogRotate, os.Stdout)
	} else {
		gin.DefaultWriter = accessLogRotate
	}

	errLogRotate, err := rotatelogs.New(errLog)
	if err != nil {
		log.Fatalln("failed to create error log file ", errLog, err)
	}

	if toConsole {
		gin.DefaultErrorWriter = io.MultiWriter(errLogRotate, os.Stdout)
	} else {
		gin.DefaultErrorWriter = errLogRotate
	}
}

func setSysLog(sysLog, logLevel string, toConsole bool) {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println("failed to parse log level:", logLevel, "use default :", "info")
	} else {
		lvl = log.InfoLevel
	}

	log.SetLevel(lvl)
	log.SetFormatter(new(LogFormatter))
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	logrotate, err := rotatelogs.New(sysLog)
	if err != nil {
		log.Fatalln("failed to create sys log file ", sysLog, err)
	}

	if toConsole {
		log.SetOutput(io.MultiWriter(os.Stdout, logrotate))
	} else {
		log.SetOutput(logrotate)
	}

}

// LogFormatter 日志自定义格式
type LogFormatter struct{}

// Format 格式详情
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format(config.TimeFormat)
	var file string
	var lenLine int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		lenLine = entry.Caller.Line
	}
	//fmt.Println(entry.Data)
	msg := fmt.Sprintf("%s [%s:%d][GOID:%d][%s] %s\n", timestamp, file, lenLine, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
