// Command app is the operator CLI for the template service.
//
// Subcommands live in their own files. This binary is what you use to
// migrate, seed, or open a one-off shell — not to serve requests. For that,
// use appd.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "migrate":
		fmt.Println("migrate: not implemented in skeleton; wire via hanzoai/base/tools/migrations")
	case "seed":
		fmt.Println("seed: not implemented in skeleton")
	case "version":
		fmt.Println("template dev")
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: app <migrate|seed|version>")
}
