package main

import (
	"log"
	"net"

	"github.com/dgageot/docker-machine-dns/zone"
	"github.com/hashicorp/mdns"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		return
	}

	select {}
}

func run() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	zone := zone.NewDockerMachineZone()

	for _, iface := range interfaces {
		if iface.Flags&net.FlagBroadcast == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		if len(addrs) == 0 {
			continue
		}

		log.Printf("Starting server for iface %s", iface.Name)

		_, err = mdns.NewServer(&mdns.Config{Zone: zone, Iface: &iface})

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
