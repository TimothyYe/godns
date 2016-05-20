CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o godns-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o godns-i386
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o godns-mac
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o godns-win

echo 'compile success!'