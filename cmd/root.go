// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	debugger2 "github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/slices"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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
		debugger = debugger2.New(slices.MapSlice(viper.GetStringSlice("debug"), func(s string) debugger2.Detail {
			detail, _ := debugger2.Parse(s)
			return detail
		})...)
		debugger.MessageF(debugger2.Configuration, "%v", viper.AllSettings())
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
	userId := ""
	if found, err := user.Current(); err == nil {
		userId = found.Username
	}
	rootCmd.PersistentFlags().String("openAiUserId", userId, "User to send to OpenAI")
	if err := viper.BindPFlag("openAiUserId", rootCmd.PersistentFlags().Lookup("openAiUserId")); err != nil {
		log.Fatal("Unable to find flag userEmail", err)
	}
	rootCmd.PersistentFlags().String("openAiModel", "gpt-4", "Model to use if not specified")
	if err := viper.BindPFlag("openAiModel", rootCmd.PersistentFlags().Lookup("openAiModel")); err != nil {
		log.Fatal("Unable to find flag openAiModel", err)
	}
	rootCmd.PersistentFlags().String("openAiUrlPrefix", "https://api.openai.com", "Prefix of OpenAI Urls")
	if err := viper.BindPFlag("openAiUrlPrefix", rootCmd.PersistentFlags().Lookup("openAiUrlPrefix")); err != nil {
		log.Fatal("Unable to find flag openAiUrlPrefix", err)
	}
	var gcloudBinary = "gcloud"
	if found, err := exec.LookPath("gcloud"); err == nil {
		gcloudBinary = found
	}
	rootCmd.PersistentFlags().String("gcloudBinary", gcloudBinary, "Gcloud Binary")
	if err := viper.BindPFlag("gcloudBinary", rootCmd.PersistentFlags().Lookup("gcloudBinary")); err != nil {
		log.Fatal("Unable to find flag gcloudBinary", err)
	}
	rootCmd.PersistentFlags().String("vertexAiProjectId", "", "Vertex Project ID")
	if err := viper.BindPFlag("vertexAiProjectId", rootCmd.PersistentFlags().Lookup("vertexAiProjectId")); err != nil {
		log.Fatal("Unable to find flag vertexAiProjectId", err)
	}
	rootCmd.PersistentFlags().String("vertexAiModel", "text-bison@001", "Model to use if not specified")
	if err := viper.BindPFlag("vertexAiModel", rootCmd.PersistentFlags().Lookup("vertexAiModel")); err != nil {
		log.Fatal("Unable to find flag vertexAiModel", err)
	}
	rootCmd.PersistentFlags().String("vertexAiLocation", "us-central1", "Locstion to use if not specified")
	if err := viper.BindPFlag("vertexAiLocation", rootCmd.PersistentFlags().Lookup("vertexAiLocation")); err != nil {
		log.Fatal("Unable to find flag vertexAiLocation", err)
	}
	// Find home directory.
	promptsLibraryDir := ""
	home, err := os.UserHomeDir()
	if err == nil {
		promptsLibraryDir = filepath.Join(home, "prompts-library")
	}
	rootCmd.PersistentFlags().String("promptsLibraryDir", promptsLibraryDir, "Prompts library Dir")
	if err := viper.BindPFlag("promptsLibraryDir", rootCmd.PersistentFlags().Lookup("promptsLibraryDir")); err != nil {
		log.Fatal("Unable to find flag promptsLibraryDir", err)
	}
	rootCmd.PersistentFlags().String("userAgent", "SPANDigital codeassistant", "HTTP User-Agent")
	if err := viper.BindPFlag("userAgent", rootCmd.PersistentFlags().Lookup("userAgent")); err != nil {
		log.Fatal("Unable to find flag userAgent", err)
	}
	rootCmd.PersistentFlags().StringArray("debug", nil, "details to debug, recognized values: request-header,sent-prompt,request-payload,response-header,request-time,first-reponse-time,last-response-time")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatal("Unable to find flag debugDetails", err)
	}
	rootCmd.PersistentFlags().String("backend", "openai", "backend openai or vertexai")
	if err := viper.BindPFlag("backend", rootCmd.PersistentFlags().Lookup("backend")); err != nil {
		log.Fatal("Unable to find flag backend", err)
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
