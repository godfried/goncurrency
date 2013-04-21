package main

import( 
 "net"
 "bytes"
 "strings"
"sync"
"fmt"
)
var exit = []byte("/exit") 
var clients = make(map[string]*Client)
var clientM = new(sync.Mutex)	
type Client struct{
	name string
	snd chan string
	rcv chan string
	quit chan bool
	conn net.Conn
}

func Add(client *Client) (ret bool){
//	clientM.Lock()
	if clients[client.name] == nil{
		clients[client.name] = client
		ret = true
	}
//	clientM.Unlock()
	return ret
}

func Remove(client *Client){
	clientM.Lock()
	delete(clients, client.name)
	clientM.Unlock()
}

/*
Reads data from a connection and writes it to standard output.
*/
func ConnHandler(conn net.Conn, ch chan string) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		conn.Write([]byte("Client connection error: "+ err.Error()))
		return
	}
	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	client := &Client{name, make(chan string), ch, make(chan bool), conn}
	if Add(client){
		go ClientReader(client)
		go ClientSender(client)
	} else{
		conn.Write([]byte("Username taken"))
	}
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
			fmt.Println("sender", buffer)
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
func IOHandler(dataChan chan string){
	for data := range dataChan{
		fmt.Println(data)
		//clientM.Lock()
		for _, client := range clients {
			client.rcv <-data
			fmt.Println("io", data)
		}	
		//clientM.Unlock()
	}
}




/*
Listens for new connections and spawns goroutines to handle them.
*/
func main() {
	ch := make(chan string)
	go IOHandler(ch)
	fmt.Println("Server started")
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
				fmt.Println("Listening for clients")
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Client error: "+error.Error())
				} else {
					go ConnHandler(conn, ch)
				}
			}
		}
	}
}
