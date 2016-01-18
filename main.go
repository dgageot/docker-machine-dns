package main

import (
	"flag"
	"log"
	"net"

	"github.com/dgageot/docker-machine-dns/zone"

	"github.com/hashicorp/mdns"
)

func main() {
	interfaceName := flag.String("interfaceName", "vboxnet0", "Network interface name")

	flag.Parse()

	if err := runServer(*interfaceName); err != nil {
		log.Fatal(err)
	}
}

func runServer(interfaceName string) error {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	server, err := mdns.NewServer(&mdns.Config{
		Zone:  &zone.DockerMachineZone{},
		Iface: iface,
	})
	if err != nil {
		return err
	}

	defer server.Shutdown()

	select {}
}
