package example

import (
	"errors"
	"os"
	"syscall"
	"webserver/server"
	"webserver/svc"

	log "github.com/sirupsen/logrus"
)

type MainProc struct {
	singleMode bool
	quit       chan bool
}

func (m *MainProc) Stop() {
	m.quit <- true
}

func NewProc(singleMode bool) *MainProc {
	return &MainProc{singleMode: singleMode, quit: make(chan bool)}
}

func (m *MainProc) Run(svc *svc.Svc) {
	s, err := startServer(m.singleMode)
	if err != nil {
		return
	}

	<-m.quit
	_ = s.Stop()
}

func (m *MainProc) SigHandlers() map[os.Signal]svc.SignalHandlerFunc {
	return map[os.Signal]svc.SignalHandlerFunc{
		syscall.SIGTERM: m.handleSigTerm,
		os.Interrupt:    m.handleSigTerm,
	}
}

func (m *MainProc) handleSigTerm(sig os.Signal) (err error) {
	m.quit <- true
	return errors.New("quit by signal " + sig.String())
}

func startServer(toConsole bool) (*server.Server, error) {

	s, err := server.NewServer()
	if err != nil {
		log.Fatalln("failed to start:", err)
	}

	err = s.SetLogConfig(toConsole)
	if err != nil {
		log.Fatalln("failed to set log:", err)
	}

	//s.SetRecovery(uv.DefaultRecovery(false))
	//err = s.SetRouter(view.SetView)
	//if err != nil {
	//	return nil, err
	//}

	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}

	return s, nil
}
