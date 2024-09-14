package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	command := os.Args[1]

	if command == "search" {
		fileHash := os.Args[2]
		serverIp := "localhost:8080" // Modificar quando o real ip do servidor for definido
		search(fileHash, serverIp)
	}
}

func search(fileHash string, serverIp string) {
	conn, err := net.Dial("tcp", serverIp)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write([]byte(fileHash + "\n"))

	if err != nil {
		fmt.Println(err)
		return
	}

}
