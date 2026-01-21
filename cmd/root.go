/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/vfleuryrodeofx/gotractor/internal/pkg/requests"
	"github.com/vfleuryrodeofx/gotractor/internal/pkg/ui"
	"github.com/vfleuryrodeofx/gotractor/internal/pkg/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotractor",
	Short: "A simple CLI to query tractor 🚜, powered by Go !",
	Long: `Go tractor lets you query Tractor jobs, tasks, retrieve logs and display
them in a nice TUI.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚜")
		url := args[0]
		slog.Info("Url passed is ", "url", url)
		jid := requests.ExtractJID(url)

		data, tasksData, err := requests.GetTaskTree(jid)
		if err != nil {
			return fmt.Errorf("Could not get the task tree for %s, err : %w", jid, err)
		}
		utils.GetListFromTreeTask(tasksData)
		ui.Show(data, tasksData, jid)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
