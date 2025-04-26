package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "your-cli-name",
	Short: "A brief description of your application",
	Long:  `A longer description...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			fmt.Println("Running in verbose mode...")
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}

func Execute() error {
	return rootCmd.Execute()
}
