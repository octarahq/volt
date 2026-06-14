package utils

import "fmt"

func ErrorMissingArg(command string, argument string) {
	fmt.Printf("✘ Command `%s` need missing argument : %s\n", command, argument)
}

func ErrorCannotReadFile(path string) {
	fmt.Printf("✘ Error while reading the file : %s\n", path)
}

func ErrorFileIsNotYaml() {
	fmt.Printf("✘ The file has to be a yaml file.\n")
}
