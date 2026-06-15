package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"volt/internal/types"

	"github.com/invopop/jsonschema"
)

func main() {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
	}

	schema := reflector.Reflect(&types.VoltScript{})

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Println("Error generating schema:", err)
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}
	scriptsDir := wd

	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		fmt.Println("Error creating scripts directory:", err)
		os.Exit(1)
	}

	schemaPath := filepath.Join(scriptsDir, "volt-schema.json")
	err = os.WriteFile(schemaPath, schemaJSON, 0644)
	if err != nil {
		fmt.Println("Error writing schema file:", err)
		os.Exit(1)
	}

	fmt.Printf("Schema successfully generated at %s\n", schemaPath)
}
