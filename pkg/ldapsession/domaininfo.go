package ldapsession

var FunctionalityLevelsMapping = map[string]string{
	"0": "2000",
	"1": "2003 Interim",
	"2": "2003",
	"3": "2008",
	"4": "2008 R2",
	"5": "2012",
	"6": "2012 R2",
	"7": "2016",
	"": "Unknown",
}

type DomainInfo struct {
	DomainFunctionalityLevel string
	ForestFunctionalityLevel string
	DomainControllerFunctionalityLevel string
	ServerDNSName string
}