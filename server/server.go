package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	conf    ConfigServer
	servers []*http.Server
	router  *gin.Engine
}

func (s *Server) LoadGrpcServerApi(loadFunc func(*grpc.Server)) {
	LoadGserverApiFunc = loadFunc
}

func NewServer() (*Server, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return newServerWithConf(cfg), nil
}

func newServerWithConf(c ConfigServer) *Server {
	r := gin.New()
	// r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/"})))
	if c.Gzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	r.Use(gin.Logger())
	r.Use(setResponseHeader(c))
	s := &Server{conf: c, router: r}
	return s
}

func (s *Server) SetLogConfig(toConsole bool) error {
	//gin.SetMode(gin.ReleaseMode)
	if !toConsole {
		gin.SetMode(gin.ReleaseMode)
	}
	return SetLogConfig(toConsole)
}

type hAddView func(r *gin.Engine) error

func (s *Server) SetRouter(h hAddView) error {
	return h(s.router)
}

// SetRecovery should be called before SetRouter
func (s *Server) SetRecovery(h gin.HandlerFunc) {
	s.router.Use(h)
}

// todo run multi server
//https://github.com/gin-gonic/gin#run-multiple-service-using-gin

func (s *Server) getListener() ([]net.Listener, error) {
	var listener []net.Listener
	if s.conf.Bind6 != "" {
		addr := fmt.Sprintf("[%s]:%d", s.conf.Bind6, s.conf.Port)
		l, err := net.Listen("tcp6", addr)
		if err != nil {
			return listener, fmt.Errorf("failed to listen %s:%v", addr, err)
		}
		listener = append(listener, l)
	}
	if s.conf.Bind != "" {
		addr := fmt.Sprintf("%s:%d", s.conf.Bind, s.conf.Port)
		l, err := net.Listen("tcp4", addr)
		if err != nil {
			return listener, fmt.Errorf("failed to listen %s:%v", addr, err)
		}
		listener = append(listener, l)
	}
	return listener, nil
}

func (s *Server) GetListeners() ([]net.Listener, error) {
	return s.getListener()
}

func (s *Server) Start() error {
	listeners, err := s.getListener()
	if err != nil {
		return err
	}
	for _, l := range listeners {
		server := s.buildServer(s.router, s.conf.ReadTimeout, s.conf.WriteTimeout,
			s.conf.MaxHeaderBytes, s.conf.MinTlsVersion)
		s.servers = append(s.servers, server)
		// TODO error handling
		go s.startServer(server, l)

		if s.conf.EnableGrpc {
			go func() {
				gs := NewGrpcServer(l)
				gs.StartGrpc(LoadGserverApiFunc)
			}()
		}
	}
	return nil
}

func (s *Server) startServer(server *http.Server, l net.Listener) error {
	if s.conf.SSL {
		if err := server.ServeTLS(l, s.conf.CertFile, s.conf.KeyFile); err != nil {
			return fmt.Errorf("failed to start server:%v", err)
		}
	} else {
		if err := server.Serve(l); err != nil {
			return fmt.Errorf("failed to start server:%v", err)
		}
	}
	return nil
}

var _tlsVersion = map[string]uint16{
	"1.0": tls.VersionTLS10,
	"1.1": tls.VersionTLS11,
	"1.2": tls.VersionTLS12,
	"1.3": tls.VersionTLS13,
}

func (s *Server) getTlsConfig(minTlsVersion uint16) *tls.Config {
	tlsConf := &tls.Config{
		MinVersion: minTlsVersion,
	}

	if s.conf.RequireCert {
		tlsConf.ClientAuth = tls.RequestClientCert
	}

	// 大于 tls 1.3 版本，不做处理
	if minTlsVersion >= tls.VersionTLS13 {
		return tlsConf
	}

	// 小于1.3版本， 指定tls加密算法
	tlsConf.PreferServerCipherSuites = true

	// 下面中列出的算法，剔除了DES，解决  CVE-2016-2183(SSL/TLS)漏洞
	tlsConf.CipherSuites = []uint16{
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_CHACHA20_POLY1305_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	}

	return tlsConf
}

func (s *Server) getTlsVersion(conf string) uint16 {
	if v, ok := _tlsVersion[conf]; ok {
		return v
	}
	return tls.VersionTLS12
}

func (s *Server) buildServer(r *gin.Engine, readTimeout, writeTimeout,
	maxHeaderBytes int, minTlsVersion string) *http.Server {
	return &http.Server{
		Handler:        r,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
		//TLSConfig:      &tls.Config{MinVersion: s.getTlsVersion(minTlsVersion)},
		TLSConfig: s.getTlsConfig(s.getTlsVersion(minTlsVersion)),
	}
}

func (s *Server) Stop() error {
	var wg sync.WaitGroup
	wg.Add(len(s.servers))
	for _, server := range s.servers {
		go s.stopWaitServer(server, &wg)
		if s.conf.EnableGrpc {
			go StopGrpc()
		}
	}
	return nil
}

func (s *Server) stopWaitServer(server *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.conf.StopWait)*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorln("server not shutdown gracefully:", err)
	}
}

// add response header
func setResponseHeader(conf ConfigServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range conf.ResponseHeader {
			c.Writer.Header().Set(key, value)
		}
		c.Next()
	}
}
