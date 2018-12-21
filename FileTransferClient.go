package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "io"
  "strings"
  "strconv"
  "bytes"
)

const Buffer_Size = 1024
func main(){
  Start_Ftclient()
}

func Start_Ftclient(){
  conn, err := net.Dial("tcp", ":2000")
  if err != nil{
    fmt.Println(err)
    os.Exit(1)
  }
  defer conn.Close()
  fmt.Println("Enter filename...")
  filename, _ := bufio.NewReader(os.Stdin).ReadString('\n')
  conn.Write([]byte(strings.TrimSpace(filename)))
  filesize := make([]byte, Buffer_Size)
  _, e := conn.Read(filesize)
  if e != nil{
	  fmt.Println(err)
	  conn.Close()
	  return
  }
  filesize_int, _ := strconv.ParseInt(string(bytes.Trim(filesize, "\x00")), 10, 64)
  f, _ := os.Create(strings.TrimSpace(filename))
  defer f.Close()
  var received_bytes int64
  for {
    if filesize_int-received_bytes > Buffer_Size{
      io.CopyN(f, conn, Buffer_Size)
    }else{
      io.CopyN(f, conn, filesize_int-received_bytes)
      conn.Read(make([]byte, (received_bytes+Buffer_Size)-filesize_int)) //reading extra garbage bytes
      break
    }
    received_bytes += Buffer_Size
  }
  fmt.Println(strings.TrimSpace(filename) + "received")
}
