package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/go-objectsid"
)


func WindowsGuidFromBytes(b []byte) (string, error) {
	if len(b) != 16 {
		return "", fmt.Errorf("GUID must be 16 bytes")
	}
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		binary.LittleEndian.Uint32(b[:4]),
		binary.LittleEndian.Uint16(b[4:6]),
		binary.LittleEndian.Uint16(b[6:8]),
		b[8:10],
		b[10:]), nil
}

func WindowsSIDFromBytes(b []byte) (string, error) {
	if len(b) < 12 {
		return "", fmt.Errorf("windows SID seems too short")
	}
	sid := objectsid.Decode(b)
	return sid.String(), nil
}