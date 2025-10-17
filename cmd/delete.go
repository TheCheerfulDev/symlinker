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

			sourcePath := expandPath(link.Source)
			targetPath := expandPath(link.Target)
			// check if the link exists
			fmt.Println("Deleting link:", link.Source, "->", link.Target)

			target, err := os.Lstat(targetPath)
			if os.IsNotExist(err) {
				fmt.Println("  Symlink does not exist, skipping")
				continue
			} else if err != nil {
				fmt.Println("  Error deleting symlink:", err)
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
					continue
				}

				// remove the symlink
				err = os.Remove(targetPath)
				if err != nil {
					fmt.Println("  Error deleting symlink:", err)
					return err
				}
				fmt.Println("  Symlink deleted successfully")
			} else {
				fmt.Printf("  Target is not a symlink: %s\n", link.Target)
			}
		}

		// end of RunE
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
