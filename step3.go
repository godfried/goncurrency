package main

import( "fmt"
 "net"
 "bytes"
 "strings"
"os"
)

 
/*
Represents a single client connection
*/
type Client struct{
	Name string
	Ch chan string
	Conn net.Conn
}


/*
 Handles a new client connection. Reads client username and spawns routine for 
sending client data.
*/
func ConnHandler(conn net.Conn, ch chan string) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		ch <- "Client connection error: "+ err.Error()
		conn.Close()
		return
	}
	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	client := &Client{name, ch, conn}
	go ClientReader(client)
}

/*
Reads data from the client. Sends the data to the server routines for handling client input. 
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
	client.Ch <- client.Name + " has left chat"
	client.Conn.Close()
}


/*
Read input from channel and writes to standard output
*/
func IOHandler(ch chan string){
	for{
		read := <- ch
		fmt.Println(read)
	}
}


/*
Listens for new connections and spawns goroutines to handle them.
*/
func main() {
	ch := make(chan string)
	go IOHandler(ch)
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
