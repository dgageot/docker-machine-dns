BIN := docker-machine-dns

export GO15VENDOREXPERIMENT = 1

.DEFAULT_GOAL := build

run: build
	./$(BIN)

build: $(BIN)

$(BIN): main.go zone/zone.go
	go build .

deps:
	godep save

clean:
	rm -f $(BIN)

