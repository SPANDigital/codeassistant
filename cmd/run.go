// SPDX-License-Identifier: MIT

package cmd

import (
	"bufio"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/openai"
	"github.com/spandigitial/codeassistant/client/vertexai"
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
	Run: func(cmd *cobra.Command, args []string) {
		commandInstance, err := model.NewCommandInstance(true, map[string]string{}, args...)
		if err == nil {
			backend := viper.GetString("backend")
			if backend == "" {
				backend = "openai"
			}
			var llmClient client.LLMClient
			switch backend {
			case "openai":
				openAiApiKey := viper.GetString("openAiApiKey")
				user := viper.GetString("userEmail")
				userAgent := viper.GetString("userAgent")
				if userAgent == "" {
					userAgent = "SPAN Digital codeassistant"
				}
				llmClient = openai.New(openAiApiKey, debugger, rate.NewLimiter(rate.Every(60*time.Second), 20), openai.WithUser(user), openai.WithUserAgent(userAgent))
			case "vertexai":
				vertexAiProjectId := viper.GetString("vertexAiProjectId")
				vertexAiLocation := viper.GetString("vertexAiLocation")
				llmClient = vertexai.New(vertexAiProjectId, vertexAiLocation, debugger)
			}
			f := bufio.NewWriter(os.Stdout)
			defer f.Flush()
			messages := make(chan client.MessagePart)
			go func() {
				err = llmClient.Completion(commandInstance, messages)
			}()
			for message := range messages {
				if message.Type == "Part" {
					fmt.Fprint(f, message.Delta)
				}
			}
			fmt.Fprintln(f)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
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
