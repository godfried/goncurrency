package main

import( 
 "net"
 "bytes"
 "strings"
"fmt"
"os"
)

/*
Represents a single client connection
*/
type Client struct{
	Name string
	Ch chan string
	Quit chan bool
	Conn net.Conn
}

/*
 Handles a new client connection. Reads client username and spawns routines for 
sending data to and receiving data from other clients.
*/
func ConnHandler(conn net.Conn, ch chan string) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		conn.Write([]byte("Client connection error: "+ err.Error()))
		conn.Close()
		return
	}
	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	client := &Client{name, ch, make(chan bool), conn}
	go ClientReader(client)
	go ClientSender(client)
}

/*
Reads client send data. Sends the data to the server routines for handling client input.
*/
func ClientReader(client *Client){
	buffer := make([]byte, 2048)
	client.Ch <- "Client " +  client.Name + " has joined"
	exit := []byte("/exit") 
	for {
		bytesRead, err := client.Conn.Read(buffer)
		if err != nil || bytes.Equal(buffer[:len(exit)], exit) {
			break
		}
		client.Ch <- client.Name+"> "+strings.TrimSpace(string(buffer[:bytesRead]))
	}
	client.Quit <- true
}


/*
Listens for data and exit signal sent to this client.
*/
func ClientSender(client *Client) {
	for {
		select {
		case buffer := <-client.Ch:
			client.Conn.Write([]byte(buffer +"\n"))
		case <-client.Quit:
			client.Ch <- "Client " + client.Name +  " has left chat."
			client.Conn.Close()
			break
		}
	}
}


/*
Listens for new connections and spawns goroutines to handle them.
*/
func main() {
	fmt.Println("Server started")
	service := "localhost:3000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("Could not resolve: ", service)
		os.Exit(1)
	} else {
		listener, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
		if err != nil {
			fmt.Println("Could not listen on: ", tcpAddr)
			os.Exit(1)
		} else {
			ch := make(chan string)
			defer listener.Close()
			for {
				fmt.Println("Listening for clients")
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Client error: ", err)
				} else {
					//Create routine for each connected client
					go ConnHandler(conn, ch)
				}
			}
		}
	}
}