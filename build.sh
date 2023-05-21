rm -f dkgo
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o dkgo main.go
scp dkgo root@192.168.217.10:/root
