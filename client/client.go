package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
    directory := "/tmp/dataset"
    command := os.Args[1]
    serverIp := "localhost:5000" // Modificar quando o real ip do servidor for definido
    //serverIp := "150.165.42.145:5000" // Modificar quando o real ip do servidor for definido
    if command == "search" {
        fileHash := os.Args[2]
        search(fileHash, serverIp)
    } else if command == "publish" {
        // operationAndFileHashs := strings.Join(os.Args[2:], " ")
        cmd := exec.Command("../hash/hash")
        _, errCmd := cmd.Output()
        if errCmd != nil {
            fmt.Println(errCmd)
        }
        operationAndFileHashs, errHashs := getFileHashs(directory)
        if errHashs != nil {
            fmt.Println(errHashs)
            os.Exit(1)
        }
        fmt.Println(operationAndFileHashs)
        publish(operationAndFileHashs, serverIp)
    }
}

func getFileHashs(directory string) (string, error) {
    files, errDir := os.ReadDir(directory)
    if errDir != nil {
        return "", errDir
    }
    output := ""
    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".hash") {
            continue
        }
        //originalName := strings.Split(file.Name(), ".hash")[0]
        //if existFile(directory, originalName) {
        //    output += "a,"
        //} else {
        //    output += "r,"
        //}
        hash, errHash := getHashOfFile(directory, file.Name())
        if errHash != nil {
            return "", errHash
        }
        output += hash + " "
    }
    return strings.TrimSpace(output), nil
}

func existFile(directory string, fileName string) bool {
    _, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
    return err == nil
}

func getHashOfFile(directory string, fileName string) (string, error) {
    bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
    if err != nil {
        return "", err
    }
    return string(bytes), nil 
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
	_, err = conn.Write([]byte(clientIp + " publish " + operationAndFileHashs)) // nao precisa mais ser operation e filehashs, apenas todos os hashs que o cliente tem, espacados ex ip + publish + 3 4 5 6
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
