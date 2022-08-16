package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
    "strings"
   
)

type server struct{
    channels map[string]*channel
    files map[string]*file
    commands chan command
}

func newServer() *server{
    return &server{
        channels: make(map[string]*channel),
        files: make(map[string]*file),
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
        case cmd_file:
            s.file(cmd.client, cmd.args)
        case cmd_files:
            s.listFiles(cmd.client, cmd.args)
        case cmd_save:
            s.saveFile(cmd.client, cmd.args)
        }


    }
}

func (s *server) newClient(conn net.Conn){
    log.Printf("New client has connected: %s", conn.RemoteAddr().String())
    
    c := &client{
        conn : conn,
        nick: "anonymous",
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
    log.Printf("%s has entered the channel %s", c.nick, ch.name)
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
        c.err(errors.New("[!] you must join the channel first"))
        return
    }

    c.channel.broadcast(c, c.nick + ":" + strings.Join(args[1:len(args)], " "))
}

func (s *server) file(c *client, args []string){
    fileN := args[1]
    fl, ok := s.files[fileN]

    if !ok {
        fl = &file{
            name: fileN,
            members: make(map[net.Addr]*client),
        }
        s.files[fileN] = fl
    }

    f, err := os.Open("./"+fileN)
    //evisar este error
    if err != nil{
        log.Fatal(err.Error())
    }
    

    defer f.Close()
    _, err = io.Copy(c.conn, f)
    if err != nil{
        log.Fatal(err.Error())
    }
    
   
    if c.channel == nil{
        c.err(errors.New("[!] you must join the channel first"))
        return
    }

    c.msg(fmt.Sprintf("[*] yo have sent file %s", fileN))
    log.Printf("%s has sent %s",c.nick, fileN)
    if c.channel !=nil{
        c.channel.broadcast(c, fmt.Sprintf("[*] %s has shared %s", c.nick, fileN))
    }

}
func (s *server) listFiles(c *client, args []string){
    var files []string
    for name := range s.files{
        files = append(files, name)
    }

    c.msg(fmt.Sprintf("Available files are: %s", strings.Join(files, ",")))

}

func (s *server) saveFile(c *client, args[]string){
   
    fileN := args[1]
    newFile, err := os.Create(c.nick + "_" + fileN)

    if err != nil{
        panic(err.Error())
    }
    
    defer newFile.Close()
   
    _, err = io.Copy(c.conn, newFile)
    if err != nil{
        log.Fatal(err.Error())
    }

    log.Printf("%s has downloaded %s",c.nick, fileN)
    c.msg(fmt.Sprintf("[!] You have downloaded a file as %s",c.nick + "_" + fileN ))
    
}

func (s *server) exit(c *client, args []string){
    log.Printf("[!] %s has disconnected: %s",c.nick, c.conn.RemoteAddr().String())

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


