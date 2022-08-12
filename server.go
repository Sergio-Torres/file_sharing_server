package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type server struct{
    channels map[string]*channel
    commands chan command
}

func newServer() *server{
    return &server{
        channels: make(map[string]*channel),
        commands: make(chan command),
    }
}
func (s *server) newClient(conn net.Conn){
    log.Printf("New client has connected: %s", conn.RemoteAddr().String())
    
    c := &client{
        conn : conn,
        nick: "none",
        commands: s.commands,
    }
    

}

/*
var count = 0

func handleConnection(c net.Conn){
    fmt.Print(".")
    for{
        netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil{
            fmt.Println(err)
            return
        }

        temp := strings.TrimSpace(string(netData))
        if temp == "STOP"{
            break
        }

        fmt.Println(temp)
        counter := strconv.Itoa(count) + "\n"
        c.Write([]byte(string(counter)))
    }
    c.Close()
}

func main(){
    arguments := os.Args
    if len(arguments) ==1{
        fmt.Println("Please provide a port number!")
        return
    }

    PORT := ":" + arguments[1]
    l, err := net.Listen("tcp4", PORT)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()

    for {
        c, err := l.Accept()
        if err != nil{
            fmt.Println(err)
            return
        }

        go handleConnection(c)
        count++
    }

}*/
