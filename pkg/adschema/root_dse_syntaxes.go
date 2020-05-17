package adschema

import (
	"encoding/json"
)

// Manually creating mappings for Root DSE attributes, which are documented here:
// https://docs.microsoft.com/en-us/windows/win32/adschema/rootdse

// I'm only implementing ones I care about at this point
var RootDSEAttributeMap = map[string]bool {
	"defaultNamingContext": true,
	"dnsHostName": true,
	"domainFunctionality": true,
	"forestFunctionality": true,
	"domainControllerFunctionality": true,
	"rootDomainNamingContext": true,
	"currentTime": true,
}

func marshalRootDSEAttribute(e *ADAttribute) ([]byte, error) {
	switch e.Name {
	case "defaultNamingContext", "dnsHostName", "rootDomainNamingContext":
		return json.Marshal(string(e.ByteValues[0]))
	case "domainFunctionality", "forestFunctionality", "domainControllerFunctionality":
		level, ok := FunctionalityLevelsMapping[string(e.ByteValues[0])]
		if ok {
			return json.Marshal(level)
		} else {
			return json.Marshal(printable(e.ByteValues[0]))
		}
	case "currentTime":
		b, err := ConvertGeneralizedTime(e.Name, e.ByteValues[0])
		if err != nil {
			return nil, err
		}
		return json.Marshal(b)
	}
	return marshalUnknownAttribute(e)
}

