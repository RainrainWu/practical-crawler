package config

import (
	"log"

	"github.com/joho/godotenv"
)

const (

	// RootURL indicates the url start to scrape
	RootURL string = "https://www.sina.com.tw/"

	// URLPattern indicates the regular expression of valid URL
	URLPattern string = "^http[s]?://[a-z0-9-]+(.[a-z0-9-]+)+"

	// BrokerJobsMaximum indicates the maximums of queueing jobs
	BrokerJobsMaximum int = 512

	// BrokerCacheSize indicates the maximums size of broker lru cache
	BrokerCacheSize int = 256

	// WorkerTimeout indicate the time limit of a request
	WorkerTimeout int = 2

	// WorkerAmount indicates the size of workers pool
	WorkerAmount int = 64
)

var (

	// URLExcludePostfix sort out all postfix to be discard
	URLExcludePostfix []string = []string{"jpg", "png", "pdf"}

	// URLDiscardParameter list all parameter tag to be truncate
	URLDiscardParameter []string = []string{"@", "#", "?"}
)

func init() {

	loadEnv()
}

func loadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
