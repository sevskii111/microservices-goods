package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sevskii111/microservices-goods/pkg/models/mysql"
)

type Config struct {
	Addr          string
	Secret        string
	TokenLifetime int
}

type application struct {
	cfg      *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	goods    *mysql.GoodModel
}

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.Secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	dsn := flag.String("dsn", "root:@/goods?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		cfg:      cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		goods:    &mysql.GoodModel{DB: db},
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
