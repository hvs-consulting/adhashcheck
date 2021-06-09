TARGET_NAME=adhashcheck

all: prepare build-windows build-linux

prepare:
	mkdir -p ./build/

clean:
	rm -rf ./build/

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "./build/${TARGET_NAME}"

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "./build/${TARGET_NAME}.exe"