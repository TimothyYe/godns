# Binary name
BINARY=godns
# Builds the project
build:
		GO111MODULE=on go build -ldflags "-X main.Version=${VERSION}" -o ${BINARY} cmd/godns/godns.go 
# Installs our project: copies binaries
install:
		GO111MODULE=on go install
image:
		# Build docker image
		go clean
		docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t timothyye/godns:${VERSION} . --push
		docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t timothyye/godns:latest . --push
release:
		# Clean
		go clean
		rm -rf *.gz
		# Build for mac
		GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}" cmd/godns/godns.go
		tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}" cmd/godns/godns.go 
		tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
		# Build for arm
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}" cmd/godns/godns.go 
		tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}" cmd/godns/godns.go 
		tar czvf ${BINARY}-arm-${VERSION}.tar.gz ./${BINARY}
		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}" cmd/godns/godns.go
		tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
		make image
# Cleans our projects: deletes binaries
clean:
		go clean
		rm -rf ./godns
		rm -rf ./godns.exe
		rm -rf *.gz

.PHONY:  clean build
