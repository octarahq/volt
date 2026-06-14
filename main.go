package main

import (
	"os"
	"volt/internal/commands"
	"volt/internal/utils"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		os.Exit(0)
	}

	switch args[1] {
	case "run":
		if len(args) == 2 {
			utils.ErrorMissingArg("run", "path")
			return
		}
		commands.Run(args[2])
	case "check":
		if len(args) == 2 {
			utils.ErrorMissingArg("check", "path")
			return
		}
		commands.Check(args[2])
	case "create":
		commands.Create()
	}
}
