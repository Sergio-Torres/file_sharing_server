package main

import (
	"errors"
	"fmt"
	"io"
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
        c.err(errors.New("[!] you musy join the channel first"))
        return
    }

    c.channel.broadcast(c, c.nick + ":" + strings.Join(args[1:len(args)], " "))
}

func (s *server) file(c *client, args []string){
    /*
    go func (x c.conn){
        defer func(){
            x.Close()
            done <- struct{}{}
        }()
        buffer := make([]byte, BUFFERSIZE)
    }(c.conn)
    f, err := os.Open(fileN)
    //revisar este error
    if err != nil{
        log.Fatal(err.Error())
    }
    pr, pw := io.Pipe()
    w, err := gzip.NewWriterLevel(pw,7)
    if err != nil{
            log.Fatal(err.Error())
    }

    go func(){
        n, err := io.Copy(w, f)
        if err != nil{
            log.Fatal(err.Error())
        }
        w.Close()
        pw.Close()
        log.Printf("se ha copiado y escrito la vaina %s", n)
    }()
    n, err := io.Copy(c.conn, pr)
    if err != nil{
        log.Fatal(err.Error())
    }
    log.Printf("se ha copaido la conexion: %s", n)
    */
    fileN := args[1]
    f, err := os.Open(fileN)
    //revisar este error
    if err != nil{
        log.Fatal(err.Error())
    }

    fileInfo, err := f.Stat()
    if err != nil{
        log.Fatal(err.Error())
    }
    
    fileSize := strconv.FormatInt(fileInfo.Size(),10)
    fileName := fileInfo.Name()
    c.conn.Write([]byte(fileSize))
    c.conn.Write([]byte(fileName))
    buffer := make([]byte, BUFFERSIZE)

    for{
        _, err = f.Read(buffer)
        if err  == io.EOF{
            break
        }
        c.conn.Write(buffer)

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


