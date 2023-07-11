// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/openai"
	"github.com/spandigitial/codeassistant/client/vertexai"
	"github.com/spandigitial/codeassistant/model/indexing"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var llmClient client.LLMClient
		switch viper.GetString("backend") {
		case "openai":
			openAiApiKey := viper.GetString("openAiApiKey")
			user := viper.GetString("openAiUserId")
			userAgent := viper.GetString("userAgent")
			if userAgent == "" {
				userAgent = "SPAN Digital codeassistant"
			}
			llmClient = openai.New(openAiApiKey, debugger, openai.WithUser(user), openai.WithUserAgent(userAgent))
		case "vertexai":
			vertexAiProjectId := viper.GetString("vertexAiProjectId")
			vertexAiLocation := viper.GetString("vertexAiLocation")
			vertexAiModel := viper.GetString("vertexAiModel")
			llmClient = vertexai.New(vertexAiProjectId, vertexAiLocation, vertexAiModel, debugger)
		}

		embeddingInstance, err := indexing.NewEmbeddingInstance(llmClient, true, map[string]string{}, args[:2]...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		err = embeddingInstance.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		what := args[2]
		query := args[3]

		messageParts := make(chan client.MessagePart)
		go func() {
			err = embeddingInstance.Search(llmClient, what, query, messageParts)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		}()

		for message := range messageParts {
			if message.Type == "Part" {
				fmt.Print(message.Delta)
			} else if message.Type == "Done" {
				break
			}
		}
		fmt.Println()

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
