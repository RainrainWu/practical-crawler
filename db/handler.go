package db

import (
	"fmt"
	"log"
	"os"

	// postgres database driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
)

// Record decribe the schema of record table
type Record struct {
	ID      int    `db:"id" json:"_id"`
	URLHash string `db:"url_hash" json:"url_hash"`
}

const (
	dropSQL string = `
	DROP TABLE IF EXISTS record;
	`

	createSQL string = `
	CREATE TABLE IF NOT EXISTS record (
		id INT GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
		url_hash VARCHAR(64) NOT NULL
	);
	`
)

// Handler is the export interface of dbHandler
type Handler interface {
	init()
	Push(urlHash string)
	Search(urlHash string) bool
}

type handler struct {
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
	log.Println("Instantiate dbHandler instance")
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

func (db *handler) Push(urlHash string) {

	pushSQL := `
	INSERT INTO record (url_hash)
	SELECT $1::VARCHAR(64)
	WHERE NOT EXISTS (
		SELECT id from record
		WHERE record.url_hash = $1
	)
	`
	tx := db.database.MustBegin()
	tx.MustExec(pushSQL, urlHash)
	tx.Commit()
}

func (db *handler) Search(urlHash string) bool {

	searchSQL := `
	SELECT id from record
	WHERE record.url_hash = $1
	`
	records := []Record{}
	err := db.database.Select(&records, searchSQL, urlHash)
	if err != nil {
		log.Fatal(err)
	}
	return len(records) != 0
}
