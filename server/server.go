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
		go handleConnection(conn, channelDiffs, fileHashs)
	}
}

func handleConnection(conn net.Conn, channelDiffs chan []string, fileHashs map[string][]string) {
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
		handlePublish(conn, data, channelDiffs, fileHashs)
	case "find":
		log.Print("entrou no find")
		handleFind(conn, data, fileHashs)
	default:
		log.Print("descartando comando inválido")
	}
}

func handleFind(conn net.Conn, data []string, fileHashs map[string][]string) {
	if len(data) < 3 {
		_, err := conn.Write([]byte("Comando find inválido"))
		if err != nil {
			log.Print("Erro ao enviar mensagem de comando inválido ao cliente:", err)
		}
		return
	}

	hash := data[2]
	ips, found := fileHashs[hash]
	if !found {
		_, err := conn.Write([]byte("Hash não encontrado"))
		if err != nil {
			log.Print("Erro ao enviar mensagem de hash não encontrado ao cliente:", err)
		}
		return
	}

	response := strings.Join(ips, ",")
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Print("Erro ao enviar resposta ao cliente:", err)
	}
}

func handlePublish(conn net.Conn, data []string, channelDiffs chan []string, fileHashs map[string][]string) {
	// percorrer cada hash do mapa e verificar se o ip do cliente está presente
	diffs := []string{}
	for hash, ips := range fileHashs {
        ipPresente := false
        for _, ip := range ips {
            if ip == data[0] {
                ipPresente = true
                break
            }
        }

		if ipPresente && contains(data[2:], hash) {
			diffs = append(diffs, "a,"+hash)
		} else if !ipPresente && contains(data[2:], hash) {
			diffs = append(diffs, "a,"+hash)
		} else if ipPresente && !contains(data[2:], hash) {
			diffs = append(diffs, "r,"+hash)
		}
    }
	
	// Verificar se há hashes em data[2:] que não estão no mapa fileHashs
    for _, hash := range data[2:] {
        if _, exists := fileHashs[hash]; !exists {
            // Adicionar lógica para lidar com hashes que não existem no mapa
            diffs = append(diffs, "a,"+hash)
        }
    }
	
	for _, diff := range diffs {
		channelDiffs <- []string{data[0], diff}
	}

	var err interface{}
	// Envia uma mensagem de confirmação ao cliente
	_, err = conn.Write([]byte("Itens atualizados com sucesso"))

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
			if contains(fileHashs[fileHash], ipAddress) {
				fileHashs[fileHash] = removeElem(fileHashs[fileHash], ipAddress)
				if len(fileHashs[fileHash]) == 0 {
					delete(fileHashs, fileHash)
				}
				log.Print("file removed")
			} else {
				log.Print("IP não está presente no array")
			}
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
