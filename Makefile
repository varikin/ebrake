all: build

build:
	go build -o bin/ebrake

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/ebrake.exe 

clean:
	go clean
	rm -rf bin/

