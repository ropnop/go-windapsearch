package modules

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
	"strings"
)

type DnsNamesModule struct{}

func init() {
	AllModules = append(AllModules, new(DnsNamesModule))
}

func (d DnsNamesModule) Name() string {
	return "dns-names"
}

func (d DnsNamesModule) Description() string {
	return "List all DNS Names"
}

func (d DnsNamesModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(d.Name(), pflag.ExitOnError)
	return flags
}

func (d DnsNamesModule) DefaultAttrs() []string {
	return []string{"name", "dnsTombstoned"}
}

// Optional function for the module interface that will be called by searchResultWorker
// This will hide by default a lot of extraneous entries we usually don't care about
// use '--ignore-display-filters' to bypass this and display everything
func (d DnsNamesModule) DisplayFilter(entry *adschema.ADEntry) bool {
	locations := []string{"CN=MicrosoftDNS,DC=DomainDnsZones,%s", "CN=MicrosoftDNS,DC=ForestDnsZones,%s", "CN=MicrosoftDNS,CN=System,%s"}
	dnContainsFilter := []string{"DC=RootDNSServers", "in-addr.arpa,", "DC=_msdcs", "..TrustAnchors"}
	dnBeginsFilter := []string{"DC=DomainDnsZones,", "DC=ForestDnsZones,", "DC=_kerberos.", "DC=_ldap.", "DC=_kpasswd.", "DC=_gc.", "DC=@", "DC=_autodiscover."}
	for _, filterEntry := range dnContainsFilter { // Filter entries like rDNS zones
		if strings.Contains(entry.DN, filterEntry) {
			return false
		}
	}
	for _, filterEntry := range dnBeginsFilter { // filter entries used for discovery, e.g. _ldap
		if strings.HasPrefix(entry.DN, filterEntry) {
			return false
		}
	}
	for _, filterEntry := range locations { // filter the zone entries themselves
		if strings.HasPrefix(entry.DN, fmt.Sprintf(filterEntry, "")) {
			return false
		}
	}
	return true

}

func (d DnsNamesModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	locations := []string{"CN=MicrosoftDNS,DC=DomainDnsZones,%s", "CN=MicrosoftDNS,DC=ForestDnsZones,%s", "CN=MicrosoftDNS,CN=System,%s"}
	var searchRequests []*ldap.SearchRequest
	for _, location := range locations {
		dn := fmt.Sprintf(location, session.BaseDN)

		searchReq := session.MakeSearchRequestWithDN(dn, "(objectClass=*)", attrs)
		searchRequests = append(searchRequests, searchReq)
	}

	return session.ExecuteBulkSearchRequest(searchRequests)
}
