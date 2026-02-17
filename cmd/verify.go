package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the symlinks defined in the configuration file",
	Long:  `Verifies the symlinks defined in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		symlinks, err := getSymlinks()
		if err != nil {
			return err
		}

		for _, link := range symlinks.Links {
			sourcePath, err := expandPath(link.Source)
			if err != nil {
				return fmt.Errorf("expand source %q: %w", link.Source, err)
			}
			targetPath, err := expandPath(link.Target)
			if err != nil {
				return fmt.Errorf("expand target %q: %w", link.Target, err)
			}

			cmd.Printf("Checking link: %s -> %s\n", link.Source, link.Target)

			if _, err := os.Lstat(sourcePath); err != nil {
				if os.IsNotExist(err) {
					cmd.Printf("  Source does not exist: %s\n", link.Source)
				} else {
					return fmt.Errorf("stat source %q: %w", sourcePath, err)
				}
			}

			targetInfo, err := os.Lstat(targetPath)
			if err != nil {
				if os.IsNotExist(err) {
					cmd.Printf("  Target does not exist: %s\n", link.Target)
					continue
				}
				return fmt.Errorf("stat target %q: %w", targetPath, err)
			}

			if isSymlink(targetInfo) {
				ok, stored, resolved, err := symlinkPointsTo(targetPath, sourcePath)
				if err != nil {
					return fmt.Errorf("read symlink %q: %w", targetPath, err)
				}
				if !ok {
					cmd.Printf("  Target mismatch: expected %s, got %s (resolved %s)\n", sourcePath, stored, resolved)
				}
			} else {
				cmd.Printf("  Target is not a symlink: %s\n", link.Target)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
