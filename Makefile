export GO15VENDOREXPERIMENT = 1

.DEFAULT_GOAL := build

build: docker-machine-dns

docker-machine-dns: main.go zone/zone.go
	go build .

deps:
	godep save

clean:
	rm docker-machine-dns

