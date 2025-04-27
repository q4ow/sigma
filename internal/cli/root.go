package cli

import (
	"fmt"

	"github.com/q4ow/sigma/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigma",
	Short: "very sigma cli",
	Long:  `A very sigma command line interface`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		config.SetVerbose(verbose)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		if config.IsVerbose() {
			fmt.Println("Running in verbose mode...")
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(systemCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
