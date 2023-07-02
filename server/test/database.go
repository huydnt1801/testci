package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/enttest"
	"github.com/huydnt1801/chuyende/internal/ent/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Database *sql.DB
	Client   *ent.Client
}

const (
	dbUsername = "root"
	dbPassword = "password"
)

// mysqlContainer creates an instance of the mysql container type
func MysqlContainer(ctx context.Context) (testcontainers.Container, func()) {
	req := testcontainers.ContainerRequest{
		Name:         "test_mysql",
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": dbPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		panic(err)
	}

	return container, func() {
		container.Terminate(ctx)
	}
}

func NewTestEnv(t *testing.T, c testcontainers.Container, dbName string) *TestEnv {
	var (
		env = &TestEnv{}
		err error
	)
	ctx := context.Background()
	host, _ := c.Host(ctx)

	p, _ := c.MappedPort(ctx, "3306/tcp")
	port := p.Int()
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&parseTime=true",
		dbUsername, dbPassword, host, port, dbName)

	// Create connection and create new database for eact test case/test suite
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		dbUsername,
		dbPassword,
		host,
		port)
	conn, err := sql.Open("mysql", uri)
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		t.Fatal(err)
	}

	// Mock database schema
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}
	env.Client = enttest.Open(t, "mysql", dbUrl, opts...)

	err = env.Client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(false), // Disable foreign keys.
	)
	if err != nil {
		return nil
	}

	// Create db connection
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil
	}
	env.Database = db
	return env
}

func (e *TestEnv) CleanUp(a *assert.Assertions) {
	ctx := context.Background()
	if _, err := e.Client.Trip.Delete().Exec(ctx); err != nil {
		a.FailNow(err.Error(), "failed to truncate trip")
	}
	if _, err := e.Client.Session.Delete().Exec(ctx); err != nil {
		a.FailNow(err.Error(), "failed to truncate session")
	}
	if _, err := e.Client.Otp.Delete().Exec(ctx); err != nil {
		a.FailNow(err.Error(), "failed to truncate otp")
	}
	if _, err := e.Client.User.Delete().Exec(ctx); err != nil {
		a.FailNow(err.Error(), "failed to truncate user")
	}
	if _, err := e.Client.VehicleDriver.Delete().Exec(ctx); err != nil {
		a.FailNow(err.Error(), "failed to truncate driver")
	}
}
