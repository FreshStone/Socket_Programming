package main

import (
  "fmt"
  "net"
  "os"
  "io"
  "bytes"
  "strconv"
)

const Buffer_Size = 1024

func main(){
  Start_Ftserver()
}

func Start_Ftserver(){
    listener, err := net.Listen("tcp", ":2000")
    if err != nil{
      fmt.Println(err)
      os.Exit(1)
    }
    for {
      conn, err := listener.Accept()
      fmt.Println("new client added")
      if err != nil{
        fmt.Println(err)
        os.Exit(1)
      }
      go handleconnection(conn)
    }
  }

  func handleconnection(conn net.Conn){
    defer conn.Close()
    file := make([]byte, Buffer_Size)
    conn.Read(file)
    filename := string(bytes.Trim(file, "\x00"))
    f, err := os.Open(filename)
    if err != nil {
	    fmt.Println(err)
	    return
    }
    defer f.Close()
    stat, _ := f.Stat()
    conn.Write([]byte(strconv.FormatInt(stat.Size(), 10)))
    buffer := make([]byte, Buffer_Size)
    var bytes_read int64 = 0
    fmt.Println("sending " + filename)
    for {
      n, err := f.ReadAt(buffer, bytes_read)
      bytes_read += int64(n)
      if n == Buffer_Size{
        conn.Write(buffer)
      }else {
        conn.Write(buffer[:n])
      }
      if err == io.EOF{
        break
      }
    }
    fmt.Printf("%v bytes sent\n", stat.Size())
  }
