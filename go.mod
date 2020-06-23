module github.com/ropnop/go-windapsearch

replace github.com/ropnop/ldap/v3 => /Users/RonnieFlathers/go/src/github.com/ropnop/ldap/v3

//replace github.com/Azure/go-ntlmssp => /Users/RonnieFlathers/go/src/github.com/ropnop/go-ntlmssp

//replace github.com/ropnop/go-ntlm => /Users/RonnieFlathers/go/src/github.com/ropnop/go-ntlm

//replace github.com/go-asn1-ber/asn1-ber => /Users/RonnieFlathers/go/src/github.com/ropnop/asn1-ber

go 1.13

require (
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c
	github.com/audibleblink/msldapuac v0.2.0
	github.com/bwmarrin/go-objectsid v0.0.0-20191126144531-5fee401a2f37
	github.com/go-asn1-ber/asn1-ber v1.5.0
	github.com/magefile/mage v1.9.0
	github.com/ropnop/ldap/v3 v3.1.11-0.20200611014906-485c70f019f1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/tcnksm/go-input v0.0.0-20180404061846-548a7d7a8ee8
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
)
