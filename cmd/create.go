package cmd

import (
"github.com/spf13/cobra"
"uosDevops/config"
"uosDevops/install"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Simplest way to create your project",
	Long:  `uosdevops create --project projectname`,
	Run: func(cmd *cobra.Command, args []string) {
		install.Create()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&config.ProjectName, "project", "hellworld", "For project build")
}