# Binary name
BINARY=godns
# Builds the project
build:
		go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
# Installs our project: copies binaries
install:
		go install
release:
		# Clean
		go clean
		rm -rf *.gz
		# Build for mac
		go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
		# Build for arm
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
		go clean
# Cleans our projects: deletes binaries
clean:
		go clean
		rm -rf *.gz

.PHONY:  clean build
