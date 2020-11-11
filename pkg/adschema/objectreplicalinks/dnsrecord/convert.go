package dnsrecord

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/lunixbochs/struc"
	"log"
)

type DnsRecordInfo struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func ConvertDnsRecord(b []byte) interface{} {
	var record DnsRecord
	buf := bytes.NewReader(b)

	err := struc.UnpackWithOptions(buf, &record, &struc.Options{Order: binary.LittleEndian})
	if err != nil {
		log.Fatal(err)
	}
	data, err := DNSRecordDataToString(record.Type, record.Data)
	if err != nil {
		data = base64.StdEncoding.EncodeToString(record.Data)
	}
	return &DnsRecordInfo{
		Type: DnsTypes[record.Type],
		Data: data,
	}
}

func DNSRecordDataToString(dnsType uint16, data []byte) (string, error) {
	buf := bytes.NewReader(data)
	switch dnsType {
	case DNS_TYPE_A:
		var rpcRecord DNS_RPC_RECORD_A
		if err := struc.Unpack(buf, &rpcRecord); err != nil {
			return "", err
		}
		return rpcRecord.Ipv4Address.String(), nil
	case DNS_TYPE_AAAA:
		var rpcRecord DNS_RPC_RECORD_AAAA
		if err := struc.Unpack(buf, &rpcRecord); err != nil {
			log.Fatal(err)
			return "", err
		}
		return rpcRecord.Ipv6Address.String(), nil
	case DNS_TYPE_PTR, DNS_TYPE_NS, DNS_TYPE_CNAME, DNS_TYPE_DNAME, DNS_TYPE_MB, DNS_TYPE_MR, DNS_TYPE_MG, DNS_TYPE_MD, DNS_TYPE_MF:
		var nameRecord DNS_COUNT_NAME
		if err := struc.Unpack(buf, &nameRecord); err != nil {
			return "", err
		}
		return nameRecord.String(), nil

	case DNS_TYPE_SRV:
		var srvRecord DNS_RPC_RECORD_SRV
		if err := struc.Unpack(buf, &srvRecord); err != nil {
			return "", err
		}
		return srvRecord.String(), nil
	case DNS_TYPE_HINFO, DNS_TYPE_ISDN, DNS_TYPE_TXT, DNS_TYPE_X25, DNS_TYPE_LOC:
		var nameRecord DNS_RPC_NAME
		if err := struc.Unpack(buf, &nameRecord); err != nil {
			return "", err
		}
		return nameRecord.DnsName, nil
	case DNS_TYPE_SOA:
		var soaRecord DNS_RPC_RECORD_SOA
		if err := struc.Unpack(buf, &soaRecord); err != nil {
			return "", err
		}
		return soaRecord.String(), nil

	default:
		return "", fmt.Errorf("unimplemented type: %s", DnsTypes[dnsType])
	}
}
