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

// GetSearchResults is a synchronous operation that will populate and return an ldap.SearchResult object
func (w *LDAPSession) GetSearchResults(request *ldap.SearchRequest) (result *ldap.SearchResult, err error) {
	return w.LConn.SearchWithPaging(request, 1000)
}

// ExecuteSearchRequest performs a paged search and writes results to the LDAPsession's defined results channel.
// it only returns an err
func (w *LDAPSession) ExecuteSearchRequest(searchRequest *ldap.SearchRequest) (err error) {
	if w.resultsChan == nil {
		return fmt.Errorf("no channel defined. Call SetChannel first, or use GetSearchResults instead")
	}

	defer close(w.resultsChan)

	// basically a re-implementation of the standard function: https://github.com/go-ldap/ldap/blob/master/v3/search.go#L253
	// but writes entries to a channel as it gets them instead of waiting for all pages to complete


	var pagingControl *ldap.ControlPaging
	control := ldap.FindControl(searchRequest.Controls, ldap.ControlTypePaging)
	if control == nil {
		pagingControl = ldap.NewControlPaging(w.PageSize)
		searchRequest.Controls = append(searchRequest.Controls, pagingControl)
	} else {
		castControl, ok := control.(*ldap.ControlPaging)
		if !ok {
			return fmt.Errorf("expected paging control to be of type *ControlPaging, got %v", control)
		}
		if castControl.PagingSize != w.PageSize {
			return fmt.Errorf("paging size given in search request (%d) conflicts with size given in search call (%d)", castControl.PagingSize, w.PageSize)
		}
		pagingControl = castControl
	}

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
			w.resultsChan <- entry
		}

		// I don't use these, but keeping them here just in case
		// TODO: add support for Referrals and Controls channels
		//for _, referral := range result.Referrals {
		//	//todo
		//}
		//for _, control := range result.Controls {
		//	//todo
		//}

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
