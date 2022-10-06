package main

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
	"github.com/vagmi/hellopb/exporter"
)

func main() {
	app := pocketbase.New()
	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Exports the collection to JSON on STDOUT",
		Run: func(cmd *cobra.Command, args []string) {
			colls, err := exporter.ExportCollections(app)
			if err != nil {
				panic(err)
			}
			fmt.Println(colls)
		},
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
