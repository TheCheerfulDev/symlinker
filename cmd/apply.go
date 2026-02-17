package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the symlinks defined in the configuration file",
	Long: `Applies the symlinks defined in the configuration file.
If the source doe not exist, the link is skipped.
If the target exists and is a symlink, it is verified.
If the target exists and is not a symlink, a warning is printed and the link is skipped.

Note: The target will NEVER be overwritten or deleted by this command.
`,
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

			cmd.Printf("Creating link: %s -> %s\n", link.Source, link.Target)

			if _, err := os.Lstat(sourcePath); err != nil {
				if os.IsNotExist(err) {
					cmd.Printf("  Skipping because source does not exist: %s\n", link.Source)
					continue
				}
				return fmt.Errorf("stat source %q: %w", sourcePath, err)
			}

			targetInfo, err := os.Lstat(targetPath)
			if err != nil {
				if os.IsNotExist(err) {
					if err := os.Symlink(sourcePath, targetPath); err != nil {
						return fmt.Errorf("create symlink %q -> %q: %w", targetPath, sourcePath, err)
					}
					cmd.Println("  Symlink created")
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
	rootCmd.AddCommand(applyCmd)
}
