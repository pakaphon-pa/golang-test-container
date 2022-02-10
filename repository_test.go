package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx" // connect DB can change anything
	_ "github.com/lib/pq"     // driver db required
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	container "github.com/testcontainers/testcontainers-go" // container test
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	_repo *Repo
	_conn *sqlx.DB
)

func TestMain(m *testing.M) {
	log.Println("Prepare DB....")
	postgresPort := nat.Port("5432/tcp")
	postgres, err := container.GenericContainer(context.Background(),
		container.GenericContainerRequest{
			ContainerRequest: container.ContainerRequest{
				Image:        "postgres",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "pass",
					"POSTGRES_USER":     "user",
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true, // auto-start the container
		})

	if err != nil {
		log.Fatal("start:", err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatal("map:", err)
	}

	postgresURLTemplate := "postgres://user:pass@localhost:%s?sslmode=disable"
	postgresURL := fmt.Sprintf(postgresURLTemplate, hostPort.Port())

	log.Printf("Postgres container started, running at:  %s\n", postgresURL)
	_conn, err = sqlx.Connect("postgres", postgresURL)

	if err != nil {
		log.Fatal("connect:", err)
	}

	if err := RunMigrations(_conn); err != nil {
		log.Fatal("runMigrations:", err)
	}

	_repo = NewRepo(_conn)
	os.Exit(m.Run())
}

func TestRepoImp(t *testing.T) {
	t.Run("Create Test", func(t *testing.T) {
		user, err := _repo.CreateUser("abcd")
		require.NoError(t, err)

		getUser, err := _repo.GetAllUser()
		require.NoError(t, err)
		assert.Equal(t, user, getUser[len(getUser)-1])
	})
}
