all: build

build:
	go build -o bin/ebrake

clean:
	go clean
	rm -f bin/

