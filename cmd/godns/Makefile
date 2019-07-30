# Binary name
BINARY=godns
# Builds the project
build:
		GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
# Installs our project: copies binaries
install:
		GO111MODULE=on go install
image:
		# Build docker image
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		docker build -t timothyye/godns:${VERSION} .
		docker tag timothyye/godns:${VERSION} timothyye/godns:latest
		docker push timothyye/godns:${VERSION}
		docker push timothyye/godns:latest
release:
		# Clean
		go clean
		rm -rf *.gz
		# Build for mac
		GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
		# Build for arm
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-arm-${VERSION}.tar.gz ./${BINARY}
		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ${BINARY}.exe -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
		make image
# Cleans our projects: deletes binaries
clean:
		go clean
		rm -rf *.gz

.PHONY:  clean build
