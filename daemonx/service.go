package daemonx

import (
	"errors"
	"fmt"
	"github.com/localhostjason/webserver/svc"

	log "github.com/sirupsen/logrus"
)

const (
	START  = "start"
	STOP   = "stop"
	STATUS = "status"
)

func createService(singleMode bool) (*svc.Svc, error) {
	prc := NewProc(singleMode)
	svcx, err := NewService(prc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create program:%v", err))
	}
	return svcx, nil
}

func runService(singleMode bool, cmd string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	s, err := createService(true)
	if err != nil {
		log.Fatalln("failed to start", err)
	}

	if singleMode {
		s.RunSingle()
	} else if cmd == START {
		s.Run()
	} else if cmd == STOP {
		s.Stop()
	} else if cmd == STATUS {
		if s.Status() == nil {
			fmt.Println("running")
		} else {
			fmt.Println("not running")
		}
	}
}
