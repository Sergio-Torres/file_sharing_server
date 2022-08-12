package main

import (
	"errors"
	"fmt"
	"log"
	"net"
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
func (s *server) run(){
    for cmd := range s.commands{
        switch cmd.id{
        case cmd_nick:
            s.nick(cmd.client, cmd.args)
        case cmd_join:
            s.join(cmd.client, cmd.args)
        case cmd_channels:
            s.listChannels(cmd.client, cmd.args)
        case cmd_msg:
            s.msg(cmd.client, cmd.args)
        case cmd_exit:
            s.exit(cmd.client, cmd.args)
        }


    }
}

func (s *server) newClient(conn net.Conn){
    log.Printf("New client has connected: %s", conn.RemoteAddr().String())
    
    c := &client{
        conn : conn,
        nick: "none",
        commands: s.commands,
    }
    
    c.readInput()
}

func (s *server) nick(c *client, args []string){
    c.nick = args[1]
    c.msg(fmt.Sprintf("[*] All right, your nickname will be %s", c.nick))
}

func (s *server) join(c *client, args []string){
    channelName := args[1]
    ch, ok := s.channels[channelName]

    if !ok {
        ch = &channel{
            name: channelName,
            members: make(map[net.Addr]*client),
        }
        s.channels[channelName] = ch
    }

    ch.members[c.conn.RemoteAddr()] = c

    s.exitCurrentChannel(c)   
    c.channel = ch
    
    ch.broadcast(c, fmt.Sprintf("[*] %s has joined the channel", c.nick))
    c.msg(fmt.Sprintf("[!] welcome to %s", ch.name))
}

func (s *server) listChannels(c *client, args []string){
    var channels []string
    for name := range s.channels{
        channels = append(channels, name)
    }

    c.msg(fmt.Sprintf("Available channels are: %s", strings.Join(channels, ",")))
}

func (s *server) msg(c *client, args []string){
    if c.channel == nil{
        c.err(errors.New("[*] you musy join the room first"))
        return
    }

    c.channel.broadcast(c, c.nick + ":" + strings.Join(args[1:len(args)], " "))
}
func (s *server) exit(c *client, args []string){
    log.Printf("[!] client has disconnected: %s", c.conn.RemoteAddr().String())

    s.exitCurrentChannel(c)
    c.msg("Bye..")
    c.conn.Close()

}

func (s *server) exitCurrentChannel(c *client){
    if c.channel !=nil{
        delete(c.channel.members, c.conn.RemoteAddr())
        c.channel.broadcast(c, fmt.Sprintf("[*] %s has left the channel", c.nick))
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
