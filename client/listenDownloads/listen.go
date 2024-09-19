package main

import (
    "fmt"
    "net"
    "os"
    "io"
    "log"
)

func main() {
    directory := "/tmp/dataset"
    listener, errListen := net.Listen("tcp", ":5050")
    if errListen != nil {
        fmt.Println(errListen)
        os.Exit(1)
    }
    defer listener.Close()
    for {
        conn, errAcc := listener.Accept()
        if errAcc != nil {
            fmt.Println(errAcc)
            continue
        }
        
        go handleDownload(conn, directory)
    }
}

func handleDownload(conn net.Conn, directory string) {
    defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Print(err)
		}
		return
	}
	fileName := string(buffer[:n])
    bytes, errRead := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
    if errRead != nil {
        fmt.Println(errRead)
        conn.Write([]byte("file not find"))
        os.Exit(1)
    }
    _, errWrite := conn.Write(bytes)
    if errWrite != nil {
        fmt.Println(errWrite)
    }
}
