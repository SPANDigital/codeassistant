// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model/prompts"
	"github.com/spf13/cobra"
)

// listPromptsCmd represents the listPrompts command
var listPromptsCmd = &cobra.Command{
	Use:   "list",
	Short: "List prompts from prompt database",
	Long:  `List prompts from prompt database.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, library := range prompts.BuildLibraries() {
			for _, command := range library.Commands {
				if !command.Abstract {
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
					fmt.Printf("codeassistant run %s %s%s\n", library.Name, command.Name, params)
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
