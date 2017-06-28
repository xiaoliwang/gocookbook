package svc

import "os/signal"

var signalNotify = signal.Notify

type Service interface {
    Init(Environment)
    Start()
    Stop()
}

type Environment interface {
    IsWindowsService() bool
}