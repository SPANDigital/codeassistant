// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	debugger2 "github.com/spandigitial/codeassistant/client/debugger"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugger *debugger2.Debugger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codeassistant",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debugger = debugger2.New(viper.GetStringSlice("debug")...)
		if debugger.IsRecording("configuration") {
			debugger.Message("configuration", fmt.Sprintf("%v", viper.AllSettings()))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codeassistant.yaml)")
	rootCmd.PersistentFlags().String("openAiApiKey", "", "OpenAI API Key")
	if err := viper.BindPFlag("openAiApiKey", rootCmd.PersistentFlags().Lookup("openAiApiKey")); err != nil {
		log.Fatal("Unable to find flag openAiApiKey", err)
	}
	rootCmd.PersistentFlags().String("defaultOpenAiModel", "", "Model to use if not specified (defaqults to gpt=3.5-turbo)")
	if err := viper.BindPFlag("defaultOpenAiModel", rootCmd.PersistentFlags().Lookup("defaultOpenAiModel")); err != nil {
		log.Fatal("Unable to find flag defaultOpenAiModel", err)
	}
	rootCmd.PersistentFlags().String("vertexAiApiKey", "", "Vertex API Key")
	if err := viper.BindPFlag("vertexAiApiKey", rootCmd.PersistentFlags().Lookup("vertexAiApiKey")); err != nil {
		log.Fatal("Unable to find flag vertexAiApiKey", err)
	}
	rootCmd.PersistentFlags().String("defaultVertexAiModel", "", "Model to use if not specified (defaqults to text-bison@001")
	if err := viper.BindPFlag("defaultVertexAiModel", rootCmd.PersistentFlags().Lookup("defaultVertexAiModel")); err != nil {
		log.Fatal("Unable to find flag defaultVertexAiModel", err)
	}
	rootCmd.PersistentFlags().String("defaultVertexAiLocation", "", "Locstion to use if not specified (defaqults to us-central1")
	if err := viper.BindPFlag("defaultVertexAiLocation", rootCmd.PersistentFlags().Lookup("defaultVertexAiLocation")); err != nil {
		log.Fatal("Unable to find flag defaultVertexAiLocation", err)
	}
	rootCmd.PersistentFlags().String("userEmail", "", "User to send to ChatGPT")
	if err := viper.BindPFlag("userEmail", rootCmd.PersistentFlags().Lookup("userEmail")); err != nil {
		log.Fatal("Unable to find flag userEmail", err)
	}
	rootCmd.PersistentFlags().String("promptsLibraryDir", "", "Prompt library Dir")
	if err := viper.BindPFlag("promptsLibraryDir", rootCmd.PersistentFlags().Lookup("promptsLibraryDir")); err != nil {
		log.Fatal("Unable to find flag promptsLibraryDir", err)
	}
	rootCmd.PersistentFlags().String("userAgent", "", "HTTP User-Agent (default is SPANDigital codeassistant)")
	if err := viper.BindPFlag("userAgent", rootCmd.PersistentFlags().Lookup("userAgent")); err != nil {
		log.Fatal("Unable to find flag userAgent", err)
	}
	rootCmd.PersistentFlags().StringArray("debug", nil, "details to debug, recognized values: request-header,sent-prompt,request-payload,response-header,request-time,first-reponse-time,last-response-time")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatal("Unable to find flag debugDetails", err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".codeassistant" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".codeassistant")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
