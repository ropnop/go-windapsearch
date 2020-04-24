package ldapsession

import (
	"fmt"
	"gopkg.in/ldap.v3"
)

func (w *LDAPSession) MakeSimpleSearchRequest(filter string, attrs []string) *ldap.SearchRequest {
	return ldap.NewSearchRequest(
		w.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		attrs,
		nil)
}

func (w *LDAPSession) GetSearchResults(request *ldap.SearchRequest) (result *ldap.SearchResult, err error) {
	return w.lConn.SearchWithPaging(request, 1000)
}

func (w *LDAPSession) AddExtraFilter(filter, extra string) string {
	return fmt.Sprintf("(&(%s)(%s))", filter, extra)
}


