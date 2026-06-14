package commands

import (
	"fmt"
	"os"
	"volt/internal/utils"
)

func Check(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		utils.ErrorCannotReadFile(path)
		os.Exit(1)
		return
	}

	script, err := utils.ParseScriptFile(path, data)
	if err != nil {
		fmt.Printf("✘ Error parsing file: %v\n", err)
		os.Exit(1)
		return
	}

	utils.CheckScriptConsole(&script)
}
