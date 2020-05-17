package dns

import (
	"fmt"
	"net"
	"strings"
)

// FindLDAPServers attempts to find LDAP servers in a domain via DNS. First it attempts looking up LDAP via SRV records,
// if that fails, it will just resolve the domain to an IP and return that.
func FindLDAPServers(domain string) (servers []string, err error) {
	_, srvs, err := net.LookupSRV("ldap", "tcp", domain)
	if err != nil {
		if strings.Contains(err.Error(), "No records found") {
			return net.LookupHost(domain)
		}
	}

	for _, s := range srvs {
		servers = append(servers, s.Target)
	}
	// also resolve the domain itself and return that IP
	domain_ips, _ := net.LookupHost(domain)
	for _, s := range domain_ips {
		servers = append(servers, s)
	}
	if len(servers) == 0 {
		err = fmt.Errorf("no LDAP servers found for domain: %s", domain)
		return
	}
	return servers, nil
}
