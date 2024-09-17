package main

import (
    "fmt"
    "os"
    "io"
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
    
    for _, file := range files {
        go calculateHash(file.Name(), directory, join)
    }
    
    for i := 0; i < len(files); i++ {
        <- join
    }
    fmt.Println("Done :)")
}

func calculateHash(fileName string, directory string, join chan bool) {
    bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    var sum int
    for _, byteReaded := range bytes {
        fmt.Println(byteReaded)
        sum += int(byteReaded)
    }

    filePath := fmt.Sprintf("%s/%s.hash", directory, fileName)
    fmt.Println(fmt.Sprintf("writing %d hash", sum))
    fmt.Println(string(sum))
    var buf [] byte
    binary.Encoder(buf, uint8, sum) 
    errWrite := os.WriteFile(filePath, buf, 0777)
    if errWrite != nil {
        fmt.Println(errWrite)
        os.Exit(1)
    }
    join <- true
}
