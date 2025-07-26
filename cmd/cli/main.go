package main

import (
	"cornyk/gin-template/internal/commands"
	"os"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
