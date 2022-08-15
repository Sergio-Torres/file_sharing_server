package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct{
    conn net.Conn
    nick string
    channel *channel
    file *file
    commands chan<- command
}

func (c *client) readInput(){
    for {
        msg, err := bufio.NewReader(c.conn).ReadString('\n')
        if err != nil{
            return
        }
        
        msg = strings.Trim(msg, "\r\n")

        args := strings.Split(msg, " ")
        cmd := strings.TrimSpace(args[0])
        
        switch cmd{
        case "/nick": 
            c.commands <- command{
                id: cmd_nick,
                client: c,
                args: args,
            }
        case "/join":
             c.commands <- command{
                id: cmd_join,
                client: c,
                args: args,
            }
        case "/channels":
            c.commands <- command{
                id: cmd_channels,
                client: c,
                args: args,
            }
        case "/msg":
            c.commands <- command{
                id: cmd_msg,
                client: c,
                args: args,
            }
        case "/exit":
            c.commands <- command{
                id: cmd_exit,
                client: c,
                args: args,
            }
        case "/file":
            c.commands <-command{
                id: cmd_file,
                client: c,
                args: args,
            }
        case "/files":
            c.commands <-command{
                id: cmd_files,
                client: c,
                args: args,
            }
        case "/save":
            c.commands <-command{
                id: cmd_save,
                client: c,
                args: args,
            }

        default:
            c.err(fmt.Errorf("[!]Unknown command: %s", cmd))
        }
    }
}

func (c *client) err(err error){
    c.conn.Write([]byte("[!]ERR: "+ err.Error() + "\n"))
}

func (c *client) msg(msg string){
    c.conn.Write([]byte(">> "+ msg + "\n"))
}


