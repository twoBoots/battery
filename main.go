package main

import (
	"os"

	"github.com/twoboots/battery/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
