package main

import (
	"log"
	"net"
)


func main(){
    s := newServer()
    listener, err := net.Listener("tcp", "1234")
    if err != nil{
        log.Fatal("unable to start server: %s", err.Error())
    }

    defer listener.Close()
    log.Printf("started server on :1234")

    for{
        conn, err := listener.Accept()
        if err != nil{
            log.Printf("unable to accepy connection %s", err.Error())
            continue
        }
        
        s.newClient(conn)
    }
}
