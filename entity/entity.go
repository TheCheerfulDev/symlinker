package entity

import (
	"fmt"
	"strings"
)

type Symlink struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Symlinks struct {
	Links []Symlink `yaml:"symlinks"`
}

func (s Symlink) Validate() error {
	if strings.TrimSpace(s.Source) == "" {
		return fmt.Errorf("source is required")
	}
	if strings.TrimSpace(s.Target) == "" {
		return fmt.Errorf("target is required")
	}
	return nil
}

func (ss Symlinks) Validate() error {
	seenTargets := map[string]struct{}{}
	for _, link := range ss.Links {
		if err := link.Validate(); err != nil {
			return err
		}
		if _, ok := seenTargets[link.Target]; ok {
			return fmt.Errorf("duplicate target found: %s", link.Target)
		}
		seenTargets[link.Target] = struct{}{}
	}
	return nil
}
