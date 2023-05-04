// SPDX-License-Identifier: MIT

package cmd

import (
	"bufio"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"os"
	"time"
)

// runPromptsCmd represents the runPrompts command
var runPromptsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run prompts from prompt database",
	Long:  `Run prompts from prompt database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		commandInstance, err := model.NewCommandInstance(args)
		if err != nil {
			return err
		}
		openAiApiKey := viper.GetString("openAiApiKey")
		user := viper.GetString("userEmail")
		chatGPT := client.New(openAiApiKey, rate.NewLimiter(rate.Every(60*time.Second), 20), client.WithUser(user))
		choices, err := chatGPT.Completion(commandInstance)
		if err != nil {
			return err
		}
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		for _, choice := range choices {
			fmt.Fprintln(f, choice.Message.Content)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runPromptsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runPromptsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runPromptsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
