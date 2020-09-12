package cmd

import (
	"github.com/spf13/cobra"
	"uosDevops/config"
	"uosDevops/install"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Simplest way to build your project",
	Long:  `uosdevops build --project projectname`,
	Run: func(cmd *cobra.Command, args []string) {
		install.Build()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&config.ProjectName, "project", "hellworld", "For project build")
}