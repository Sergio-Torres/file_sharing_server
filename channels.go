package main

import "net"


type channel struct{
    name string
    members map[net.Addr]*client
}
