package main

import (
	"fmt"
	"os"

	cmd "github.com/Eagerod/dupester/cmd/dupester"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
