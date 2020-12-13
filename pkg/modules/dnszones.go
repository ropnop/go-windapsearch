package modules

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type DnsZonesModule struct{}

func init() {
	AllModules = append(AllModules, new(DnsZonesModule))
}

func (d DnsZonesModule) Name() string {
	return "dns-zones"
}

func (d DnsZonesModule) Description() string {
	return "List all DNS Zones"
}

func (d DnsZonesModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(d.Name(), pflag.ExitOnError)
	return flags
}

func (d DnsZonesModule) DefaultAttrs() []string {
	return []string{"name"}
}

func (d DnsZonesModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	locations := []string{"CN=MicrosoftDNS,DC=DomainDnsZones,%s", "CN=MicrosoftDNS,DC=ForestDnsZones,%s", "CN=MicrosoftDNS,CN=System,%s"}
	baseDN := session.BaseDN
	var requests []*ldap.SearchRequest
	for _, location := range locations {
		dn := fmt.Sprintf(location, baseDN)

		searchReq := session.MakeSearchRequestWithDN(dn, "(&(objectClass=dnsZone)(!name=RootDNSServers)(!name=*.in-addr.arpa)(!name=_msdcs.*)(!name=..TrustAnchors))", attrs)
		requests = append(requests, searchReq)
	}
	return session.ExecuteBulkSearchRequest(requests)
}
