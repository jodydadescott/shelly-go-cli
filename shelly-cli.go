package main

import (
	"os"

	"github.com/jodydadescott/shelly-go-cli/cmd"
)

func main() {
	err := cmd.NewCmd().Execute()

	if err != nil {
		os.Exit(1)
	}
}
