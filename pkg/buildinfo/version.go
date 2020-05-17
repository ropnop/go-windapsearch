package buildinfo

import "fmt"

var (
	Version     string
	GitSHA      string
	BuildDate   string
	BuildNumber string
)

func FormatVersionString() string {
	return fmt.Sprintf("Version: %s | GitSHA: %s | BuildDate: %s\n", Version, GitSHA, BuildDate)
}