package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type service interface {
	Build() error
	Run() error
	Shutdown() error
}

func MustRun(s service) {
	if err := s.Build(); err != nil {
		log.Fatal(err.Error())
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	if err := s.Run(); err != nil {
		log.Fatal(err.Error())
	}
	<-stopChan
	if err := s.Shutdown(); err != nil {
		log.Fatal(err.Error())
	}
}
