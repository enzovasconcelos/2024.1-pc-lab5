package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	command := os.Args[1]
	serverIp := "localhost:8080" // Modificar quando o real ip do servidor for definido
	if command == "search" {
		fileHash := os.Args[2]

		search(fileHash, serverIp)
	} else if command == "publish" {
		operationAndFileHashs := os.Args[2]
		publish(operationAndFileHashs, serverIp)

	}

}

func search(fileHash string, serverIp string) {
	conn, err := net.Dial("tcp", serverIp)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte("search " + fileHash))

	if err != nil {
		fmt.Println("Erro ao enviar dados para o servidor:", err)
		return
	}
}

func publish(operationAndFileHashs string, serverIp string) {
	conn, err := net.Dial("tcp", serverIp)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("publish " + operationAndFileHashs))

	if err != nil {
		fmt.Println("Erro ao enviar dados para o servidor:", err)
		return
	}

}
