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

			sourcePath := expandPath(link.Source)
			targetPath := expandPath(link.Target)
			// check if the link exists
			fmt.Println("Creating link:", link.Source, "->", link.Target)
			_, err := os.Lstat(sourcePath)
			if os.IsNotExist(err) {
				fmt.Printf("  Skipping because source does not exist: %s\n", link.Source)
				continue
			} else if err != nil {
				fmt.Println("  Error creating link:", err)
				return err
			}

			target, err := os.Lstat(targetPath)
			if os.IsNotExist(err) {
				err := os.Symlink(sourcePath, targetPath)
				if err != nil {
					fmt.Println("  Error creating symlink:", err)
					return err
				}
				fmt.Println("  Symlink created")
				continue
			} else if err != nil {
				fmt.Println("  Error creating symlink:", err)
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

		// end of RunE
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
