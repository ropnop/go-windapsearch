package ldapsession

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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
	w.Log.Infof("making search request with filter: %q", request.Filter)
	return w.LConn.SearchWithPaging(request, 1000)
}

func (w *LDAPSession) ManualWriteSearchResultsToChan(results *ldap.SearchResult) {
	for _, entry := range results.Entries {
		w.resultsChan <- entry
	}
	close(w.resultsChan)
}

// ExecuteSearchRequest performs a paged search and writes results to the LDAPsession's defined results channel.
// it only returns an err
func (w *LDAPSession) ExecuteSearchRequest(searchRequest *ldap.SearchRequest) error {
	w.Log.WithFields(logrus.Fields{"filter": searchRequest.Filter, "attributes": searchRequest.Attributes}).Infof("sending LDAP search request")

	if w.Channels == nil {
		return fmt.Errorf("no channels defined. Call SetChannels first, or use GetSearchResults instead")
	}

	defer func() {
		w.Log.Debugf("search finished. closing channels...")
		w.CloseChannels()
	}()

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
	pageNumber := 0

PagedSearch:
	for {
		select {
		case <-w.ctx.Done():
			w.Log.Warn("cancel received. aborting remaining pages")
			return nil
		default:
			w.Log.Debugf("making paged request...\n")
			result, err := w.LConn.Search(searchRequest)
			w.Log.Debugf("Looking for Paging Control...\n")
			pageNumber++
			if err != nil {
				return err
			}
			if result == nil {
				return ldap.NewError(ldap.ErrorNetwork, errors.New("ldap: packet not received"))
			}

			for _, entry := range result.Entries {
				w.Channels.Entries <- entry
			}

			w.Log.Infof("Received page %d with %d LDAP entries...", pageNumber, len(result.Entries))

			for _, referral := range result.Referrals {
				w.Channels.Referrals <- referral
			}

			for _, control := range result.Controls {
				w.Channels.Controls <- control
			}

			w.Log.Debugf("Looking for Paging Control...")
			pagingResult := ldap.FindControl(result.Controls, ldap.ControlTypePaging)
			if pagingResult == nil {
				pagingControl = nil
				w.Log.Debugf("Could not find paging control.  Breaking...")
				break PagedSearch
			}

			cookie := pagingResult.(*ldap.ControlPaging).Cookie
			if len(cookie) == 0 {
				pagingControl = nil
				w.Log.Debugf("Could not find cookie.  Breaking...")
				break PagedSearch
			}
			pagingControl.SetCookie(cookie)
		}
	}

	if pagingControl != nil {
		w.Log.Debugf("Abandoning Paging...")
		pagingControl.PagingSize = 0
		w.LConn.Search(searchRequest)
	}
	return nil
}
