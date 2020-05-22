// +build mage

// magefile inspired/copied from Hugo's: https://github.com/gohugoio/hugo/blob/master/magefile.go

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	packageName = "github.com/ropnop/go-windapsearch"
)

var (
	curDir    string
	binDir    string
	distDir   string
	hash      string
	buildDate string
	buildNum  string
	version   string
)

var ldflags = `-w -s` +
	` -X $PKG/pkg/buildinfo.GitSHA=$GIT_SHA` +
	` -X $PKG/pkg/buildinfo.BuildDate=$DATE` +
	` -X $PKG/pkg/buildinfo.Version=$VERSION` +
	` -X $PKG/pkg/buildinfo.BuildNumber=$BUILDNUM`

var targets = "linux/amd64 darwin/amd64 windows/amd64 linux/arm64"

var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")

	curDir, err := os.Getwd()
	if err != nil {
		curDir = "." //hack
	}
	binDir = path.Join(curDir, "bin")
	distDir = path.Join(curDir, "dist")
}

// Build Compile all cmd packages for current GOOS and GOARCh and put in ./bin
func Build() error {
	err := sh.RunWith(flagEnv(), goexe, "install", "-ldflags", ldflags, "$PKG/cmd/...")
	if err != nil {
		return err
	}
	fmt.Printf("[+] Compiled binary to %s\n", binDir)
	return nil
}

var gox = sh.RunCmd("gox")

// Dist Cross-compile for Windows, Linux, Mac x64 and put in ./dist
func Dist() error {
	ldflags += " -extldflags \"-static\""
	mg.Deps(checkGox)
	fmt.Printf("[+] Cross compiling for: %q\n", targets)
	err := sh.RunWith(
		flagEnv(),
		"gox",
		"-parallel=3",
		"-output",
		"$DISTDIR/{{.Dir}}-{{.OS}}-{{.Arch}}",
		"--osarch=$TARGETS",
		"-ldflags",
		ldflags,
		"$PKG/cmd/...")
	if err != nil {
		return err
	}
	fmt.Printf("[+] Cross compiled binaries in: %s\n", distDir)
	return nil

}

func checkGox() error {
	_, err := exec.LookPath("gox")
	if err != nil {
		return sh.Run(goexe, "get", "-u", "github.com/mitchellh/gox")
	}
	return nil
}

// Clean Delete bin and dist dirs
func Clean() {
	fmt.Println("[+] Removing bin and dist...")
	os.RemoveAll(binDir)
	os.RemoveAll(distDir)
}

// set up environment variables
func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	if version = os.Getenv("VERSION"); version == "" {
		version = "dev"
	}
	if buildNum = os.Getenv("BUILDNUM"); buildNum == "" {
		buildNum = "local"
	}

	return map[string]string{
		"PKG":         packageName,
		"GOBIN":       binDir,
		"GIT_SHA":     hash,
		"DATE":        time.Now().Format("01/02/06"),
		"VERSION":     version,
		"BUILDNUM":    buildNum,
		"DISTDIR":     distDir,
		"CGO_ENABLED": "1", //bug: when this is disabled, DNS gets wonky
		"TARGETS":     targets,
	}
}
