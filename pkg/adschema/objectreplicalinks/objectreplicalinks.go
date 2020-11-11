package objectreplicalinks

import "github.com/ropnop/go-windapsearch/pkg/adschema/objectreplicalinks/dnsrecord"

type ConvertObjectReplicaLink func([]byte) interface{}

var ObjectReplicaLinkFuncs = map[string]ConvertObjectReplicaLink{
	"dnsRecord": dnsrecord.ConvertDnsRecord,
}
