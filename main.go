package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"github.com/vagmi/hellopb/exporter"
	"github.com/vagmi/hellopb/invitations"
	_ "github.com/vagmi/hellopb/migrations"
)

func main() {
	app := pocketbase.New()
	addCommands(app)
	setupHooks(app)
	addEndpoints(app)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func addEndpoints(app core.App) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/hello",
			Handler: func(c echo.Context) error {
				return c.String(200, "Hello world!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireAdminOrUserAuth(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/websocket",
			Handler: func(c echo.Context) error {
				upgrader := websocket.Upgrader{}
				ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
				if err != nil {
					return err
				}
				msg := []byte("hello client")
				defer ws.Close()
				for {
					err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Server %s", string(msg))))
					if err != nil {
						log.Default().Println(err)
					}

					// Read
					_, msg, err = ws.ReadMessage()
					if err != nil {
						log.Default().Println(err)
					}
					fmt.Printf("%s\n", msg)
				}
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireAdminOrUserAuth(),
			},
		})
		return nil

	})
}

func setupHooks(app core.App) {
	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "invitations" {
			err := invitations.SendInvitation(app, e.Record)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func addCommands(app *pocketbase.PocketBase) {
	var filename string
	exportCmd := &cobra.Command{
		Use:       "export",
		Short:     "Exports the collection to JSON on STDOUT",
		ValidArgs: []string{"filename"},
		Run: func(cmd *cobra.Command, args []string) {
			collectionFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatal(err)
			}
			defer collectionFile.Close()
			err = exporter.ExportCollections(app, collectionFile)
			if err != nil {
				panic(err)
			}
		},
	}

	exportCmd.Flags().StringVarP(&filename, "output", "o", "collections.json", "Export to file")

	importCmd := &cobra.Command{
		Use:       "import",
		Short:     "Imports the collection from JSON",
		ValidArgs: []string{"filename"},
		Run: func(cmd *cobra.Command, args []string) {
			err := exporter.ImportCollections(app, filename)
			if err != nil {
				panic(err)
			}
		},
	}
	importCmd.Flags().StringVarP(&filename, "from", "f", "collections.json", "Import from file")
	app.RootCmd.AddCommand(exportCmd)
	app.RootCmd.AddCommand(importCmd)
}
