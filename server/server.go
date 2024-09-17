package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	channelDiffs := make(chan []string)
	fileHashs := make(map[string][]string)

	go listenDiffs(channelDiffs, fileHashs)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn, channelDiffs)
	}
}

func handleConnection(conn net.Conn, channelDiffs chan []string) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Print(err)
		}
		return
	}
	data := strings.Split(string(buffer[:n]), " ")
	if len(data) < 2 {
		log.Print("descartando comando inválido")
		return
	}
	switch command := data[1]; command {
	case "publish":
		log.Print("entrou no publish")
		handlePublish(conn, data, channelDiffs)
	case "find":
		log.Print("entrou no find")
		// Implementar handleFind aqui
	default:
		log.Print("descartando comando inválido")
	}
}

func handlePublish(conn net.Conn, data []string, channelDiffs chan []string) {
	for _, diff := range data[2:] {
		channelDiffs <- []string{data[0], diff}
	}

	fmt.Println("teste")
	// Envia uma mensagem de confirmação ao cliente
	_, err := conn.Write([]byte("Itens adicionados com sucesso"))
	if err != nil {
		log.Print("Erro ao enviar confirmação ao cliente:", err)
	}
}

func removeElem(list []string, elem string) []string {
	newList := []string{}
	for _, current := range list {
		if current != elem {
			newList = append(newList, current)
		}
	}
	return newList
}

func listenDiffs(channelDiffs chan []string, fileHashs map[string][]string) {
	for {
		diff := <-channelDiffs
		ipAddress := diff[0]
		command := strings.Split(diff[1], ",")
		fileHash := command[1]
		switch command[0] {
		case "a":
			if !contains(fileHashs[fileHash], ipAddress) {
				fileHashs[fileHash] = append(fileHashs[fileHash], ipAddress)
				log.Print("file added")
			} else {
				log.Print("IP já presente no array")
			}
		case "r":
			fileHashs[fileHash] = removeElem(fileHashs[fileHash], ipAddress)
			log.Print("file removed")
		}
		fmt.Println(fileHashs)
	}
}

func contains(list []string, elem string) bool {
	for _, item := range list {
		if item == elem {
			return true
		}
	}
	return false
}
