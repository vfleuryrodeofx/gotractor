/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/victorfleury/gotractor/internal/pkg/requests"
	"github.com/victorfleury/gotractor/internal/pkg/ui"
	"github.com/victorfleury/gotractor/internal/pkg/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotractor",
	Short: "A simple CLI to query tractor 🚜, powered by Go !",
	Long: `Go tractor lets you query Tractor jobs, tasks, retrieve logs and display
them in a nice TUI.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚜")
		if len(args) != 1 {
			slog.Warn("You need to provide 1 jid or URL with a jid")
			os.Exit(1)
		}
		url := args[0]
		slog.Info("Url passed is ", "url", url)
		jid := requests.ExtractJID(url)

		data, tasksData := requests.GetTaskTree(jid)
		utils.GetListFromTreeTask(tasksData)
		ui.Show(data, tasksData, jid)
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotractor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
