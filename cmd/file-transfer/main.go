package main

import (
	"os"

	ft "file-transfer/internal/file-transfer"
)

func main() {
	command := ft.NewCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
