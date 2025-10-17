package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
