# Binary name
BINARY=goethspeedup
# Builds the project
build:
		go build -o bin/${BINARY} cmd/main.go
# Installs our project: copies binaries
install:
		go install
release:
		# Clean
		go clean
		rm -rf bin/*
		# Build for mac
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY} cmd/main.go
		tar czvf bin/${BINARY}-darwin-amd64-${VERSION}.tar.gz bin/${BINARY}
		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o bin/${BINARY} cmd/main.go
		tar czvf bin/${BINARY}-linux-amd64-${VERSION}.tar.gz bin/${BINARY}
		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/${BINARY} cmd/main.go
		tar czvf bin/${BINARY}-windows-amd64-${VERSION}.tar.gz bin/${BINARY}
		go clean
# Cleans our projects: deletes binaries
clean:
		go clean

.PHONY:  clean build