package exporter

import (
	"bytes"
	"encoding/json"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func ExportCollections(app *pocketbase.PocketBase) (string, error) {
	var collections []models.Collection
	err := app.Dao().CollectionQuery().All(&collections)
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(collections); err != nil {
		return "", nil
	}
	return buf.String(), nil
}
