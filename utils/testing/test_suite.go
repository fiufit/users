package testing

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuite struct {
	DB       *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	models   []interface{}
}

func NewTestSuite(models ...interface{}) TestSuite {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=testpgpassword",
			"POSTGRES_USER=testpguser",
			"POSTGRES_DB=testpgdb",
			"listen_addresses= '*'",
		},
	},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://testpguser:testpgpassword@%s/testpgdb?sslmode=disable", hostPort)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	var sqlDB *sql.DB
	if err = pool.Retry(func() error {
		sqlDB, err = sql.Open("pgx", databaseUrl)
		if err != nil {
			return err
		}

		return sqlDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("Could not create gorm DB from dockertest sql connection: %s", err)
	}

	suite := TestSuite{
		DB:       gormDB,
		pool:     pool,
		resource: resource,
		models:   models,
	}
	suite.setUpModels()

	return suite
}

func (ts TestSuite) setUpModels() {
	if err := ts.DB.AutoMigrate(ts.models...); err != nil {
		log.Fatalf("Could not migrate test models: %s", err)
	}
}

func (ts TestSuite) TearDown() {
	if err := ts.pool.Purge(ts.resource); err != nil {
		log.Fatalf("Could not TearDown testing DB: %s", err)
	}
}

func (ts TestSuite) TruncateModels() {
	for _, model := range ts.models {
		truncateStatement := &gorm.Statement{DB: ts.DB}
		err := truncateStatement.Parse(model)

		if err != nil {
			panic(err)
		}
		tableName := truncateStatement.Schema.Table
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)
		err = ts.DB.Exec(query).Error
		if err != nil {
			panic(err)
		}
	}
}
