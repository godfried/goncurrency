package main

import( "fmt"
 "net"
 "bytes"
 "strings"
)
var exit = []byte("/exit") 
type Client struct{
	name string
	snd chan string
	rcv chan string
	quit chan bool
	conn net.Conn
}


/*
Reads data from a connection and writes it to standard output.
*/
func ConnHandler(conn net.Conn, rcv chan string) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		ch <- "Client connection error: "+ err.Error()
		return
	}
	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	client := &Client{name, make(chan string), ch, make(chan bool), conn}
	go ClientReader(client)
	go ClientSender(client)
}

func ClientReader(client *Client){
	buffer := make([]byte, 2048)
	client.snd <- "Client " +  client.name + " has joined"
	for {
		bytesRead, err := client.conn.Read(buffer)
		if err != nil || bytes.Equal(buffer[:len(exit)], exit) {
			break
		}
		client.snd <- client.name+"> "+strings.TrimSpace(string(buffer[:bytesRead]))
	}
	client.snd <- client.name + " has left chat"
}


func ClientSender(client *Client) {
	for {
		select {
		case buffer := <-client.rcv:
			client.conn.Write([]byte(buffer +"\n"))
		case <-client.quit:
			client.snd <- "Client " + client.name +  " quiting"
			client.conn.Close()
			break
		}
	}
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
	ch <- "Server started"
	service := "0.0.0.0:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		panic(err)
	} else {
		listener, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
		if error != nil {
			panic(error)
		} else {
			defer listener.Close()
			for {
				ch <- "Listening for clients"
				conn, err := listener.Accept()
				if err != nil {
					ch <- "Client error: "+error.Error()
				} else {
					go ConnHandler(conn, ch)
				}
			}
		}
	}
}
