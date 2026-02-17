package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes the symlinks defined in the configuration file",
	Long: `Deletes the symlinks defined in the configuration file.

Note: The target will only be deleted if it is a symlink pointing to the source.
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

			cmd.Printf("Deleting link: %s -> %s\n", link.Source, link.Target)

			targetInfo, err := os.Lstat(targetPath)
			if err != nil {
				if os.IsNotExist(err) {
					cmd.Println("  Symlink does not exist, skipping")
					continue
				}
				return fmt.Errorf("stat target %q: %w", targetPath, err)
			}

			if !isSymlink(targetInfo) {
				cmd.Printf("  Target is not a symlink: %s\n", link.Target)
				continue
			}

			ok, stored, resolved, err := symlinkPointsTo(targetPath, sourcePath)
			if err != nil {
				return fmt.Errorf("read symlink %q: %w", targetPath, err)
			}
			if !ok {
				cmd.Printf("  Target mismatch: expected %s, got %s (resolved %s)\n", sourcePath, stored, resolved)
				continue
			}

			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("remove %q: %w", targetPath, err)
			}
			cmd.Println("  Symlink deleted successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
