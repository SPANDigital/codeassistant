/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/assistant"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
	"os"
	"time"
)

var schemaFilename string
var entitiesDirectory string
var servicesDirectory string

// schema2entitiesCmd represents the schema2entities command
var schema2entitiesCmd = &cobra.Command{
	Use:   "schema2entities",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schema2entities called")

		data, err := os.ReadFile(schemaFilename)
		if err != nil {
			panic(err)
		}

		// convert the []byte buffer to a string
		railsSchema := string(data)

		chatGPT := client.New(chatGptApiKey, rate.NewLimiter(rate.Every(60*time.Second), 20))
		codeAssistant := assistant.New(chatGPT)

		codeAssistant.RailsSchemaToEntities(railsSchema, model.SourceCodeHandlers(func(code model.SourceCode) model.SourceCode {
			if code.Language == "typescript" && code.Content != "" {
				code.Save(entitiesDirectory)
			}
			return code
		}), model.SourceCodeHandlers(func(code model.SourceCode) model.SourceCode {
			if code.Language == "typescript" && code.Content != "" {
				code.Save(servicesDirectory)
			}
			return code
		}))
	},
}

func init() {
	rails2nextjsCmd.AddCommand(schema2entitiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schema2entitiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schema2entitiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	schema2entitiesCmd.Flags().StringVar(&schemaFilename, "schemaFilename", "", "schema file name")
	schema2entitiesCmd.Flags().StringVar(&entitiesDirectory, "entitiesDirectory", "", "entities directory")
	schema2entitiesCmd.Flags().StringVar(&servicesDirectory, "servicesDirectory", "", "services directory")
}
