package db

import (
	"fmt"
	"log"
	"os"
	"practical-crawler/config"

	// postgres database driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Record decribe the schema of record table
type Record struct {
	ID      int    `db:"id" json:"_id"`
	URLHash string `db:"url_hash" json:"url_hash"`
}

const (
	dropSQL string = `
	DROP INDEX IF EXISTS url_idx;
	DROP TABLE IF EXISTS record;
	`

	createSQL string = `
	CREATE TABLE IF NOT EXISTS record (
		id INT GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
		url_hash VARCHAR(64) UNIQUE
	);
	CREATE UNIQUE INDEX url_idx ON record (url_hash);
	`
)

// Handler is the export interface of dbHandler
type Handler interface {
	init()
	Count() int
	Push(urlHash string) error
	Search(urlHash string) bool
}

type handler struct {
	logger   *zap.SugaredLogger
	database *sqlx.DB
	drop     bool
}

// Option is the abstract configure option
type Option interface {
	apply(*handler)
}

type optionFunc func(*handler)

func (f optionFunc) apply(db *handler) {

	f(db)
}

// LoggerOption is a setter of logger member
func LoggerOption(l *zap.SugaredLogger) Option {
	return optionFunc(func(db *handler) {
		db.logger = l
	})
}

// DropOption is a setter of drop member
func DropOption(drop bool) Option {
	return optionFunc(func(db *handler) {
		db.drop = drop
	})
}

// DatabseOption is a setter of database member
func DatabseOption(database *sqlx.DB) Option {
	return optionFunc(func(db *handler) {
		db.database = database
	})
}

// NewHandler instantiate a new Handler
func NewHandler(opts ...Option) Handler {

	instance := &handler{}
	for _, opt := range opts {
		opt.apply(instance)
	}
	if instance.database == nil {
		connStr := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"
		connStr = fmt.Sprintf(
			connStr,
			os.Getenv("PG_HOST"),
			os.Getenv("PG_PORT"),
			os.Getenv("PG_USERNAME"),
			os.Getenv("PG_DBNAME"),
			os.Getenv("PG_PASSWORD"),
			os.Getenv("PG_SSLMODE"),
		)
		conn, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		conn.SetMaxOpenConns(config.DBHandlerMaxConn)
		instance.database = conn
	}
	instance.init()
	return instance
}

func (db *handler) init() {

	if db.drop {
		db.database.MustExec(dropSQL)
	}
	db.database.MustExec(createSQL)
}

// Count return the amount of found urls
func (db *handler) Count() int {

	searchSQL := `
	SELECT id from record
	`
	records := []Record{}
	err := db.database.Select(&records, searchSQL)
	if err != nil {
		db.logger.Warnf("[DB Handler] %s", err)
	}
	return len(records)
}

// Push will insert a url into record table
func (db *handler) Push(urlHash string) error {

	pushSQL := `
	INSERT INTO record (url_hash)
	SELECT $1::VARCHAR(64)
	WHERE NOT EXISTS (
		SELECT id from record
		WHERE record.url_hash = $1
	)
	`
	tx := db.database.MustBegin()
	_, err := tx.Exec(pushSQL, urlHash)
	if err != nil {
		db.logger.Warnf("[DB Handler] %s", err)
	} else {
		tx.Commit()
	}
	return err
}

// Search will check if a urlhash is already exists
func (db *handler) Search(urlHash string) bool {

	searchSQL := `
	SELECT id from record
	WHERE record.url_hash = $1
	`
	records := []Record{}
	err := db.database.Select(&records, searchSQL, urlHash)
	if err != nil {
		db.logger.Warnf("[DB Handler] %s", err)
	}
	return len(records) != 0
}
