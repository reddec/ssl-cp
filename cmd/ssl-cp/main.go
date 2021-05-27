package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/jessevdk/go-flags"
	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/api/server"
	"github.com/reddec/ssl-cp/db"
	"github.com/reddec/ssl-cp/service"
	"github.com/reddec/ssl-cp/ui"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const version = "dev"

type Config struct {
	Bind            string        `long:"bind" env:"BIND" description:"Binding address" default:"127.0.0.1:8000"`
	TLS             bool          `long:"tls" env:"TLS" description:"Enable TLS"`
	TLSCert         string        `long:"tls-cert" env:"TLS_CERT" description:"Server TLS certificate" default:"cert.pem"`
	TLSKey          string        `long:"tls-key" env:"TLS_KEY" description:"Server TLS private key" default:"key.pem"`
	GracefulTimeout time.Duration `long:"graceful-timeout" env:"GRACEFUL_TIMEOUT" description:"Timeout for graceful shutdown" default:"10s"`
	CORS            bool          `long:"cors" env:"CORS" description:"Enable Allow-All CORS policy"`
	EncryptionKey   string        `long:"encryption-key" env:"ENCRYPTION_KEY" description:"Key used to encrypt private keys in database" default:""`
	KeySize         int           `long:"key-size" env:"KEY_SIZE" description:"Certificate key size in bits" default:"2048"`
	ExposeCA        bool          `long:"expose-ca" env:"EXPOSE_CA" description:"Allow query private key of CA"`
	Import          string        `long:"import" env:"IMPORT" description:"Import directory for initial setup"`
	DB              DBConfig      `group:"Database config" namespace:"db" env-namespace:"DB"`
}

type DBConfig struct {
	Provider string `long:"provider" env:"PROVIDER" description:"Database provider" default:"sqlite" choice:"postgres" choice:"sqlite"`
	DSN      string `long:"dsn" env:"DSN" description:"Database DSN" default:"sslcp.db?_foreign_keys=on"`
}

func main() {
	var config Config
	parser := flags.NewParser(&config, flags.Default)
	parser.LongDescription = "ssl-cp\nCertificate manager\nAuthor: Baryshnikov Aleksandr <dev@baryshnikov.net>\nVersion: " + version
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if config.EncryptionKey == "" {
		log.Println("WARNING! ENCRYPTION KEY IS NOT CHANGED AND BLANK - consider to use --encryption-key flag or ENCRYPTION_KEY environment")
	}
	if config.CORS {
		log.Println("WARNING! CORS enabled - several browser-based attacks possible")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err = run(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg Config) error {
	database, err := cfg.DB.Database()
	if err != nil {
		return fmt.Errorf("initialize db: %w", err)
	}

	var opts []service.Option
	if cfg.ExposeCA {
		opts = append(opts, service.ExposeCAKey(true))
	}
	opts = append(opts, service.Key(cfg.KeySize))

	apiImpl, err := service.New(database, cfg.EncryptionKey, opts...)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	err = service.ImportFromDir(ctx, cfg.Import, apiImpl)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	var apiHandler = server.New(apiImpl)

	router := http.NewServeMux()
	router.Handle(api.Prefix+"/", http.StripPrefix(api.Prefix, apiHandler))
	router.Handle("/ui/", http.StripPrefix("/ui", ui.Handler()))
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "ui/"+ui.Path, http.StatusTemporaryRedirect)
	})

	var handler http.Handler = router
	if cfg.CORS {
		handler = cors.AllowAll().Handler(handler)
	}

	srv := &http.Server{
		Addr:    cfg.Bind,
		Handler: gziphandler.GzipHandler(handler),
	}

	go func() {
		<-ctx.Done()
		tctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
		defer cancel()
		err := srv.Shutdown(tctx)
		if err != nil {
			log.Println("failed shutdown server gracefully:", err)
		} else {
			log.Println("server stopped")
		}
	}()
	log.Println("started on", cfg.Bind)
	if cfg.TLS {
		return srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
	}
	return srv.ListenAndServe()
}

var unknownDB = errors.New("unknown database provider")

func (cfg DBConfig) Database() (*gorm.DB, error) {
	switch cfg.Provider {
	case "sqlite":
		return db.New(sqlite.Open(cfg.DSN))
	case "postgres":
		return db.New(postgres.Open(cfg.DSN))
	default:
		return nil, unknownDB
	}
}
