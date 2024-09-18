# 2024.1-pc-lab5

## Arquitetura utilizada 
Cliente-Servidor

## Integrantes

- Enzo Diniz Vasconcelos - 120211072
- Rafael de Sousa Cavalcante - 121210299
- Gabriel Yuri Ramalho Ferreira - 121210381
- Lucas Emmanuel de Sousa Alves - 121210586


### Primeiramente é preciso buildar cada arquivo.

```bash
go build nomearquivo.go
```

### Comandos

É necessário possuir um diretório com arquivos e colocar seu 
caminho na variável **directory** no client/client.go

Exemplo de publicação para o cliente

Atualizar todos os arquivos de um determinado diretório no server.
É necessário antes rodar o script de client/hash.go para gerar os
hashs dos arquivos. Também é preciso informar o diretório dos arquivos
no script de hash, na variável **directory**.
```bash
./hash && ./client publish
```

- Exemplo de search para o cliente por um determinado hash

```bash
./client search 2
```

- Inicializar o servidor.
```bash
./server
```

