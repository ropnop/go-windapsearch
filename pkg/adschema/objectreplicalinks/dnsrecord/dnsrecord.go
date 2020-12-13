package dnsrecord

import (
	"fmt"
	"net"
	"strings"
)

// See https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dnsp/6912b338-5472-4f59-b912-0edb536b6ed8
// for structure of dnsRecord
// thanks to @dirkjanm for https://github.com/dirkjanm/adidnsdump/blob/master/adidnsdump/dnsdump.py for inspiration

type DnsRecord struct {
	DataLength uint16 `struc:"uint16,sizeof=Data"`
	Type uint16 `struc:"uint16"`
	Version uint8 `struc:"uint8"`
	Rank uint8 `struc:"uint8"`
	Flags uint16 `struc:"uint16"`
	Serial uint32 `struc:"uint32"`
	TtlSeconds uint32 `struc:"uint32"`
	Reserved []byte `struc:"pad,[4]byte"`
	TimeStamp uint32 `struc:"uint32"`
	Data []byte `struc:"[]byte"`
}

const (
	DNS_TYPE_ZERO uint16  = 0x0
	DNS_TYPE_A uint16     = 0x1
	DNS_TYPE_NS uint16    = 0x2
	DNS_TYPE_MD uint16    = 0x3
	DNS_TYPE_MF uint16    = 0x4
	DNS_TYPE_CNAME uint16 = 0x5
	DNS_TYPE_SOA uint16   = 0x6
	DNS_TYPE_MB uint16    = 0x7
	DNS_TYPE_MG uint16    = 0x8
	DNS_TYPE_MR uint16    = 0x9
	DNS_TYPE_NULL uint16  = 0xA
	DNS_TYPE_WKS uint16   = 0xB
	DNS_TYPE_PTR uint16   = 0xC
	DNS_TYPE_HINFO uint16 = 0xD
	DNS_TYPE_MINFO uint16 = 0xE
	DNS_TYPE_MX uint16    = 0xF
	DNS_TYPE_TXT uint16   = 0x10
	DNS_TYPE_RP uint16    = 0x11
	DNS_TYPE_AFSDB uint16 = 0x12
	DNS_TYPE_X25 uint16   = 0x13
	DNS_TYPE_ISDN uint16  = 0x14
	DNS_TYPE_RT uint16    = 0x15
	DNS_TYPE_SIG uint16   = 0x18
	DNS_TYPE_KEY uint16   = 0x19
	DNS_TYPE_AAAA uint16  = 0x1C
	DNS_TYPE_LOC uint16   = 0x1D
	DNS_TYPE_NXT uint16   = 0x1E
	DNS_TYPE_SRV uint16   = 0x21
	DNS_TYPE_ATMA uint16  = 0x22
	DNS_TYPE_NAPTR uint16 = 0x23
	DNS_TYPE_DNAME uint16 = 0x27
	DNS_TYPE_DS uint16    = 0x2B
	DNS_TYPE_RRSIG uint16 = 0x2E
	DNS_TYPE_NSEC uint16  = 0x2F
	DNS_TYPE_DNSKEY uint16 = 0x30
	DNS_TYPE_DHCID uint16 = 0x31
	DNS_TYPE_NSEC3 uint16 = 0x32
	DNS_TYPE_NSEC3PARAM uint16 = 0x33
	DNS_TYPE_TLSA uint16 = 0x34
	DNS_TYPE_ALL uint16   = 0xFF
	DNS_TYPE_WINS uint16  = 0xFF01
	DNS_TYPE_WINSR uint16 = 0xFF02
)

var DnsTypes = map[uint16]string{
	DNS_TYPE_ZERO:       "DNS_TYPE_ZERO",
	DNS_TYPE_A:          "A",
	DNS_TYPE_NS:         "NS",
	DNS_TYPE_MD:         "MD",
	DNS_TYPE_MF:         "MF",
	DNS_TYPE_CNAME:      "CNAME",
	DNS_TYPE_SOA:        "SOA",
	DNS_TYPE_MB:         "MB",
	DNS_TYPE_MG:         "MG",
	DNS_TYPE_MR:         "MR",
	DNS_TYPE_NULL:       "NULL",
	DNS_TYPE_WKS:        "WKS",
	DNS_TYPE_PTR:        "PTR",
	DNS_TYPE_HINFO:      "HINFO",
	DNS_TYPE_MINFO:      "MINFO",
	DNS_TYPE_MX:         "MX",
	DNS_TYPE_TXT:        "TXT",
	DNS_TYPE_RP:         "RP",
	DNS_TYPE_AFSDB:      "AFSDB",
	DNS_TYPE_X25:        "X25",
	DNS_TYPE_ISDN:       "ISDN",
	DNS_TYPE_RT:         "RT",
	DNS_TYPE_SIG:        "SIG",
	DNS_TYPE_KEY:        "KEY",
	DNS_TYPE_AAAA:       "AAAA",
	DNS_TYPE_LOC:        "LOC",
	DNS_TYPE_NXT:        "NXT",
	DNS_TYPE_SRV:        "SRV",
	DNS_TYPE_ATMA:       "ATMA",
	DNS_TYPE_NAPTR:      "NAPTR",
	DNS_TYPE_DNAME:      "DNAME",
	DNS_TYPE_DS:         "DS",
	DNS_TYPE_RRSIG:      "RRSIG",
	DNS_TYPE_NSEC:       "NSEC",
	DNS_TYPE_DNSKEY:     "DNSKEY",
	DNS_TYPE_DHCID:      "DHCID",
	DNS_TYPE_NSEC3:      "NSEC3",
	DNS_TYPE_NSEC3PARAM: "NSEC3PARAM",
	DNS_TYPE_TLSA:       "TLSA",
	DNS_TYPE_ALL:        "ALL",
	DNS_TYPE_WINS:       "WINS",
	DNS_TYPE_WINSR:      "WINSR",
}

type DNS_RPC_RECORD_A struct {
	Ipv4Address net.IP `struc:"[4]byte,little"`
}

type DNS_RPC_RECORD_AAAA struct {
	Ipv6Address net.IP `struc:"[16]byte,little"`
}

type DNS_RPC_RECORD_NODE_NAME struct {
	NameNode DNS_RPC_NAME
}

type DNS_RPC_NAME struct {
	NameLength uint8 `struc:"uint8,sizeof=DnsName"`
	DnsName string `struc:"[]byte"`
}

type DNS_COUNT_NAME struct {
	Length uint8 `struc:"uint8,little,sizeof=RawName"`
	LabelCount uint8 `struc:"uint8,little"`
	RawName []byte
}

func (d DNS_COUNT_NAME) String() string {
	// <3 you dirkjanm https://github.com/dirkjanm/adidnsdump/blob/master/adidnsdump/dnsdump.py#L107
	var ind uint8 = 0
	var labels []string
	for i := uint8(0); i <= d.LabelCount; i++ {
		nextlen := uint8(d.RawName[ind : ind+1][0])
		labels = append(labels, string(d.RawName[ind+1:ind+1+nextlen]))
		ind += nextlen + 1
	}
	return strings.Join(labels, ".")
}

type DNS_RPC_RECORD_SOA struct {
	DwSerialNo int `struc:"uint32"`
	DwRefresh int `struc:"uint32"`
	DwRetry int `struc:"uint32"`
	DwExpire int `struc:"uint32"`
	DwMinimumTtl int `struc:"uint32"`
	NamePrimaryServer DNS_COUNT_NAME
	ZoneAdministratorEmail DNS_COUNT_NAME
}

func (d *DNS_RPC_RECORD_SOA) String() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", d.NamePrimaryServer, d.ZoneAdministratorEmail, d.DwSerialNo, d.DwRefresh, d.DwRetry, d.DwExpire, d.DwMinimumTtl)
}

// can't get bitmask to work?
//type DNS_RPC_RECORD_WKS struct {
//	IpAddress net.IP `struc:"[4]byte"`
//	ChProtocol int `struc:"uint8"`
//	BBitMask DNS_COUNT_NAME
//}
//
//func (d *DNS_RPC_RECORD_WKS) String() string {
//	return fmt.Sprintf("%s %d %d", d.IpAddress, d.ChProtocol, d.BBitMask)
//}

type DNS_RPC_RECORD_SRV struct {
	Priority int `struc:"uint16"`
	Weight int `struc:"uint16"`
	Port int `struc:"uint16"`
	NameTarget DNS_COUNT_NAME
}

func (d *DNS_RPC_RECORD_SRV) String() string {
	return fmt.Sprintf("%d %d %d %s", d.Priority, d.Weight, d.Port, d.NameTarget)
}

