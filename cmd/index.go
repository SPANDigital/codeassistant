// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/openai"
	"github.com/spandigitial/codeassistant/client/vertexai"
	"github.com/spandigitial/codeassistant/model/indexing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// embeddingsCmd represents the saveEmbeds command
var embeddingsCmd = &cobra.Command{
	Use:   "index",
	Short: "Create an index based on vector indexing",
	Long:  `Create an index based on vector indexing`,
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

		embeddingInstance, err := indexing.NewEmbeddingInstance(llmClient, true, map[string]string{}, args...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		err = embeddingInstance.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		err = embeddingInstance.Collect()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		err = embeddingInstance.Fetch()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		err = embeddingInstance.Save()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(embeddingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveEmbedsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveEmbedsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
