package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"symlinker/entity"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the symlinks defined in the configuration file",
	Long:  `Verifies the symlinks defined in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// read the yaml file from this directory, names symlinker.yaml

		symlinks, err := getSymlinks()
		if err != nil {
			return err
		}
		for _, link := range symlinks.Links {

			sourcePath := expandPath(link.Source)
			targetPath := expandPath(link.Target)
			// check if the link exists
			fmt.Println("Checking link:", link.Source, "->", link.Target)
			_, err := os.Lstat(sourcePath)
			if os.IsNotExist(err) {
				fmt.Printf("  Source does not exist: %s\n", link.Source)
			} else if err != nil {
				fmt.Println("  Error checking link:", err)
				return err
			}

			target, err := os.Lstat(targetPath)
			if os.IsNotExist(err) {
				fmt.Printf("  Target does not exist: %s\n", link.Target)
				continue
			} else if err != nil {
				fmt.Println("  Error checking link:", err)
				return err
			}

			// check if the target is a symlink
			if isSymlink(target) {
				// read the symlink
				actualTarget, err := os.Readlink(targetPath)
				if err != nil {
					fmt.Println("  Error reading symlink:", err)
					return err
				}
				// compare the actual target with the expected target
				if actualTarget != sourcePath {
					fmt.Printf("  Target mismatch: expected %s, got %s\n", sourcePath, actualTarget)
				}
			} else {
				fmt.Printf("  Target is not a symlink: %s\n", link.Target)
			}
		}

		// if we reach here, all links checked
		return nil
	},
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	// handle relative paths
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return ""
		}
		return filepath.Join(cwd, path)
	}
	return path
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

func getSymlinks() (entity.Symlinks, error) {
	file, err := os.ReadFile(symlinkerFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return entity.Symlinks{}, err
	}

	// parse the yaml file

	symlinks := entity.Symlinks{}

	err = yaml.Unmarshal(file, &symlinks)
	if err != nil {
		fmt.Println("Error parsing yaml:", err)
		return entity.Symlinks{}, err
	}

	// check if there are any duplicate targets
	targets := make(map[string]bool)
	for _, link := range symlinks.Links {
		if targets[link.Target] {
			return entity.Symlinks{}, fmt.Errorf("duplicate target found: %s", link.Target)
		}
		targets[link.Target] = true
	}

	return symlinks, nil
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
