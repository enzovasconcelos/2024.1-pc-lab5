package main

import (
	"log"
	"net"
    "io"
    "strings"
    "fmt"
)

func main() {
    fileHashs := make(map [string][]string)
    channelDiffs := make(chan []string)
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
    go listenDiffs(channelDiffs, fileHashs)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
        request, errRead := io.ReadAll(conn)
        if errRead != nil {
            log.Print(errRead)
        }
        data := strings.Split(string(request), " ")
        switch command := data[1]; command {
            case "publish":
                log.Print("entrou no publish")
                go handlePublish(data, channelDiffs)
            case "find":
                // TODO: implementar aqui
                log.Print("entrou no find")
                // go handleFind(conn, data, fileHashs)
            default:
                log.Print("descartando comando inválido")
        }
	}
}

func listenDiffs(channelDiffs chan []string, fileHashs map[string][]string) {
    for {
        diff, _ := <- channelDiffs
        ipAddress := diff[0] 
        command := strings.Split(diff[1], ",")
        fileHash := command[1]
        switch command[0] {
            case "a":
                fileHashs[fileHash] = append(fileHashs[fileHash], ipAddress)
                log.Print("file added")
            case "r":
                fileHashs[fileHash] = removeElem(fileHashs[fileHash], ipAddress)
                log.Print("file removed")
        }
        fmt.Println(fileHashs)
    }
}

func removeElem(list []string, elem string) [] string {
    newList := []string {}
    for _, current := range list {
        if current != elem {
            newList = append(newList, current)
        }
    } 
    return newList
}

func handlePublish(data []string, channelDiffs chan []string) {
    for _, diff := range(data[2:]) {
        channelDiffs <- []string { data[0], diff }
    }
}
