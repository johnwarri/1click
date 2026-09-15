package main

import (
	"fmt"
	"time"
)

func Run(p Platform) error {
	fmt.Println("=== System Test ===")
	fmt.Println()

	// What OS and chip is this running on?
	fmt.Printf("Operating system : %s\n", p.OS)
	fmt.Printf("Architecture     : %s\n", p.Arch)
	fmt.Printf("Time             : %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Who is running this?
	user, err := p.Capture("whoami")
	if err != nil {
		return fmt.Errorf("whoami failed: %w", err)
	}
	fmt.Printf("Current user     : %s\n", user)
	fmt.Println()

	// Where is this program sitting?
	whereCmd := "pwd"
	if p.IsWindows() {
		whereCmd = "cd"
	}
	location, err := p.Capture(whereCmd)
	if err != nil {
		return fmt.Errorf("location check failed: %w", err)
	}
	fmt.Printf("Running from     : %s\n", location)
	fmt.Println()

	fmt.Println("=== All good! ===")
	return nil
}
