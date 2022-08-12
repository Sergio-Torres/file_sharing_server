package main

type cmdID int

const (
    cmd_nick cmdID = iota
    cmd_join
    cmd_channels
    cmd_msg
    cmd_exit
)

type command struct{
    id cmdID
    client *client
    args []string
}
