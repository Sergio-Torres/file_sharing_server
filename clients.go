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


/*
func main(){
    arguments := os.Args
    if len(arguments)==1{
        fmt.Println("Please provide host:port")
        return
    }

    CONNECT := arguments[1]
    c, err := net.Dial("tcp", CONNECT)
    if err != nil{
        fmt.Println(err)
        return
    }

    for{
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(c, text+"\n")

        message, _ := bufio.NewReader(c).ReadString('\n')
        fmt.Print("->: " + message) 
        if strings.TrimSpace(string(text)) == "STOP"{
            fmt.Println("TCP client exiting...")
            return
        }
    }

}

*/
