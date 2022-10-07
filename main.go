package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"github.com/vagmi/hellopb/exporter"
	"github.com/vagmi/hellopb/invitations"
)

func main() {
	app := pocketbase.New()
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
	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "invitations" {
			err := invitations.SendInvitation(app, e.Record)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
