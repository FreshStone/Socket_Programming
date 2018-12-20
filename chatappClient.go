//based on the blog "https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/"
package main

import (
      "fmt"
      "os"
      "net"
      "bufio"
      "strings"
    )

func main(){
    StartClient()
}

func StartClient() {
    fmt.Println("Starting client...Dialing on port - 1200")
    connection, error := net.Dial("tcp", "localhost:1200")
    if error != nil {
        fmt.Println(error)
    }
    go receive(&connection)
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
	connection.Write([]byte(strings.TrimRight(message, "\n")))
    }
}

func receive(conn *net.Conn){
        for{
                message := make([]byte, 4096)
		            length, err := (*conn).Read(message)
                if err != nil{
			                  (*conn).Close()
			                   break
                }
                if length > 0 {
                        fmt.Println("RECEIVED:" + string(message))
                }
        }
}
