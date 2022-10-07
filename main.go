package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"github.com/vagmi/hellopb/exporter"
	"github.com/vagmi/hellopb/invitations"
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

func addEndpoints(app *pocketbase.PocketBase) {
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

		return nil
	})
}

func setupHooks(app *pocketbase.PocketBase) {
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
			err := exporter.ExportCollections(app, filename)
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
