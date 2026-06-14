package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"volt/internal/types"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

func ParseScriptFile(path string, data []byte) (types.VoltScript, error) {
	var script types.VoltScript
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		err := yaml.Unmarshal(data, &script)
		return script, err
	case ".json":
		err := json.Unmarshal(data, &script)
		return script, err
	case ".toml":
		err := toml.Unmarshal(data, &script)
		return script, err
	default:
		return script, fmt.Errorf("unsupported format %s", ext)
	}
}
