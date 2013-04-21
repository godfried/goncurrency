package main

// Imports required packages
import( "fmt"
 "net"
 "bytes"
 "strings"
)

/*
Reads data from a connection and writes it to standard output.
*/
func ConnHandler(conn net.Conn, ch chan string) {
	buffer := make([]byte, 2048)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		ch <- "Client connection error: "+ err.Error()
		return
	}

	name := strings.TrimSpace(string(buffer[0:bytesRead]))
	ch <- "Client " +  name + " has joined"
	exit := []byte("exit") 
	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil || bytes.Equal(buffer[:len(exit)], exit) {
			break
		}
		ch <- name+"> "+strings.TrimSpace(string(buffer))
		for i := 0; i < 2048; i++ {
			buffer[i] = 0x00
		}
	}
	ch <- name + " has left chat"
	conn.Close()
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
