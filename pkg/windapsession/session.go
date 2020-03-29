package windapsession

import (
	"crypto/tls"
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/dns"
	"gopkg.in/ldap.v3"
)



type LDAPSessionOptions struct {
	Domain string
	DomainController string
	Username string
	Password string
	Port int
	Secure bool
}

type Windapsession struct {
	lConn  *ldap.Conn
	BaseDN string
	attrs  []string
	DomainInfo DomainInfo
}

func NewWindapSession(options *LDAPSessionOptions) (sess Windapsession, err error) {
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
	sess = Windapsession{
		lConn: lConn,
	}
	err = sess.Bind(options.Username, options.Password)
	if err != nil {
		return
	}
	err = sess.getMetaData()
	if err != nil {
		return
	}
	return sess, nil
}

func (w *Windapsession) Bind(username, password string) (err error) {
	if username == "" {
		err = w.lConn.UnauthenticatedBind("")
	} else {
		err = w.lConn.Bind(username, password)
	}
	if err != nil {
		return
	}
	return
}

func (w *Windapsession) Close() {
	w.lConn.Close()
}

func (w *Windapsession) getMetaData() (err error) {
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"*"},
		nil)
	res, err := w.lConn.Search(sr)
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
	w.DomainInfo.DomainFunctionalityLevel = FunctionalityLevelsMapping[domainFunctionality]
	w.DomainInfo.ForestFunctionalityLevel = FunctionalityLevelsMapping[forestFunctionality]
	w.DomainInfo.DomainControllerFunctionalityLevel = FunctionalityLevelsMapping[domainControllerFunctionality]
	w.DomainInfo.ServerDNSName = res.Entries[0].GetAttributeValue("dnsHostName")
	w.BaseDN = defaultNamingContext
	return nil
}

