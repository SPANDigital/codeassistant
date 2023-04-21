/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spf13/cobra"
)

// listPromptsCmd represents the listPrompts command
var listPromptsCmd = &cobra.Command{
	Use:   "list-prompts",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, library := range model.Libraries {
			for _, command := range library.Commands {
				if !command.System {
					params := ""
					for param, value := range command.Params {
						var display string
						if value == "" {
							display = "<value>"
						} else {
							display = fmt.Sprintf("<value, default: %s>", value)
						}
						params = fmt.Sprintf("%s %s:%s>", params, param, display)
					}
					fmt.Printf("%s %s%s\n", library.Name, command.Name, params)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listPromptsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listPromptsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listPromptsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
