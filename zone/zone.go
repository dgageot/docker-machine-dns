package zone

import (
	"errors"
	"net"

	"fmt"

	"strings"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"
)

var (
	errOnlyForVirtualBox = errors.New("Only for Virtualbox for now")
)

type DockerMachineZone struct{}

func (dm *DockerMachineZone) Records(q dns.Question) []dns.RR {
	// Not for us
	if strings.HasPrefix(q.Name, "_") {
		return nil
	}

	// Not for us
	if !strings.HasSuffix(q.Name, ".local.") {
		return nil
	}

	machineName := q.Name[0 : len(q.Name)-len(".local.")]

	ip, err := findIP(machineName)
	if err != nil {
		return nil
	}

	service, err := mdns.NewMDNSService("DockerMachine", "_ssh._tcp", "local.", q.Name, 22, []net.IP{ip}, []string{"DockerMachine " + machineName})
	if err != nil {
		return nil
	}

	return service.Records(q)
}

func findIP(name string) (net.IP, error) {
	fmt.Println("Looking for virtualbox machine", name)

	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	defer api.Close()

	machine, err := api.Load(name)
	if err != nil {
		return nil, err
	}

	fmt.Println(machine.Driver.DriverName())

	if machine.Driver.DriverName() != "virtualbox" {
		return nil, errOnlyForVirtualBox
	}

	ip, err := machine.Driver.GetIP()
	if err != nil {
		return nil, err
	}

	return net.ParseIP(ip), nil
}
