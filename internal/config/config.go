package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigFilename       = ".batteryrc"
	LocalConfigFilename  = ".batteryrc.local"
	CurrentConfigVersion = "1.0.0"
)

// ProjectStructure defines the topology of the workspace.
type ProjectStructure string

const (
	StructureMultiRepo ProjectStructure = "multi-repo"
	StructureMonorepo  ProjectStructure = "monorepo"
	StructureCustom    ProjectStructure = "custom"
)

// BarrelType differentiates standard repo barrels from nested battery orchestrators.
type BarrelType string

const (
	BarrelTypeBarrel  BarrelType = "barrel"
	BarrelTypeBattery BarrelType = "battery"
)

// BarrelConfig represents a single barrel entry in configuration.
type BarrelConfig struct {
	Name  string                     `json:"name"`
	Path  string                     `json:"path"`
	Type  BarrelType                 `json:"type,omitempty"`
	Role  string                     `json:"role,omitempty"`
	Tech  string                     `json:"tech,omitempty"`
	Docs  string                     `json:"docs,omitempty"`
	Jira  string                     `json:"jira,omitempty"`
	Extra map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON handles unmarshaling BarrelConfig and capturing dynamic extra fields.
func (b *BarrelConfig) UnmarshalJSON(data []byte) error {
	type Alias BarrelConfig
	var a struct {
		Alias
	}
	if err := json.Unmarshal(data, &a.Alias); err != nil {
		return err
	}
	*b = BarrelConfig(a.Alias)

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	delete(rawMap, "name")
	delete(rawMap, "path")
	delete(rawMap, "type")
	delete(rawMap, "role")
	delete(rawMap, "tech")
	delete(rawMap, "docs")
	delete(rawMap, "jira")

	if len(rawMap) > 0 {
		b.Extra = rawMap
	}
	return nil
}

// MarshalJSON handles marshaling BarrelConfig including any dynamic extra fields.
func (b BarrelConfig) MarshalJSON() ([]byte, error) {
	type Alias BarrelConfig
	a := Alias(b)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(a); err != nil {
		return nil, err
	}
	data := bytes.TrimSpace(buf.Bytes())

	if len(b.Extra) == 0 {
		return data, nil
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}
	if rawMap == nil {
		rawMap = make(map[string]json.RawMessage)
	}

	for k, v := range b.Extra {
		rawMap[k] = v
	}

	buf.Reset()
	enc = json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(rawMap); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

// BatteryConfig represents the canonical .batteryrc file format.
type BatteryConfig struct {
	Schema    string           `json:"$schema,omitempty"`
	Version   string           `json:"version"`
	Structure ProjectStructure `json:"structure"`
	Barrels   []BarrelConfig   `json:"barrels"`
}

// LocalBatteryConfig represents optional developer overrides in .batteryrc.local.
type LocalBatteryConfig struct {
	Structure ProjectStructure `json:"structure,omitempty"`
	Barrels   []BarrelConfig   `json:"barrels,omitempty"`
}

// EffectiveBarrel combines canonical and local configurations with metadata.
type EffectiveBarrel struct {
	Name            string                     `json:"name"`
	Path            string                     `json:"path"`
	Type            BarrelType                 `json:"type,omitempty"`
	Source          string                     `json:"source"` // "canonical" or "local"
	Role            string                     `json:"role,omitempty"`
	Tech            string                     `json:"tech,omitempty"`
	Docs            string                     `json:"docs,omitempty"`
	Jira            string                     `json:"jira,omitempty"`
	Exists          bool                       `json:"exists,omitempty"`
	AbsolutePath    string                     `json:"absolutePath,omitempty"`
	CooperTechStack string                     `json:"cooperTechStack,omitempty"`
	ProfilePath     string                     `json:"profilePath,omitempty"`
	HasProfile      bool                       `json:"hasProfile,omitempty"`
	Extra           map[string]json.RawMessage `json:"extra,omitempty"`
}

// EffectiveBatteryConfig represents the final merged configuration.
type EffectiveBatteryConfig struct {
	Version   string            `json:"version"`
	Structure ProjectStructure  `json:"structure"`
	Barrels   []EffectiveBarrel `json:"barrels"`
}

// DefaultConfig returns the initial default configuration.
func DefaultConfig() BatteryConfig {
	return BatteryConfig{
		Version:   CurrentConfigVersion,
		Structure: StructureMultiRepo,
		Barrels:   []BarrelConfig{},
	}
}

// LoadConfig loads the canonical .batteryrc file from the given directory.
func LoadConfig(cwd string) (*BatteryConfig, error) {
	filePath := filepath.Join(cwd, ConfigFilename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			return &cfg, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", ConfigFilename, err)
	}

	var cfg BatteryConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", ConfigFilename, err)
	}

	if cfg.Version == "" {
		cfg.Version = CurrentConfigVersion
	}
	if cfg.Structure == "" {
		cfg.Structure = StructureMultiRepo
	}
	if cfg.Barrels == nil {
		cfg.Barrels = []BarrelConfig{}
	}

	return &cfg, nil
}

// LoadLocalConfig loads local developer overrides from .batteryrc.local if present.
func LoadLocalConfig(cwd string) (*LocalBatteryConfig, error) {
	filePath := filepath.Join(cwd, LocalConfigFilename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", LocalConfigFilename, err)
	}

	var local LocalBatteryConfig
	if err := json.Unmarshal(data, &local); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", LocalConfigFilename, err)
	}

	return &local, nil
}

// SaveConfig saves canonical or local configuration to JSON format.
func SaveConfig(cfg interface{}, cwd string, isLocal bool) (string, error) {
	filename := ConfigFilename
	if isLocal {
		filename = LocalConfigFilename
	}
	targetFile := filepath.Join(cwd, filename)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return "", fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(targetFile, buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return targetFile, nil
}

// GetEffectiveConfig merges .batteryrc with .batteryrc.local to compute effective configuration.
func GetEffectiveConfig(cwd string) (*EffectiveBatteryConfig, error) {
	canonical, err := LoadConfig(cwd)
	if err != nil {
		return nil, err
	}

	local, err := LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}

	effectiveMap := make(map[string]EffectiveBarrel)
	order := make([]string, 0)

	// Add canonical barrels
	for _, b := range canonical.Barrels {
		effectiveMap[b.Name] = EffectiveBarrel{
			Name:   b.Name,
			Path:   b.Path,
			Type:   b.Type,
			Source: "canonical",
			Role:   b.Role,
			Tech:   b.Tech,
			Docs:   b.Docs,
			Jira:   b.Jira,
			Extra:  b.Extra,
		}
		order = append(order, b.Name)
	}

	// Override or append local barrels
	if local != nil && local.Barrels != nil {
		for _, b := range local.Barrels {
			if _, exists := effectiveMap[b.Name]; !exists {
				order = append(order, b.Name)
			}
			existing := effectiveMap[b.Name]
			role := b.Role
			if role == "" {
				role = existing.Role
			}
			tech := b.Tech
			if tech == "" {
				tech = existing.Tech
			}
			docs := b.Docs
			if docs == "" {
				docs = existing.Docs
			}
			jira := b.Jira
			if jira == "" {
				jira = existing.Jira
			}
			extra := make(map[string]json.RawMessage)
			for k, v := range existing.Extra {
				extra[k] = v
			}
			for k, v := range b.Extra {
				extra[k] = v
			}
			if len(extra) == 0 {
				extra = nil
			}

			effectiveMap[b.Name] = EffectiveBarrel{
				Name:   b.Name,
				Path:   b.Path,
				Type:   b.Type,
				Source: "local",
				Role:   role,
				Tech:   tech,
				Docs:   docs,
				Jira:   jira,
				Extra:  extra,
			}
		}
	}

	structure := canonical.Structure
	if local != nil && local.Structure != "" {
		structure = local.Structure
	}

	effectiveBarrels := make([]EffectiveBarrel, 0, len(order))
	for _, name := range order {
		effectiveBarrels = append(effectiveBarrels, effectiveMap[name])
	}

	return &EffectiveBatteryConfig{
		Version:   canonical.Version,
		Structure: structure,
		Barrels:   effectiveBarrels,
	}, nil
}

// AddBarrel adds a barrel to canonical or local configuration.
func AddBarrel(barrel BarrelConfig, cwd string, isLocal bool) (*EffectiveBatteryConfig, error) {
	normalizedPath := strings.TrimSpace(barrel.Path)
	barrelName := strings.TrimSpace(barrel.Name)
	if barrelName == "" {
		barrelName = InferBarrelName(normalizedPath)
	}
	barrelName = strings.TrimSpace(barrelName)

	if barrelName == "" || normalizedPath == "" {
		return nil, errors.New("barrel name and path cannot be empty")
	}

	if isLocal {
		local, err := LoadLocalConfig(cwd)
		if err != nil {
			return nil, err
		}
		if local == nil {
			local = &LocalBatteryConfig{Barrels: []BarrelConfig{}}
		}
		if local.Barrels == nil {
			local.Barrels = []BarrelConfig{}
		}

		for _, b := range local.Barrels {
			if b.Name == barrelName || b.Path == normalizedPath {
				return nil, fmt.Errorf("barrel '%s' or path '%s' already exists in %s", barrelName, normalizedPath, LocalConfigFilename)
			}
		}

		local.Barrels = append(local.Barrels, BarrelConfig{
			Name:  barrelName,
			Path:  normalizedPath,
			Type:  barrel.Type,
			Role:  barrel.Role,
			Tech:  barrel.Tech,
			Docs:  barrel.Docs,
			Jira:  barrel.Jira,
			Extra: barrel.Extra,
		})
		if _, err := SaveConfig(local, cwd, true); err != nil {
			return nil, err
		}
	} else {
		canonical, err := LoadConfig(cwd)
		if err != nil {
			return nil, err
		}

		for _, b := range canonical.Barrels {
			if b.Name == barrelName || b.Path == normalizedPath {
				return nil, fmt.Errorf("barrel '%s' or path '%s' already exists in %s", barrelName, normalizedPath, ConfigFilename)
			}
		}

		canonical.Barrels = append(canonical.Barrels, BarrelConfig{
			Name:  barrelName,
			Path:  normalizedPath,
			Type:  barrel.Type,
			Role:  barrel.Role,
			Tech:  barrel.Tech,
			Docs:  barrel.Docs,
			Jira:  barrel.Jira,
			Extra: barrel.Extra,
		})
		if _, err := SaveConfig(canonical, cwd, false); err != nil {
			return nil, err
		}
	}

	return GetEffectiveConfig(cwd)
}

// RemoveBarrel removes a barrel by name or path from canonical or local configuration.
func RemoveBarrel(identifier string, cwd string, isLocal bool) (*EffectiveBatteryConfig, error) {
	targetID := strings.TrimSpace(identifier)
	if targetID == "" {
		return nil, errors.New("barrel name or path cannot be empty")
	}

	if isLocal {
		local, err := LoadLocalConfig(cwd)
		if err != nil {
			return nil, err
		}
		if local == nil || len(local.Barrels) == 0 {
			return nil, fmt.Errorf("barrel '%s' not found in %s", targetID, LocalConfigFilename)
		}

		found := false
		newBarrels := make([]BarrelConfig, 0, len(local.Barrels))
		for _, b := range local.Barrels {
			if b.Name == targetID || b.Path == targetID {
				found = true
			} else {
				newBarrels = append(newBarrels, b)
			}
		}
		if !found {
			return nil, fmt.Errorf("barrel '%s' not found in %s", targetID, LocalConfigFilename)
		}

		local.Barrels = newBarrels
		if _, err := SaveConfig(local, cwd, true); err != nil {
			return nil, err
		}
	} else {
		canonical, err := LoadConfig(cwd)
		if err != nil {
			return nil, err
		}

		found := false
		newBarrels := make([]BarrelConfig, 0, len(canonical.Barrels))
		for _, b := range canonical.Barrels {
			if b.Name == targetID || b.Path == targetID {
				found = true
			} else {
				newBarrels = append(newBarrels, b)
			}
		}
		if !found {
			return nil, fmt.Errorf("barrel '%s' not found in %s", targetID, ConfigFilename)
		}

		canonical.Barrels = newBarrels
		if _, err := SaveConfig(canonical, cwd, false); err != nil {
			return nil, err
		}
	}

	return GetEffectiveConfig(cwd)
}

// InferBarrelName extracts the last directory name from a given path.
func InferBarrelName(filePath string) string {
	cleaned := strings.TrimRight(filePath, "/\\")
	if cleaned == "" {
		return "barrel"
	}
	base := filepath.Base(cleaned)
	if base == "." || base == "/" || base == "\\" || base == "" {
		return "barrel"
	}
	return base
}
