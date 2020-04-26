package ldapsession

import (
	"errors"
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
	return w.LConn.SearchWithPaging(request, 1000)
}

func (w *LDAPSession) AddExtraFilter(filter, extra string) string {
	return fmt.Sprintf("(&(%s)(%s))", filter, extra)
}

func (w *LDAPSession) SearchWithPagingToChannel(searchRequest *ldap.SearchRequest, ch chan *ldap.Entry, pagingSize uint32) error {
	// basically a re-implementation of the standard function: https://github.com/go-ldap/ldap/blob/master/v3/search.go#L253
	// but writes entries to a channel as it gets them instead of waiting for all pages to complete

	defer close(ch)



	var pagingControl *ldap.ControlPaging
	control := ldap.FindControl(searchRequest.Controls, ldap.ControlTypePaging)
	if control == nil {
		pagingControl = ldap.NewControlPaging(pagingSize)
		searchRequest.Controls = append(searchRequest.Controls, pagingControl)
	} else {
		castControl, ok := control.(*ldap.ControlPaging)
		if !ok {
			return fmt.Errorf("expected paging control to be of type *ControlPaging, got %v", control)
		}
		if castControl.PagingSize != pagingSize {
			return fmt.Errorf("paging size given in search request (%d) conflicts with size given in search call (%d)", castControl.PagingSize, pagingSize)
		}
		pagingControl = castControl
	}

	searchResult := new(ldap.SearchResult)
	for {
		result, err := w.LConn.Search(searchRequest)
		w.LConn.Debug.Printf("Looking for Paging Control...")
		if err != nil {
			return err
		}
		if result == nil {
			return ldap.NewError(ldap.ErrorNetwork, errors.New("ldap: packet not received"))
		}

		for _, entry := range result.Entries {
			ch <- entry
			searchResult.Entries = append(searchResult.Entries, entry)
		}
		for _, referral := range result.Referrals {
			searchResult.Referrals = append(searchResult.Referrals, referral)
		}
		for _, control := range result.Controls {
			searchResult.Controls = append(searchResult.Controls, control)
		}

		w.LConn.Debug.Printf("Looking for Paging Control...")
		pagingResult := ldap.FindControl(result.Controls, ldap.ControlTypePaging)
		if pagingResult == nil {
			pagingControl = nil
			w.LConn.Debug.Printf("Could not find paging control.  Breaking...")
			break
		}

		cookie := pagingResult.(*ldap.ControlPaging).Cookie
		if len(cookie) == 0 {
			pagingControl = nil
			w.LConn.Debug.Printf("Could not find cookie.  Breaking...")
			break
		}
		pagingControl.SetCookie(cookie)
	}

	if pagingControl != nil {
		w.LConn.Debug.Printf("Abandoning Paging...")
		pagingControl.PagingSize = 0
		w.LConn.Search(searchRequest)
	}
	return nil
}
