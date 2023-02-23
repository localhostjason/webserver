package server

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var LoadGserverApiFunc func(server *grpc.Server)
var GServers []*grpc.Server

type GrpcServer struct {
	Listen net.Listener
}

func NewGrpcServer(listen net.Listener) *GrpcServer {
	return &GrpcServer{Listen: listen}
}

func (g *GrpcServer) StartGrpc(loadFunc func(*grpc.Server)) {
	//logrus.Info("start grpc server")

	grpcServer := grpc.NewServer()
	GServers = append(GServers, grpcServer)

	//logrus.Info("load grpc server handler")
	if LoadGserverApiFunc != nil {
		loadFunc(grpcServer)
	}

	go func() {
		err := grpcServer.Serve(g.Listen)
		if err != nil {
			logrus.Fatalln("err:", err)
		}
	}()
}

func StopGrpc() {
	for _, s := range GServers {
		if s == nil {
			continue
		}
		s.Stop()
	}
}
