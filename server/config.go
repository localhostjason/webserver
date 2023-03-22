package server

import (
	"github.com/localhostjason/webserver/server/config"
	"os"
	"path/filepath"
)

const _key = "webserver"

func PathExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func InitDefaultServerConfigFile(configFile string) string {
	if configFile != "" {
		return configFile
	}

	execPath, _ := GetExeDir()
	configDir := filepath.Join(execPath, "config")
	logDir := filepath.Join(execPath, "log")
	file := filepath.Join(configDir, "server.json")

	if !PathExists(configDir) {
		_ = os.MkdirAll(configDir, os.ModePerm)
	}

	if !PathExists(logDir) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	if PathExists(file) {
		return file
	}

	content, _ := config.GeneDefaultConfig()
	_ = os.WriteFile(file, content, os.ModePerm)
	return file
}

func init() {
	execPath, _ := GetExeDir()
	logPath := filepath.Join(execPath, "log")
	configPath := filepath.Join(execPath, "config")
	certFile := filepath.Join(configPath, "web.crt")
	keyFile := filepath.Join(configPath, "web.key")

	c := ConfigServer{
		LogLevel:       "info",
		LogPath:        logPath,
		AccessLog:      "access-%Y%m%d.log",
		ErrorLog:       "error-%Y%m%d.log",
		SysLog:         "sys-%Y%m%d.log",
		SSL:            false,
		RequireCert:    false,
		MinTlsVersion:  "1.2",
		CertFile:       certFile,
		KeyFile:        keyFile,
		Bind6:          "::",
		Bind:           "0.0.0.0",
		Port:           8088,
		ReadTimeout:    60,
		WriteTimeout:   60,
		MaxHeaderBytes: 2 << 20, // 2MB
		StopWait:       5,
		EnableGrpc:     false,

		ResponseHeader: map[string]string{
			"X-Frame-Options":           "SAMEORIGIN",
			"X-Content-Type-Options":    "nosniff",
			"X-XSS-Protection":          "1;mode=block",
			"Content-Security-Policy":   "default-src 'self' 'unsafe-inline' 'unsafe-eval';img-src 'self' data:;font-src 'self';frame-ancestors 'self'",
			"Strict-Transport-Security": "max-age=172800",
			"Referrer-Policy":           "strict-origin-when-cross-origin",
		},
	}
	_ = config.RegConfig(_key, c)
}

type ConfigServer struct {
	// log
	LogLevel  string `json:"log_level"`
	LogPath   string `json:"log_path"`
	AccessLog string `json:"access_log"`
	ErrorLog  string `json:"error_log"`
	SysLog    string `json:"sys_log"`

	// listen
	SSL           bool   `json:"ssl"`
	RequireCert   bool   `json:"require_cert"`
	MinTlsVersion string `json:"min_tls_version"`
	CertFile      string `json:"cert_file"`
	KeyFile       string `json:"key_file"`
	Bind6         string `json:"bind_6"`
	Bind          string `json:"bind"`
	Port          int    `json:"port"`

	// http
	ReadTimeout    int `json:"read_timeout"`
	WriteTimeout   int `json:"write_timeout"`
	MaxHeaderBytes int `json:"max_header_bytes"`

	// gzip
	Gzip bool `json:"gzip"`

	//
	StopWait int `json:"stop_wait"`

	// response Header 响应头
	ResponseHeader map[string]string `json:"response_header"`

	// grpc
	EnableGrpc bool `json:"enable_grpc"`
}

func GetConfig() (ConfigServer, error) {
	var c ConfigServer
	err := config.GetConfig(_key, &c)
	// TODO validation , port, bind etc..
	return c, err
}
