package cli

import (
	"fmt"
	"os"
)

// Execute runs the command-line interface.
func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("usage: auditexport <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if err := Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
