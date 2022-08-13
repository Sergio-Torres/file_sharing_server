package main

import (
	"fmt"
	"log"
	"net"

)

const BUFFERSIZE = 1024

func main(){
    s := newServer()
    go s.run()

    l, err := net.Listen("tcp", ":8888")
    if err != nil{
        log.Fatalf("unable to start server: %s", err.Error())
    }

    defer l.Close()
    log.Printf("started server on :8888")

    for{
        conn, err := l.Accept()
        if err != nil{
            fmt.Printf("unable to accept connection %s", err.Error())
            continue
        }
        
        go s.newClient(conn)
    }
}
