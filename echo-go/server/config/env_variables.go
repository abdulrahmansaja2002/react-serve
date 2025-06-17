package config

import (
	"sync"

	"echo-react-serve/helpers/config"

	"log"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Name string `env:"APP_NAME"`

		// CHANGE LATER IF YOU HAVE ACCESS TO .env FILE
		Environment string `env-default:"development"`
		// Environtment            string `env:"APP_ENV" env-default:"production"`

		BaseURL                 string `env:"APP_BASE_URL" env-default:"http://localhost:3000"`
		FrontendBaseURL         string `env:"APP_FRONTEND_BASE_URL" env-default:"http://localhost:3000"`
		Port                    string `env:"APP_PORT" env-default:"5000"`
		LogLevel                string `env:"APP_LOG_LEVEL" env-default:"debug"`
		LogFile                 string `env:"APP_LOG_FILE" env-default:"./logs/app.log"`
		LocalStoragePublicPath  string `env:"LOCAL_STORAGE_PUBLIC_PATH" env-default:"./storage/public"`
		LocalStoragePrivatePath string `env:"LOCAL_STORAGE_PRIVATE_PATH" env-default:"./storage/private"`
		EmailAppPassword        string `env:"EMAIL_APP_PASSWORD"`
		EmailServiceAddress     string `env:"EMAIL_SERVICE_ADDRESS"`
		TimeoutDuration         int    `env:"TIMEOUT_DURATION" env-default:"30"`
		TZ                      string `env:"TZ" env-default:"Asia/Jakarta"`
		CookieSecure            bool
		CookieSameSite          string
		FileServerUrl           string `env:"FILE_SERVER_URL" env-default:"http://localhost:3233/api/scholarship_form/file/?paths="`
		ProtectAPI              bool   `env:"APP_PROTECT_API" env-default:"true"`
		RealBackendUrl          string `env:"APP_REAL_BACKEND_URL" env-default:"http://localhost:3000"`
		CookieSecret            string `env:"APP_COOKIE_SECRET" env-default:"cookiesecret"`
		ClientBuildsFolder      string `env:"APP_CLIENT_BUILDS_FOLDER" env-default:"dist"`
	}
	DB struct {
		ConnectionTimeout int    `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int    `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int    `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int    `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
		CmsSchemaName     string `env:"DB_CMS_SCHEMA_NAME" env-default:"cms"`
	}
	Guard struct {
		JwtPrivateKey             string `env:"JWT_PRIVATE_KEY"`
		JwtAccessTokenExpiration  int    `env:"JWT_ACCESS_TOKEN_EXPIRATION" env-default:"24"`   // in hours
		JwtRefreshTokenExpiration int    `env:"JWT_REFRESH_TOKEN_EXPIRATION" env-default:"120"` // in hours
	}
	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432"`
		Username string `env:"POSTGRES_USER" env-default:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"POSTGRES_DB" env-default:"humas_test"`
		SslMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
	}
	Mongo struct {
		Host     string `env:"MONGO_HOST" env-default:"localhost"`
		Port     string `env:"MONGO_PORT" env-default:"27017"`
		Username string `env:"MONGO_USER" env-default:""`
		Password string `env:"MONGO_PASSWORD" env-default:""`
		Database string `env:"MONGO_DB" env-default:"humas_test"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST" env-default:"localhost"`
		Port     int    `env:"REDIS_PORT" env-default:"6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
		Database int    `env:"REDIS_DATABASE" env-default:"0"`
	}
	Storage struct {
		AccessKey string `env:"STORAGE_ACCESSKEY"`
		SecretKey string `env:"STORAGE_SECRETKEY"`
		Endpoint  string `env:"STORAGE_ENDPOINT"`
		SslMode   bool   `env:"STORAGE_SSL_MODE" env-default:"false"`
	}
	SMTP struct {
		Host     string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
		Port     int    `env:"SMTP_PORT" env-default:"587"`
		Email    string `env:"SMTP_EMAIL"`
		Password string `env:"SMTP_PASSWORD"`
	}
	Elasticsearch struct {
		Host       string `env:"ELASTICHOST" env-default:"localhost"`
		Port       int    `env:"ELASTICPORT" env-default:"9200"`
		MaxRetries int    `env:"ES_MAX_RETRIES" env-default:"30"`
		RetryDelay string `env:"ES_RETRY_DELAY" env-default:"10s"`
	}
	Security struct {
		GetAllKey string `env:"SECURITY_GET_ALL_KEY" env-default:"getallkey" env-description:"key for get all data"`
		AdminKey  string `env:"SECURITY_ADMIN_KEY" env-default:"adminkey" env-description:"key for admin access"`
	}
}

type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	log.Println("initialize config...")
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal("get config error:", err)
		}

		if Envs.App.Environment == "production" {
			Envs.App.CookieSecure = true
			Envs.App.CookieSameSite = "Lax"
		} else {
			Envs.App.CookieSecure = false
			Envs.App.CookieSameSite = "Lax"
		}
	})
	log.Println("config initialized.")
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
