package main

type cmdID int

const (
    cmd_nick cmdID = iota
    cmd_join
    cmd_channels
    cmd_msg
    cmd_exit
    cmd_file
    cmd_files
    cmd_save
)

type command struct{
    id cmdID
    client *client
    args []string
}
