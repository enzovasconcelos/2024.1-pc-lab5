package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type ResultCalculate struct {
    fileName string
    sum int
}

func main() {
    // diretório onde estão os arquivos
    directory := "/tmp/dataset"
    join := make(chan bool)

    files, errRead := os.ReadDir(directory)
    if errRead != nil {
        fmt.Println(errRead)
        os.Exit(1)
    }

    numFiles := len(files)
    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".hash") {
            numFiles -= 1
            continue
        }
        go calculateHash(file.Name(), directory, join)
    }
    
    for i := 0; i < numFiles; i++ {
        <- join
    }
    fmt.Println("Done :)")
}

func calculateHash(fileName string, directory string, join chan bool) {
    _, errHash := os.ReadFile(fmt.Sprintf("%s/%s.hash", directory, fileName))
    if errHash == nil {
        //já existe o arquivo
        os.Exit(0)
    }
    bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    var sum int
    for _, byteReaded := range bytes {
        sum += int(byteReaded)
    }

    filePath := fmt.Sprintf("%s/%s.hash", directory, fileName)
    fmt.Println(fmt.Sprintf("saving hash in %s...", filePath))
    errWrite := os.WriteFile(filePath, []byte(strconv.Itoa(sum)), 0777)
    if errWrite != nil {
        fmt.Println(errWrite)
        os.Exit(1)
    }
    join <- true
}
