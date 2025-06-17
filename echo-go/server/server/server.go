package server

import (
	"log"

	"echo-react-serve/config"
	// "echo-react-serve/helpers/renderer"

	"echo-react-serve/server/routes"

	"github.com/labstack/echo/v4"
)

/*
 * this file is used to initialize all components
 * and create a new echo instance
 */

const (
	TEMPLATES_DIR = "./templates/*.html"
	DEBUG         = true
	ENV_PATH      = "./"
	ENV_FILENAME  = ".env"
)

func migrate() {
	log.Println("Startting auto migration...")
	// TODO: add gorm autho migration here... (if needed)
	log.Println("Auto migration completed.")
}

func loadEnv() {
	log.Printf("Initializing configuration with config: %s", ENV_PATH+ENV_FILENAME)

	config.Configuration(
		config.WithPath(ENV_PATH),
		config.WithFilename(ENV_FILENAME),
	).Initialize()
}

func init() {
	// load environment variables
	loadEnv()

	// init database
	err := config.InitMongo()
	if err != nil {
		panic(err)
	}
	migrate()
	// config.InitFirebase()

	// init ES client
	// if err := config.InitES(); err != nil {
	// 	panic(err)
	// }

	// init Storage (Minio Client)
	config.InitMinioClient()
}

// New creates a new echo instance
func New() *echo.Echo {
	e := echo.New()
	// e.Renderer = renderer.NewRenderer(TEMPLATES_DIR, DEBUG)
	routes.SetupRoutes(e)
	return e
}

// Close closes database connection
func Close() {
	config.CloseMongo()
	// config.CloseSqlServerDB()
}
