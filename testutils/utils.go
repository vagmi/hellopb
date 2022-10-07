package testutils

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

func GetTestApp() core.App {
	os.MkdirTemp("", "temp_pb_data_*")
	app := pocketbase.NewWithConfig(pocketbase.Config{
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
	log.Default().Printf("ran migrations %v", mList)
	return app
}

func CleanUpTestDir(app core.App) {
	os.RemoveAll(app.DataDir())
}
