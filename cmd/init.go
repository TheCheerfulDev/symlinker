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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if symlinkerFile exists
		_, err := os.Stat(symlinkerFile)
		if os.IsNotExist(err) {
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
				panic(err)
			}

			err = os.WriteFile(symlinkerFile, out, 0644)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s created successfully\n", symlinkerFile)
			return
		} else if err != nil {
			panic(err)
		}

		fmt.Printf("%s already exists\n", symlinkerFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
