package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

const rootDescription = "khayyam is a tool for managing uniform TypeScript monorepos."

var rootCmd = &cobra.Command{
	Use:          "khayyam",
	Short:        rootDescription,
	Long:         rootDescription,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		internal.Logger.Error(err.Error())
		os.Exit(1)
	}
}

func mustLoadRepository(checkEntryPoints bool) *internal.Repository {
	cwd, err := os.Getwd()
	if err != nil {
		internal.Logger.ErrorObj(err)
		os.Exit(1)
	}
	repo, err := internal.LoadRepository(cwd, checkEntryPoints)
	if err != nil {
		internal.Logger.ErrorObj(err)
		os.Exit(1)
	}
	return repo
}
