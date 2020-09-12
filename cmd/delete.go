package cmd
import (
"github.com/spf13/cobra"
"uosDevops/config"
"uosDevops/install"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Simplest way to delete your project",
	Long:  `uosdevops delete--project projectname`,
	Run: func(cmd *cobra.Command, args []string) {
		install.Delete()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&config.ProjectName, "project", "hellworld", "For project delete")
}