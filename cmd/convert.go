/*
Copyright © 2023 richard.wooding@spandigital.com
*/
package cmd

import (
	"fmt"
	"github.com/spandigitial/codeassistant/assistant"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/model"
	"golang.org/x/time/rate"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var railsType string
var nestJsType string
var src string
var dest string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("convert called")

		chatGPT := client.New(chatGptApiKey, rate.NewLimiter(rate.Every(60*time.Second), 20))
		codeAssistant := assistant.New(chatGPT)

		reader, err := os.Open(src)
		if err != nil {
			return err
		}

		absDest, _ := filepath.Abs(dest)

		println("Abs dest is ", absDest)

		return codeAssistant.Convert(reader, railsType, nestJsType, absDest, model.SourceCodeHandlers(func(code model.SourceCode) model.SourceCode {
			code.Save(path.Dir(absDest))
			return code
		}))
	},
}

func init() {
	rails2nextjsCmd.AddCommand(convertCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	convertCmd.Flags().StringVar(&src, "src", "", "Source filename")
	convertCmd.Flags().StringVar(&dest, "dest", "", "Destination filename")
	convertCmd.Flags().StringVar(&railsType, "railstype", "", "Ruby On Rails type")
	convertCmd.Flags().StringVar(&nestJsType, "nestjstype", "", "NestJS type")
}
