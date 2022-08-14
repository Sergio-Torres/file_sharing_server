package main

import "net"

type channel struct{
    name string
    members map[net.Addr]*client
}

type file struct{
    name string
    members map[net.Addr]*client
}
func (ch *channel) broadcast(sender *client, msg string){
    for addr, m := range ch.members{
        if addr != sender.conn.RemoteAddr(){
            m.msg(msg)
        }

    }
} 

/*func (ch *channel) fileBroadcast(sender *client, file string){
    for addr, f := range ch.members{
        if addr != sender.conn.RemoteAddr(){
            f.file(file)
        }
    }
}*/
