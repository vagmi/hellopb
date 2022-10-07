package exporter

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func ExportCollections(app core.App, writeTo io.Writer) error {
	var collections []models.Collection
	err := app.Dao().CollectionQuery().All(&collections)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writeTo)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(collections); err != nil {
		return err
	}
	return nil
}

func ImportCollections(app core.App, filename string) error {
	var collections []*models.Collection
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&collections); err != nil {
		return err
	}
	importer := forms.NewCollectionsImport(app)
	importer.Collections = collections
	return importer.Submit()
}
