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
	Snd chan string
	Rcv chan string
	Quit chan bool
	Conn net.Conn
}


/*
 Handles a new client connection. Reads client username and spawns routines for 
sending data to and receiving data from other clients.
*/
func ConnHandler(conn net.Conn, snd chan string, connect chan *ClientData) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		conn.Write([]byte("Client connection error: "+ err.Error()))
		conn.Close()
		return
	}
	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	client := &Client{name, snd, make(chan string), make(chan bool), conn}
	connect <- &ClientData{client.Name, client.Rcv, true}
	success := <- client.Rcv
	if success == ERROR{
		client.Conn.Write([]byte("Username already in use"))
		client.Conn.Close()
		return
	} 
	go ClientReader(client)
	go ClientSender(client, connect)
}

/*
Reads client send data. Sends the data to the server routines for handling client input.
*/
func ClientReader(client *Client){
	buffer := make([]byte, 2048)
	client.Snd <- "Client " +  client.Name + " has joined"
	exit := []byte("/exit") 
	for {
		bytesRead, err := client.Conn.Read(buffer)
		if err != nil || bytes.Equal(buffer[:len(exit)], exit) {
			break
		}
		client.Snd <- client.Name+"> "+strings.TrimSpace(string(buffer[:bytesRead]))
	}
	client.Quit <- true
}


/*
Listens for data and exit signal sent to this client.
*/
func ClientSender(client *Client, connect chan *ClientData) {
	for {
		select {
		case buffer := <-client.Rcv:
			client.Conn.Write([]byte(buffer +"\n"))
		case <-client.Quit:
			connect <- &ClientData{client.Name, client.Rcv, false}
			client.Snd <- "Client " + client.Name +  " has left chat."
			client.Conn.Close()
			break
		}
	}
}


type ClientData struct{
	Name string
	Chan chan string
	Connect bool
}

const(
	OK = "ok"
ERROR = "err"
)
/*
Read input from channel and writes to standard output
*/
func IOHandler(msgChan chan string, connected chan *ClientData){
	listeners := make(map[string] chan string)
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg)
			for _, ch := range(listeners) {			
				ch <- msg
			}
		case data := <- connected:
			if data.Connect{
				if _, ok := listeners[data.Name]; !ok{
					data.Chan <- OK
					listeners[data.Name] = data.Chan
				} else{
					data.Chan <- ERROR
				}
			}else{
				delete(listeners, data.Name)
			}
			
		}
	}
}

/*
Listens for new connections and spawns goroutines to handle them.
*/
func main() {
	dataChan := make(chan string)
	connected := make(chan *ClientData)
	go IOHandler(dataChan, connected)
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
					go ConnHandler(conn, dataChan, connected)
				}
			}
		}
	}
}