/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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
		// check if symlinker.yaml exists
		_, err := os.Stat("symlinker.yaml")
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

			err = os.WriteFile("symlinker.yaml", out, 0644)
			if err != nil {
				panic(err)
			}
			fmt.Println("symlinker.yaml created successfully")
			return
		} else if err != nil {
			panic(err)
		}

		fmt.Println("symlinker.yaml already exists")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
