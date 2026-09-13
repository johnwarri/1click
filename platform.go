package main

// platform.go is the plumbing. You normally don't need to edit this file.
// It figures out which OS you're on and gives you a small set of helpers
// for running OS-level commands the right way on each platform.

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Platform describes, in simple terms, the machine we're running on.
type Platform struct {
	OS   string // "windows", "macos", or "linux"
	Arch string // e.g. "amd64", "arm64"
}

// Detect returns the current platform.
func Detect() Platform {
	name := "linux" // default: linux and any other unix-like OS
	switch runtime.GOOS {
	case "windows":
		name = "windows"
	case "darwin":
		name = "macos"
	}
	return Platform{OS: name, Arch: runtime.GOARCH}
}

// IsWindows is a convenience for branching your logic per OS.
func (p Platform) IsWindows() bool { return p.OS == "windows" }

// shell wraps a raw command string in the correct native shell for this OS:
// cmd.exe on Windows, /bin/sh everywhere else.
func (p Platform) shell(command string) *exec.Cmd {
	if p.IsWindows() {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("sh", "-c", command)
}

// Run executes a shell command and streams its output straight to this
// terminal. It returns an error if the command exits with a non-zero status.
func (p Platform) Run(command string) error {
	cmd := p.shell(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Capture executes a shell command and returns its combined stdout+stderr as
// a string (trimmed), instead of printing it. Use this when you need the
// output as data.
func (p Platform) Capture(command string) (string, error) {
	out, err := p.shell(command).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Exec runs a program directly, with no shell involved. Use this when you
// don't need shell features like pipes, redirects, or globbing — it's safer
// because there's no shell to misinterpret your arguments.
func (p Platform) Exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
