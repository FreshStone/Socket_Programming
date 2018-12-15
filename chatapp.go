//based on the blog "https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/"
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)


type ClientManager struct {
    clients map[*net.Conn]bool
    broadcast  chan []byte
    register chan *net.Conn
    unregister chan *net.Conn
}

func main() {
    if len(os.Args) != 2{
            fmt.Printf("Usage: %s 'mode'\nmode: client/server\n", os.Args[0])
            os.Exit(1)
    }
    if os.Args[1] == "server" {
        StartServer()
    } else if os.Args[1] == "client" {
        StartClient()
    }else {
            fmt.Println("Incorrect Usage")
            os.Exit(1)
    }
}

func (manager *ClientManager) start() {
    for {
        select {
        case connection := <-manager.register:
            manager.clients[connection] = true
            fmt.Println("Added new connection!")
        case connection := <-manager.unregister:
            if _, ok := manager.clients[connection]; ok {
                delete(manager.clients, connection)
                fmt.Println("A connection has terminated!")
            }
        case message := <-manager.broadcast:
            for connection := range manager.clients {
		_, err := (*connection).Write(message)
		if err != nil {
			delete(manager.clients, connection)
			(*connection).Close()
		}
	    }
        }
    }
}

func receive(conn *net.Conn, manager *ClientManager){
        for{
                message := make([]byte, 4096)
		length, err := (*conn).Read(message)
                if err != nil{
                        if manager != nil{
                                manager.unregister <- conn
                        }
			(*conn).Close()
			break
                }
                if length > 0 {
                        fmt.Println("RECEIVED:" + string(message))
                        if manager != nil{
                                manager.broadcast <- message
                        }
                }
        }
}

func StartServer() {
    fmt.Println("Starting server...")
    listener, error := net.Listen("tcp", ":1200")
    if error != nil {
        fmt.Println(error)
    }
    manager := ClientManager{
        clients:    make(map[*net.Conn]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *net.Conn),
        unregister: make(chan *net.Conn),
    }
    go manager.start()
    for {
        connection, error := listener.Accept()
        if error != nil {
            fmt.Println(error)
        }
	manager.register <- &connection
	go receive(&connection, &manager)
    }
}

func StartClient() {
    fmt.Println("Starting client...")
    connection, error := net.Dial("tcp", "localhost:1200")
    if error != nil {
        fmt.Println(error)
    }
    go receive(&connection, nil)
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
	connection.Write([]byte(strings.TrimRight(message, "\n")))
    }
}

