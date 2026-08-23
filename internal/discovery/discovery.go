package discovery

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/twoBoots/battery/internal/config"
)

var ProjectMarkers = []string{
	".git",
	".cooper",
	"conductor",
	"package.json",
	"go.mod",
	"Cargo.toml",
	"pyproject.toml",
	"requirements.txt",
	"pom.xml",
	"build.gradle",
	".batteryrc",
}

var MonorepoPackageDirs = []string{
	"packages",
	"apps",
	"services",
	"libs",
	"modules",
	"crates",
}

// HasProjectMarker checks if a directory contains any project marker files or directories.
func HasProjectMarker(dirPath string) bool {
	for _, marker := range ProjectMarkers {
		if _, err := os.Stat(filepath.Join(dirPath, marker)); err == nil {
			return true
		}
	}
	return false
}

// DiscoverSiblingBarrels finds sibling directories containing git, cooper, or build markers.
func DiscoverSiblingBarrels(cwd string) []config.BarrelConfig {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		absCwd = cwd
	}

	parentDir := filepath.Dir(absCwd)
	currentDirName := filepath.Base(absCwd)
	barrels := make([]config.BarrelConfig, 0)

	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return barrels
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != currentDirName && !strings.HasPrefix(entry.Name(), ".") {
			fullPath := filepath.Join(parentDir, entry.Name())
			if HasProjectMarker(fullPath) {
				barrels = append(barrels, config.BarrelConfig{
					Name: entry.Name(),
					Path: "../" + entry.Name(),
				})
			}
		}
	}

	sort.Slice(barrels, func(i, j int) bool {
		return barrels[i].Name < barrels[j].Name
	})

	return barrels
}

// DiscoverMonorepoBarrels finds subdirectories within standard monorepo package containers.
func DiscoverMonorepoBarrels(cwd string) []config.BarrelConfig {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		absCwd = cwd
	}

	barrels := make([]config.BarrelConfig, 0)

	for _, pkgDir := range MonorepoPackageDirs {
		parentPkgPath := filepath.Join(absCwd, pkgDir)
		entries, err := os.ReadDir(parentPkgPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				fullPath := filepath.Join(parentPkgPath, entry.Name())
				if HasProjectMarker(fullPath) {
					barrels = append(barrels, config.BarrelConfig{
						Name: entry.Name(),
						Path: "./" + pkgDir + "/" + entry.Name(),
					})
				}
			}
		}
	}

	sort.Slice(barrels, func(i, j int) bool {
		return barrels[i].Name < barrels[j].Name
	})

	return barrels
}

// DetectProjectStructure infers the workspace topology based on surroundings.
func DetectProjectStructure(cwd string) config.ProjectStructure {
	monorepoBarrels := DiscoverMonorepoBarrels(cwd)
	if len(monorepoBarrels) > 0 {
		return config.StructureMonorepo
	}

	siblingBarrels := DiscoverSiblingBarrels(cwd)
	if len(siblingBarrels) > 0 {
		return config.StructureMultiRepo
	}

	return config.StructureCustom
}

// CandidateBarrelsResult contains the detected structure and candidate barrels.
type CandidateBarrelsResult struct {
	Structure config.ProjectStructure `json:"structure"`
	Barrels   []config.BarrelConfig   `json:"barrels"`
}

// DiscoverCandidateBarrels discovers candidate barrels matching the inferred structure.
func DiscoverCandidateBarrels(cwd string) CandidateBarrelsResult {
	structure := DetectProjectStructure(cwd)
	var barrels []config.BarrelConfig

	switch structure {
	case config.StructureMonorepo:
		barrels = DiscoverMonorepoBarrels(cwd)
	case config.StructureMultiRepo:
		barrels = DiscoverSiblingBarrels(cwd)
	default:
		barrels = []config.BarrelConfig{}
	}

	return CandidateBarrelsResult{
		Structure: structure,
		Barrels:   barrels,
	}
}
