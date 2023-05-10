// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/spandigitial/codeassistant/web"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		dist, err := fs.Sub(web.FileSystem, "dist")

		if err == nil {
			http.Handle(
				"/", http.FileServer(http.FS(dist)),
			)
			err = http.ListenAndServe(":8080", nil)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
