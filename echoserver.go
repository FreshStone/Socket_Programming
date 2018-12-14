package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"strings"
	"os/signal"   //for intercepting interrupt
)

func main(){
	if len(os.Args) != 3{
		fmt.Printf("Usage: %s startmode('client' or 'server') serviceport\n", os.Args[0])
		os.Exit(1)
	}
	mode := os.Args[1]
	port := os.Args[2]
	if mode == "server"{
		StartServer(port)
	}else if mode == "client"{
		StartClient(port)
	}else{
		fmt.Println("Enter correct mode")
		os.Exit(1)
	}
}


func StartClient(port string){
	conn, _ := net.Dial("tcp", "127.0.0.1:"+port)
	defer conn.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<-c
		conn.Write([]byte("quit"))
		os.Exit(1)
	}()

	for{
		server_msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(strings.TrimRight(server_msg, "\r\n"))
		s, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		client_msg := strings.TrimRight(s, "\r\n")
		conn.Write([]byte(client_msg + "\n"))
		if client_msg == "quit"{
                        os.Exit(0)
                }
	}
}


func StartServer(port string){
	listner, _ := net.Listen("tcp", "127.0.0.1:"+port)
	conn, _ := listner.Accept()
	defer conn.Close()
	conn.Write([]byte("Enter text to continue or 'quit' for exit" + "\n"))
	reader := bufio.NewReader(conn)
	for{
		s, _ := reader.ReadString('\n')
		client_msg := strings.TrimRight(s, "\r\n")
		if client_msg == "quit"{
			os.Exit(0)
		}
		fmt.Println(client_msg)
		fmt.Fprintf(conn, "'%s' recieved\n", client_msg)
	}
}

