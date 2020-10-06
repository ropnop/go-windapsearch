package modules

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
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
	return []string{"dn"}
}

func (d DnsNamesModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	locations := []string{"CN=MicrosoftDNS,DC=DomainDnsZones,%s", "CN=MicrosoftDNS,DC=ForestDnsZones,%s", "CN=MicrosoftDNS,CN=System,%s"}

	// I'm too stupid for filtering these entries via LDAP, tried for about a day and failed... now they're implemented here, feel free to improve this
	dnContainsFilter := []string{"DC=RootDNSServers", "in-addr.arpa,", "DC=_msdcs", "..TrustAnchors"}
	dnBeginsFilter := []string{"DC=DomainDnsZones,", "DC=ForestDnsZones,", "DC=_kerberos.", "DC=_ldap.", "DC=_kpasswd.", "DC=_gc.", "DC=@", "DC=_autodiscover."}
	baseDN := session.BaseDN
	results := make([]*ldap.SearchResult, 0)
	for _, location := range locations {
		session.BaseDN = fmt.Sprintf(location, baseDN)

		searchReq := session.MakeSimpleSearchRequest("(objectClass=*)", attrs)
		res, err := session.GetSearchResults(searchReq)

		if err != nil {
			return err
		}
		filteredResults := make([]*ldap.Entry, 0)
	outer:
		for _, entry := range res.Entries {
			for _, filterEntry := range dnContainsFilter { // Filter entries like rDNS zones
				if strings.Contains(entry.DN, filterEntry) {
					continue outer
				}
			}
			for _, filterEntry := range dnBeginsFilter { // filter entries used for discovery, e.g. _ldap
				if strings.HasPrefix(entry.DN, filterEntry) {
					continue outer
				}
			}
			for _, filterEntry := range locations { // filter the zone entries themselves
				if strings.HasPrefix(entry.DN, fmt.Sprintf(filterEntry, "")) {
					continue outer
				}
			}
			filteredResults = append(filteredResults, entry)
		}
		res.Entries = filteredResults

		results = append(results, res)
	}
	session.BaseDN = baseDN

	session.ManualWriteMultipleSearchResultsToChan(results)
	return nil
}
