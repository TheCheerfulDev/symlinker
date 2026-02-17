package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"symlinker/entity"

	"gopkg.in/yaml.v3"
)

func expandPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home dir: %w", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get cwd: %w", err)
		}
		path = filepath.Join(cwd, path)
	}

	return filepath.Clean(path), nil
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

// readSymlinkResolved reads a symlink and also returns its resolved absolute path
// (resolving relative link targets relative to the symlink's directory).
func readSymlinkResolved(targetPath string) (stored string, resolvedAbs string, err error) {
	stored, err = os.Readlink(targetPath)
	if err != nil {
		return "", "", err
	}

	resolved := stored
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(filepath.Dir(targetPath), resolved)
	}
	return stored, filepath.Clean(resolved), nil
}

func symlinkPointsTo(targetPath, expectedSourceAbs string) (ok bool, stored string, resolved string, err error) {
	stored, resolved, err = readSymlinkResolved(targetPath)
	if err != nil {
		return false, "", "", err
	}

	expected := filepath.Clean(expectedSourceAbs)
	return resolved == expected, stored, resolved, nil
}

func getSymlinks() (entity.Symlinks, error) {
	file, err := os.ReadFile(symlinkerFile)
	if err != nil {
		return entity.Symlinks{}, fmt.Errorf("read %s: %w", symlinkerFile, err)
	}

	var symlinks entity.Symlinks
	if err := yaml.Unmarshal(file, &symlinks); err != nil {
		return entity.Symlinks{}, fmt.Errorf("parse %s: %w", symlinkerFile, err)
	}

	if err := symlinks.Validate(); err != nil {
		return entity.Symlinks{}, fmt.Errorf("invalid %s: %w", symlinkerFile, err)
	}

	return symlinks, nil
}
