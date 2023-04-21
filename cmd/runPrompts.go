/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"time"
)

// runPromptsCmd represents the runPrompts command
var runPromptsCmd = &cobra.Command{
	Use:   "run-prompts",
	Short: "Run prompts from prompt database",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		commandInstance, err := model.NewCommandInstance(args)
		if err != nil {
			return err
		}
		openAiApiKey := viper.GetString("openAiApiKey")
		user := viper.GetString("userEmail")
		chatGPT := client.New(openAiApiKey, rate.NewLimiter(rate.Every(60*time.Second), 20), client.WithUser(user))
		choices, err := chatGPT.Completion(commandInstance.Prompts...)
		if err != nil {
			return err
		}
		for _, choice := range choices {
			println(choice.Message.Content)
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
