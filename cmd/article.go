/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/assistant"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"time"

	"github.com/spf13/cobra"
)

// articleCmd represents the whatis command
var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		openAiApiKey := viper.GetString("openAiApiKey")
		user := viper.GetString("userEmail")
		chatGPT := client.New(openAiApiKey, rate.NewLimiter(rate.Every(60*time.Second), 20), client.WithUser(user))
		codeAssistant := assistant.New(chatGPT)

		return codeAssistant.Article(args[0], func(markdown string) {
			fmt.Println(markdown)
		})
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// articleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// articleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
