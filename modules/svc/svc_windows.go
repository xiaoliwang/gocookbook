package svc

import(
    "os"
)

type windowsService struct {
    isInteractive bool
}

func Run(prg Service, sig ...os.Signal) error {
    ws := &windowsService{ true }
    prg.Init(ws)
    ws.run(prg, sig...)
    return nil;
}

func (ws *windowsService) run(prg Service, sig ...os.Signal) error {
    prg.Start()

    signalChan := make(chan os.Signal, 1)
    signalNotify(signalChan, sig...)
    <-signalChan

    prg.Stop()

    return nil
}



func (ws *windowsService) IsWindowsService() bool {
    return !ws.isInteractive
}