// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/openai"
	"github.com/spandigitial/codeassistant/client/vertexai"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spandigitial/codeassistant/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var responses = make(map[uuid.UUID]chan client.MessagePart)

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

		libraries := model.BuildLibraries()

		dist, err := fs.Sub(web.FileSystem, "dist")
		if err == nil {
			router := gin.Default()
			if !debugger.IsRecording("webserver") {
				gin.SetMode(gin.ReleaseMode)
			} else {
				gin.SetMode(gin.DebugMode)
			}
			router.GET("/", func(context *gin.Context) {
				context.Redirect(http.StatusTemporaryRedirect, "/web")
			})
			httpFs := http.FS(dist)
			router.StaticFS("/web", httpFs)
			router.GET("/api/graph", func(context *gin.Context) {
				context.JSON(http.StatusOK, libraries)
			})
			router.GET("/api/receive/:uuid", func(c *gin.Context) {
				uuid, err := uuid.Parse(c.Param("uuid"))
				if err != nil {
					c.Error(err)
					return
				}
				_, found := responses[uuid]
				if !found {
					c.AbortWithStatus(404)
					return
				}
				c.Stream(func(w io.Writer) bool {
					// Stream message to client from message channel
					if msg, ok := <-responses[uuid]; ok {
						c.SSEvent("message", msg)
						return true
					}
					delete(responses, uuid)
					return false
				})
			})
			router.POST("/api/prompt/:libraryName/:commandName", func(c *gin.Context) {
				defaultParams := make(map[string]string)
				params, err := model.CommandInstanceParams(c.Param("libraryName"), c.Param("commandName"))
				if err != nil {
					c.Error(err)
					return
				}
				for _, param := range params {
					defaultParams[param] = c.PostForm(param)
				}
				commandInstance, err := model.NewCommandInstance(false, defaultParams, c.Param("libraryName"), c.Param("commandName"))
				if err != nil {
					fmt.Fprintln(os.Stderr, "Can't find command", err.Error())
					c.Error(err)
					return
				}
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
					llmClient = openai.New(openAiApiKey, debugger, openai.WithUser(user), openai.WithUserAgent(userAgent))
				case "vertexai":
					vertexAiProjectId := viper.GetString("vertexAiProjectId")
					vertexAiLocation := viper.GetString("vertexAiLocation")
					vertexAiModel := viper.GetString("vertexAiModel")
					llmClient = vertexai.New(vertexAiProjectId, vertexAiLocation, vertexAiModel, debugger)
				}
				uuid := uuid.New()
				messageParts := make(chan client.MessagePart)
				responses[uuid] = messageParts
				go func() {
					err = llmClient.Completion(commandInstance, messageParts)
				}()
				if err == nil {
					c.Header("Location", fmt.Sprintf("/api/receive/%s", uuid))
					c.Status(201)
				} else {
					c.Error(err)
				}

			})

			port := viper.GetInt("serverHttpPort")
			fmt.Fprintf(os.Stderr, "Visit http://0.0.0.0:%d/\n", port)
			router.Run(fmt.Sprintf(":%d", port))
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
	rootCmd.Flags().Int("serverHttpPort", 8989, "Server HTTP Port")
	if err := viper.BindPFlag("serverHttpPort", rootCmd.Flags().Lookup("serverHttpPort")); err != nil {
		log.Fatal("Unable to find flag serverHttpPort", err)
	}
}
