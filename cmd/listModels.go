/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/openai"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// listModelsCmd represents the listModels command
var listModelsCmd = &cobra.Command{
	Use:   "list-models",
	Short: "list-models",
	Long:  `list-models`,
	Run: func(cmd *cobra.Command, args []string) {
		openAiApiKey := viper.GetString("openAiApiKey")
		user := viper.GetString("userEmail")
		userAgent := viper.GetString("userAgent")
		if userAgent == "" {
			userAgent = "SPAN Digital codeassistant"
		}
		chatGPT := openai.New(openAiApiKey, debugger, rate.NewLimiter(rate.Every(60*time.Second), 20), openai.WithUser(user), openai.WithUserAgent(userAgent))
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		models := make(chan client.LanguageModel)
		err := chatGPT.Models(models)
		for model := range models {
			fmt.Fprintln(os.Stdout, model)
		}
		/*func(languageModel openai.LanguageModel) {

		})*/
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listModelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listModelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listModelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
