package exporter_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/vagmi/hellopb/exporter"
	_ "github.com/vagmi/hellopb/migrations"
)

var app core.App

func TestMain(m *testing.M) {

	os.RemoveAll("./test_pb_data")
	app = pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "./test_pb_data",
	})
	app.Bootstrap()
	runner, err := migrate.NewRunner(app.DB(), migrations.AppMigrations)
	if err != nil {
		log.Fatal(err)
	}
	mList, err := runner.Up()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rann the following migrations %v\n", mList)
	m.Run()
}

func TestExporter(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "collection*.json")
	fName := tmpFile.Name()
	if err != nil {
		t.Fatalf("creating temp file failed %v", err)
	}
	exporter.ExportCollections(app, tmpFile)
	tmpFile.Close()
	collJSON, err := os.Open(fName)
	if err != nil {
		t.Logf("unable to read file %s - %v", fName, err)
		t.FailNow()
	}
	defer collJSON.Close()
	type collection struct {
		Name string `json:"name"`
	}
	var collections []collection
	decoder := json.NewDecoder(collJSON)
	err = decoder.Decode(&collections)
	if err != nil {
		t.Logf("Unable to decode json %v", err)
	}
	if len(collections) != 3 {
		t.Logf("expected 3 collections found %d", len(collections))
	}
}
