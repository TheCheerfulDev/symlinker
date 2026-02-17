package cmd

import (
	"fmt"
	"os"
	"symlinker/entity"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default symlinker.yaml in the current directory",
	Long:  "Creates a starter symlinker.yaml configuration file if it doesn't already exist.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stat(symlinkerFile)
		if err == nil {
			if fi.IsDir() {
				return fmt.Errorf("%s exists and is a directory", symlinkerFile)
			}
			cmd.Printf("%s already exists\n", symlinkerFile)
			return nil
		}
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", symlinkerFile, err)
		}

		defaultContent := entity.Symlinks{
			Links: []entity.Symlink{
				{
					Source: "/path/to/source",
					Target: "/path/to/target",
				},
			},
		}

		out, err := yaml.Marshal(defaultContent)
		if err != nil {
			return fmt.Errorf("marshal default config: %w", err)
		}

		if err := os.WriteFile(symlinkerFile, out, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", symlinkerFile, err)
		}
		cmd.Printf("%s created successfully\n", symlinkerFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
