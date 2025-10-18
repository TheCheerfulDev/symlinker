package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var symlinkerFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "symlinker",
	Short: "A brief description of your application",
	Long: `SYmlinker is a tool to manage symlinks based on a configuration file.
You can use it to create, delete, and verify symlinks defined in a YAML configuration file.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&symlinkerFile, "file", "symlinker.yaml", "Symlink file (default is symlinker.yaml)")
}
