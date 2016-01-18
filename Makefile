build: docker-machine-dns
	
docker-machine-dns: main.go zone/zone.go
	GO15VENDOREXPERIMENT=1 go build .
	
clean:
	rm docker-machine-dns