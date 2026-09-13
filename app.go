package main

// app.go is YOUR file. Put your logic in Run().
//
// The Platform argument (p) tells you which OS you're on so you can issue the
// right OS-level command for each one. Delete everything under the EXAMPLE
// line and replace it with your own code.

import "fmt"

func Run(p Platform) error {
	fmt.Printf("Running on %s (%s)\n\n", p.OS, p.Arch)

	// ---------------------------------------------------------------------
	// EXAMPLE — safe to delete. Shows the three ways to run OS-level work.
	// ---------------------------------------------------------------------

	// 1) Run a command that is DIFFERENT on each OS.
	//    Same idea, different command name: list the current directory.
	listCmd := "ls -la"
	if p.IsWindows() {
		listCmd = "dir"
	}
	fmt.Println("> listing the current directory:")
	if err := p.Run(listCmd); err != nil {
		return fmt.Errorf("directory listing failed: %w", err)
	}

	// 2) Capture a command's output as data instead of printing it.
	//    ("whoami" happens to exist on all three OSes.)
	user, err := p.Capture("whoami")
	if err != nil {
		return fmt.Errorf("whoami failed: %w", err)
	}
	fmt.Printf("\n> current user: %s\n", user)

	// 3) Run a program directly, no shell. Uncomment to try:
	// if err := p.Exec("go", "version"); err != nil {
	// 	return err
	// }

	return nil
}
