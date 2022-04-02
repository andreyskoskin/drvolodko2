package main

import (
	"io"
	"log"

	"github.com/andreyskoskin/drvolodko2/api"
	"github.com/andreyskoskin/drvolodko2/api/auditor"
	"github.com/andreyskoskin/drvolodko2/datamodel"
	"github.com/andreyskoskin/drvolodko2/datasource"

	"github.com/andreyskoskin/drvolodko2/api/auditprogram"
)

var config = Config{
	Echo: api.EchoConfig{
		Address: "localhost:8080",
	},
}

type DataSource interface {
	io.Closer
	AuditPrograms() datamodel.AuditPrograms
	Auditors() datamodel.Auditors
}

func main() {
	// TODO: load config from file

	var ds, err = initDataSource(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = ds.Close()
	}()

	err = api.StartEcho(config.Echo, api.EchoBindings{
		"/audit-program": auditprogram.NewEchoAPI(ds.AuditPrograms()),
		"/auditor":       auditor.NewEchoAPI(ds.Auditors()),
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func initDataSource(c Config) (DataSource, error) {
	if c.Postgres.DBName == "" {
		return datasource.NewInMemory(), nil
	}

	return datasource.NewPostgres(c.Postgres)
}
