package main

import (
	"fmt"
	"os"
)

func main() {
	p := Detect()
	if err := Run(p); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1) // non-zero exit tells the caller something failed
	}
}
