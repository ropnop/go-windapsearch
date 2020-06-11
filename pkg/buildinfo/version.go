package buildinfo

import (
	"fmt"
	"runtime"
)

var (
	Version     = "dev"
	GitSHA      = "N/A"
	BuildDate   = "N/A"
	BuildNumber = ""
	GoVersion   string
)

func FormatVersionString() string {
	GoVersion = runtime.Version()
	return fmt.Sprintf("Version: %s (%s) | Built: %s (%s) | Ronnie Flathers @ropnop\n", Version, GitSHA, BuildDate, GoVersion)
}
