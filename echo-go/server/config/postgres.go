package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	PostgresDB *gorm.DB
	// define variable for database connection
	user, password, host, port, dbname, sslmode string
)

// InitDB initializes database connection
func InitPostgresDB() (err error) {
	// get database connection from environment variable
	user = Envs.Postgres.Username
	password = Envs.Postgres.Password
	host = Envs.Postgres.Host
	port = Envs.Postgres.Port
	dbname = Envs.Postgres.Database
	sslmode = Envs.Postgres.SslMode

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var cfg gorm.Config
	loc, err := time.LoadLocation(Envs.App.TZ)
	if err != nil {
		log.Printf("[PGSQL] Error loading location for timezone %s: %v\n", Envs.App.TZ, err)
	} else {
		log.Printf("[PGSQL] Timezone set to %s\n", Envs.App.TZ)
		cfg = gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().In(loc)
			},
		}
	}

	// open database connection
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &cfg)
	return
}

// ClosePostgresDB closes database connection
func ClosePostgresDB() {
	db, _ := PostgresDB.DB()
	db.Close()
}
