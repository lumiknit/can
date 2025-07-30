//go:build mage

// A comment on the package will be output when you list the targets of a
// magefile.
package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/sh"
)

const binName = "can.exe"

var Default = Build

func buildEnv() map[string]string {
	return map[string]string{
		"CGO_ENABLED": "0",
	}
}

// Build compiles the server binary
func Build(ctx context.Context) {
	if err := sh.RunWithV(
		buildEnv(),
		"go", "build",
		"-o", binName,
		"./cmd/server",
	); err != nil {
		fmt.Println("Build failed:", err)
		return
	}
}

// Build for release version
func BuildRelease(ctx context.Context) {
	if err := sh.RunWithV(
		buildEnv(),
		"go", "build",
		"-ldflags", "-s -w",
		"-o", binName,
		"./cmd/server",
	); err != nil {
		fmt.Println("Build failed:", err)
		return
	}
}

// Run builds and runs the server
func Run(ctx context.Context) {
	var err error
	fmt.Println("[Mage] Starting the server...\n")

	command := append([]string{
		"go", "run", "./cmd/server/",
		}, []string{}...)

	err = sh.RunV(command[0], command[1:]...)

	if err != nil {
		fmt.Println("Build failed:", err)
		return
	}
}

func Pre(ctx context.Context) {
	Prebuild(ctx)
}

// Prebuild runs the build command before running the server
func Prebuild(ctx context.Context) {
	fmt.Println("[Mage] Formatting...\n")
	sh.RunV("go", "fmt", "./...")

	fmt.Println("[Mage] Running vet...\n")
	sh.RunV("go", "vet", "./...")
}

