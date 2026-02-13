package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		direction = flag.String("direction", "up", "Migration direction: up, down, force, version")
		steps     = flag.Int("steps", 0, "Number of steps (0 = all)")
		force     = flag.Int("force", -1, "Force set version (use with -direction=force)")
	)
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "manpasik")
		pass := getEnv("DB_PASSWORD", "manpasik_dev_password")
		name := getEnv("DB_NAME", "manpasik")
		sslmode := getEnv("DB_SSLMODE", "disable")
		dbURL = fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, sslmode)
	}

	migrationsPath := getEnv("MIGRATIONS_PATH", "file://migrations")

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-(*steps))
		} else {
			err = m.Down()
		}
	case "force":
		if *force < 0 {
			log.Fatal("force requires -force=<version>")
		}
		err = m.Force(*force)
	case "version":
		v, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("Get version failed: %v", verr)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", v, dirty)
		return
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Migration %s completed successfully", *direction)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
