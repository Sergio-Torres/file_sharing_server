package main

import (
	"fmt"
	"log"
	"net"
)


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
            fmt.Printf("unable to accepy connection %s", err.Error())
            continue
        }
        
        go s.newClient(conn)
    }
}
