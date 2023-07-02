package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	atlas "ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"

	atlasschema "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/huydnt1801/chuyende/internal/config"
	"github.com/huydnt1801/chuyende/internal/ent/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/urfave/cli/v2"
	"github.com/xo/dburl"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create a migration",
				Action: func(cCtx *cli.Context) error {
					name := cCtx.Args().Get(0)
					return create(name)
				},
			},
			{
				Name:    "apply",
				Aliases: []string{"a"},
				Usage:   "apply migrations",
				Action: func(cCtx *cli.Context) error {
					return apply()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func create(name string) error {
	if name == "" {
		return errors.New("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	name = strings.ReplaceAll(name, " ", "_")

	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("./migrations")
	if err != nil {
		dir, err = atlas.NewLocalDir("../../migrations")
	}
	if err != nil {
		return fmt.Errorf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.MySQL),           // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	dburl, close := mysqlTestContainer()
	defer close()
	err = migrate.NamedDiff(context.Background(), dburl, name, opts...)
	if err != nil {
		return fmt.Errorf("failed generating migration file: %v", err)
	}
	return nil
}

// TODO: incomplete
func apply() error {
	cfg := config.MustParseConfig()
	db, err := dburl.Open(cfg.DatabaseURI())
	if err != nil {
		return fmt.Errorf("failed opening db: %v", err)
	}
	driver, err := mysql.Open(db)
	if err != nil {
		return fmt.Errorf("failed opening atlas driver: %v", err)
	}

	ctx := context.Background()
	// Inspect the created table.
	sch, err := driver.InspectSchema(ctx, "", &atlasschema.InspectOptions{
		// Tables: []string{"example"},
	})
	if err != nil {
		return fmt.Errorf("failed inspect atlas schema")
	}

	_ = sch

	fmt.Println("!!! WARNING: Apply command is NOT completed yet !!!!")

	return nil

	// driver.ApplyChanges(ctx, sch)
}

func mysqlTestContainer() (string, func()) {
	ctx := context.Background()

	dbUsername := "root"
	dbPassword := "password"
	dbName := "testdb"
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": dbPassword,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL."),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed starting mysql testcontainter: %v", err)
	}

	host, _ := mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	port := p.Int()

	dburl := fmt.Sprintf("mysql://%s:%s@%s:%d/%s?parseTime=true",
		dbUsername, dbPassword, host, port, dbName)

	return dburl, func() {
		mysqlC.Terminate(ctx)
	}
}
