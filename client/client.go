package main

import (
    "fmt"
    "io"
    "net"
    "os"
    "strings"
)

func main() {
    command := os.Args[1]
    serverIp := "localhost:5000" // Modificar quando o real ip do servidor for definido
    if command == "search" {
        fileHash := os.Args[2]
        search(fileHash, serverIp)
    } else if command == "publish" {
        operationAndFileHashs := strings.Join(os.Args[2:], " ")
        publish(operationAndFileHashs, serverIp)
    }
}

func search(fileHash string, serverIp string) {

	clientIp, err := getClientIP()
    if err != nil {
        fmt.Println("Erro ao obter o IP do cliente:", err)
        return
    }

    conn, err := net.Dial("tcp", serverIp)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    _, err = conn.Write([]byte(clientIp + " find " + fileHash))
    if err != nil {
        fmt.Println("Erro ao enviar dados para o servidor:", err)
        return
    }

	fmt.Println("Aguardando resposta do servidor...")

    response, err := io.ReadAll(conn)
    if err != nil {
        fmt.Println("Erro ao ler resposta do servidor:", err)
        return
    }

    fmt.Println("Resposta do servidor:", string(response))
}

func publish(operationAndFileHashs string, serverIp string) {
    clientIp, err := getClientIP()
    if err != nil {
        fmt.Println("Erro ao obter o IP do cliente:", err)
        return
    }

    conn, err := net.Dial("tcp", serverIp)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    fmt.Println("IP do cliente:", clientIp)

    _, err = conn.Write([]byte(clientIp + " publish " + operationAndFileHashs))
    if err != nil {
        fmt.Println("Erro ao enviar dados para o servidor:", err)
        return
    }

    fmt.Println("Dados enviados para o servidor")
    // Lê a resposta do servidor
    fmt.Println("Aguardando resposta do servidor...")

    response, err := io.ReadAll(conn)
    if err != nil {
        fmt.Println("Erro ao ler resposta do servidor:", err)
        return
    }

    fmt.Println("Resposta do servidor:", string(response))
}

func getClientIP() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        // Verifica se o endereço é do tipo IP e não é um loopback
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String(), nil
            }
        }
    }

    return "", fmt.Errorf("não foi possível obter o IP do cliente")
}