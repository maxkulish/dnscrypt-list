//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/carolynvs/magex/pkg"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binDir     = "./bin"
	project    = "dnscrypt-list"
	appVersion = "0.1"
	cmdDir     = "./cmd"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// Build compiles and lints the api-monitor daemon version
func Build() error {
	mg.Deps(Lint)
	mg.Deps(Test)

	macBuild := make(map[string]string)
	macBuild["GOOS"] = "darwin"
	macBuild["GOARCH"] = "amd64"

	ldflags := fmt.Sprintf(
		`-w -s -X main.BuildVersion=%s`,
		appVersion,
	)

	cliPath := fmt.Sprintf("%s/macos/%s", binDir, project)
	cmdPath := fmt.Sprintf("%s/%s/main.go", cmdDir, project)

	return sh.RunWithV(
		macBuild,
		"go", "build", "-ldflags", ldflags,
		"-o", cliPath, cmdPath,
	)
}

// Lint check the code
func Lint() error {
	mg.Deps(EnsureRevive)

	return sh.RunV("revive", "-formatter", "friendly", "./...")
}

// EnsureGoLint install golint if needed
func EnsureRevive() error {
	return pkg.EnsurePackage("github.com/mgechev/revive", "")
}

// Test runs all test
func Test() error {

	return sh.RunV("go", "test", "-v", "./...")
}
