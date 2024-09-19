package main

import (
    "fmt"
    "os"
    "net"
    "io"
)

func main() {
    directory := "/tmp/dataset2"
    port := ":5050"
    file := os.Args[1]

    channel1 := make(chan []byte)
    channel2 := make(chan []byte)
    channel3 := make(chan []byte)

    go requestDownload(file, os.Args[2] + port, channel1)
    go requestDownload(file, os.Args[3] + port, channel2)
    go requestDownload(file, os.Args[4] + port, channel3)

    // select
    fmt.Println("esperando select")
    select {
        case fileContent := <- channel1:
            write(fileContent, directory, file)
        case fileContent :=  <- channel2:
            write(fileContent, directory, file)
        case fileContent := <- channel3:
            write(fileContent, directory, file)
    }
}

func write(fileContent []byte, directory string, fileName string) {
    os.WriteFile(directory + "/" + fileName, fileContent, 0777)
}

func requestDownload(file string, ip string, channel chan []byte) {
    conn, err := net.Dial("tcp", ip)
    if err != nil {
        fmt.Print(err)
    }
    defer conn.Close()
    //request := fmt.Sprintf("%s %s", ip, file)
    request := file
    _, errWrite := conn.Write([]byte(request))
    if errWrite != nil {
        fmt.Println(errWrite)
        os.Exit(1)
    }

    response, errResponse := io.ReadAll(conn)
    if err != nil {
        fmt.Println(errResponse)
        os.Exit(1)
    }
    channel <- response
    fmt.Println(string(response))
}
