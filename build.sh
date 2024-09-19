cd client
go build client.go 

cd hash
go build hash.go

cd ../../server
go build server.go

cd ../client/download/
go build download.go

cd ../listenDownloads
go buid listen.go
