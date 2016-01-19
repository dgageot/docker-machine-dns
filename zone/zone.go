package zone

import (
	"net"
	"time"

	"log"

	"strings"

	"sync"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"
)

// dockerMachineZone is a dns zone that reads its information from Docker Machine.
type dockerMachineZone struct {
	Ttl      uint32
	services map[string]*serviceEntry
	lock     sync.Locker
}

type serviceEntry struct {
	expiryDate time.Time
	service    *mdns.MDNSService
}

// NewDockerMachineZone creates a dns zone that reads its information from Docker Machine.
func NewDockerMachineZone() mdns.Zone {
	return &dockerMachineZone{
		Ttl:      60,
		services: map[string]*serviceEntry{},
		lock:     &sync.Mutex{},
	}
}

// Records returns DNS records in response to a DNS question.
func (dm *dockerMachineZone) Records(q dns.Question) []dns.RR {
	// Not for us
	if strings.HasPrefix(q.Name, "_") || !strings.HasSuffix(q.Name, ".local.") {
		return nil
	}

	service, err := dm.findService(q.Name)
	if err != nil {
		log.Println("Error looking for Docker Machine host", err)
		return nil
	}

	records := SetTTL(service.Records(q), dm.Ttl)

	log.Println(records)

	return records
}

// findService finds a mdns service by its fully qualified name.
func (dm *dockerMachineZone) findService(fqn string) (*mdns.MDNSService, error) {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	entry, present := dm.services[fqn]
	if present {
		if time.Now().Before(entry.expiryDate) {
			return entry.service, nil
		}

		dm.services[fqn] = nil
	}

	machineName := fqn[0 : len(fqn)-len(".local.")]

	ip, err := findIP(machineName)
	if err != nil {
		return nil, err
	}

	service, err := mdns.NewMDNSService("DockerMachine", "_ssh._tcp", "local.", fqn, 22, []net.IP{ip}, []string{"DockerMachine " + machineName})
	if err != nil {
		return nil, err
	}

	dm.services[fqn] = &serviceEntry{
		expiryDate: time.Now().Add(time.Duration(dm.Ttl) * time.Second),
		service:    service,
	}

	return service, nil
}

// findIP finds the IP address of a Docker Machine host given its name.
func findIP(name string) (net.IP, error) {
	log.Println("Looking for virtualbox machine", name)

	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	defer api.Close()

	machine, err := api.Load(name)
	if err != nil {
		return nil, err
	}

	driver := machine.Driver.DriverName()

	ip, err := machine.Driver.GetIP()
	if err != nil {
		return nil, err
	}

	log.Printf("Found %s(%s) with IP %s\n", name, driver, ip)

	return net.ParseIP(ip), nil
}
