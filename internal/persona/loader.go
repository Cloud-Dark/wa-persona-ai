package persona

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses a persona YAML file.
func Load(path string) (*Persona, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read persona file %s: %w", path, err)
	}

	var p Persona
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse persona file %s: %w", path, err)
	}

	if err := validate(&p); err != nil {
		return nil, fmt.Errorf("invalid persona %s: %w", path, err)
	}

	return &p, nil
}

func validate(p *Persona) error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.SystemPrompt == "" {
		return fmt.Errorf("system_prompt is required")
	}
	if len(p.Traits) < 1 {
		return fmt.Errorf("at least one trait is required")
	}
	return nil
}
