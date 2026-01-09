package main

import (
	"os"

	"cli-tools/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args))
}
