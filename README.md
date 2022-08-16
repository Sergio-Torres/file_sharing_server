### `How to start`
    
- Clone te repository
- Open your terminal and in the directory run the command: `go build . `
- To start the server run: `./file_sharing_server`

Once executed, the terminal should show the following message:
```
2022/08/15 17:57:53 started server on :8888
```
This means that the server is started and ready to receive clients and notify of any event that occurs.

### `How to connect to the server (the client)`

You can connect the number of clients you want to the server, for this you must open a terminal (one terminal for each client) and execute the following command:
`telnet localhost 8888`

In the server terminal it should show a message like:
```
2022/08/15 18:04:51 New client has connected: 127.0.0.1:40182

```
### `Commands`
- `/nick` This command is used to give you a name on the server, for example `/nick Juan`
- `/join`This command is used to create channels or connect to an existing channel, for example `/join #new`
- `/channels` This command displays all available channels.
- `/msg` This command is used to send messages to other clients in the channels, for example `/msg hello everyone`
- `/file` This message is used to send files to a channel, for example `/file file.txt`
- `/files` This command displays all available channels.
- `/save` This command is used to save a file, for example `/save archivo.txt`. A file with the client name and the original name will be created and saved in the directory, `Juan_archivo.txt` for example.
- `/exit` This command is used to leave the server.

The important thing to know is that although it is not necessary for the client to connect to the directory 
or repository where the files are located, the file to be shared must be in the directory.

    

 
