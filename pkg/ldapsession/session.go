package ldapsession

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"github.com/ropnop/go-windapsearch/pkg/dns"
	"gopkg.in/ldap.v3"
)

type LDAPSessionOptions struct {
	Domain           string
	DomainController string
	Username         string
	Password         string
	Port             int
	Secure           bool
	PageSize         int
}

type LDAPSession struct {
	LConn       *ldap.Conn
	PageSize    uint32
	BaseDN      string
	attrs       []string
	DomainInfo  DomainInfo
	resultsChan chan *ldap.Entry
	ctx         context.Context
	Channels    *ResultChannels
}

type ResultChannels struct {
	Entries   chan *ldap.Entry
	Referrals chan string
	Controls  chan ldap.Control
}

type DomainInfo struct {
	Metadata                           *ldap.SearchResult
	DomainFunctionalityLevel           string
	ForestFunctionalityLevel           string
	DomainControllerFunctionalityLevel string
	ServerDNSName                      string
}

func NewLDAPSession(options *LDAPSessionOptions, ctx context.Context) (sess *LDAPSession, err error) {
	port := options.Port
	dc := options.DomainController
	if port == 0 {
		if options.Secure {
			port = 636
		} else {
			port = 389
		}
	}
	if dc == "" {
		dcs, err := dns.FindLDAPServers(options.Domain)
		if err != nil {
			return sess, err
		}
		dc = dcs[0]
	}
	var url string

	if options.Secure {
		url = fmt.Sprintf("ldaps://%s:%d", dc, port)
	} else {
		url = fmt.Sprintf("ldap://%s:%d", dc, port)
	}

	lConn, err := ldap.DialURL(url)
	if err != nil {
		return
	}
	if options.Secure {
		lConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	}
	sess = &LDAPSession{
		LConn: lConn,
	}

	sess.PageSize = uint32(options.PageSize)

	err = sess.Bind(options.Username, options.Password)
	if err != nil {
		return
	}
	_, err = sess.GetDefaultNamingContext()
	if err != nil {
		return
	}

	sess.NewChannels(ctx)
	return sess, nil
}

func (w *LDAPSession) SetChannels(chs *ResultChannels, ctx context.Context) {
	w.Channels = chs
	w.ctx = ctx
}

func (w *LDAPSession) NewChannels(ctx context.Context) {
	w.Channels = &ResultChannels{
		Entries:   make(chan *ldap.Entry),
		Referrals: make(chan string),
		Controls:  make(chan ldap.Control),
	}
	w.ctx = ctx
	return
}

func (w *LDAPSession) CloseChannels() {
	if w.Channels.Entries != nil {
		close(w.Channels.Entries)
	}
	if w.Channels.Controls != nil {
		close(w.Channels.Controls)
	}
	if w.Channels.Referrals != nil {
		close(w.Channels.Referrals)
	}
	//if w.ctx != nil {
	//	w.ctx.Done()
	//}
}

//func (w *LDAPSession) SetResultsChannel(ch chan *ldap.Entry, ctx context.Context) {
//	w.resultsChan = ch
//	w.ctx = ctx
//}

func (w *LDAPSession) Bind(username, password string) (err error) {
	if username == "" {
		err = w.LConn.UnauthenticatedBind("")
	} else {
		err = w.LConn.Bind(username, password)
	}
	if err != nil {
		return
	}
	return
}

func (w *LDAPSession) Close() {
	w.LConn.Close()
}

func (w *LDAPSession) GetDefaultNamingContext() (string, error) {
	if w.BaseDN != "" {
		return w.BaseDN, nil
	}
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"defaultNamingContext"},
		nil)
	res, err := w.LConn.Search(sr)
	if err != nil {
		return "", err
	}
	if len(res.Entries) == 0 {
		return "", fmt.Errorf("error getting metadata: No LDAP responses from server")
	}
	defaultNamingContext := res.Entries[0].GetAttributeValue("defaultNamingContext")
	if defaultNamingContext == "" {
		return "", fmt.Errorf("error getting metadata: attribute defaultNamingContext missing")
	}
	w.BaseDN = defaultNamingContext
	return w.BaseDN, nil

}

func (w *LDAPSession) getMetaData() (err error) {
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"*"},
		nil)
	res, err := w.LConn.Search(sr)
	if err != nil {
		return
	}
	if len(res.Entries) == 0 {
		return fmt.Errorf("error getting metadata: No LDAP responses from server")
	}
	defaultNamingContext := res.Entries[0].GetAttributeValue("defaultNamingContext")
	if defaultNamingContext == "" {
		return fmt.Errorf("error getting metadata: attribute defaultNamingContext missing")
	}
	domainFunctionality := res.Entries[0].GetAttributeValue("domainFunctionality")
	forestFunctionality := res.Entries[0].GetAttributeValue("forestFunctionality")
	domainControllerFunctionality := res.Entries[0].GetAttributeValue("domainControllerFunctionality")
	w.DomainInfo.DomainFunctionalityLevel = adschema.FunctionalityLevelsMapping[domainFunctionality]
	w.DomainInfo.ForestFunctionalityLevel = adschema.FunctionalityLevelsMapping[forestFunctionality]
	w.DomainInfo.DomainControllerFunctionalityLevel = adschema.FunctionalityLevelsMapping[domainControllerFunctionality]
	w.DomainInfo.ServerDNSName = res.Entries[0].GetAttributeValue("dnsHostName")
	w.BaseDN = defaultNamingContext
	w.DomainInfo.Metadata = res
	return nil
}

func (w *LDAPSession) ReturnMetadataResults() error {
	for _, entry := range w.DomainInfo.Metadata.Entries {
		w.resultsChan <- entry
	}
	return nil
}
