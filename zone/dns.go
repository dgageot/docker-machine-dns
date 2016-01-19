package zone

import "github.com/miekg/dns"

// SetTTL sets the TTL for each record.
func SetTTL(records []dns.RR, ttl uint32) []dns.RR {
	for i, _ := range records {
		records[i].Header().Ttl = ttl
	}

	return records
}
