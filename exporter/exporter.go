package exporter

import (
	"encoding/json"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func ExportCollections(app core.App, filename string) error {
	var collections []models.Collection
	err := app.Dao().CollectionQuery().All(&collections)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
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
