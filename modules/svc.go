package main

import (
    "fmt"
    "os"
    "net"
    "bufio"
    "sync"
    "github.com/gocookbook/modules/svc"
    "syscall"
)

type wait struct {
    sync.WaitGroup
}

type program struct {
    w wait
    tcp net.Listener
}

func main() {
    prg := &program{}
    // SIGINT Ctrl-C
    // SIGTERM kill
    svc.Run(prg, syscall.SIGINT, syscall.SIGTERM)
}

func (w *wait) Wrap(cb func()) {
    w.Add(1)
    go func() {
        cb()
        fmt.Println("job finished")
        w.Done()
    }()
}

func (p *program) Init(e svc.Environment) {
    if (e.IsWindowsService()) {
        fmt.Println("This is windows")
    }
}

func (p *program) Start() {
    listener, err := net.Listen("tcp", ":12345")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    p.tcp = listener
    p.w.Wrap(func(){
        Server(p.tcp)
    })

    fmt.Println("Program has started")
}

func (p *program) Stop() {
    if p.tcp != nil {
        p.tcp.Close()
    }
    fmt.Println("program stop")
    p.w.Wait()
}

func Server(listener net.Listener) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("conn:", err)
            break
        }
        go handle(conn)
    }

    fmt.Println("Accept closed")
}

func handle(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("reader:", err)
            break
        }
        fmt.Println(line)
    }

    conn.Close()
    fmt.Println("connection closed")
}