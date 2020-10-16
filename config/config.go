package config

import (
	"log"

	"github.com/joho/godotenv"
)

const (

	// URLPattern indicates the regular expression of valid URL
	URLPattern string = "^http[s]?://[a-z0-9-]+(.[a-z0-9-]+)+"

	// BrokerJobsMaximum indicates the maximums of queueing jobs
	BrokerJobsMaximum int = 1024

	// BrokerCacheSize indicates the maximums size of broker lru cache
	BrokerCacheSize int = 1024

	// WorkerTimeout indicate the time limit of a request
	WorkerTimeout int = 3

	// WorkerAmount indicates the size of workers pool
	WorkerAmount int = 128

	// BenchmarkDuration indicates the time duration for benchmark
	BenchmarkDuration int = 10
)

var (

	// URLSeed indicates the urls start to scrape
	URLSeed []string = []string{
		"https://www.openfind.com.tw/",
		"https://24h.pchome.com.tw/",
		// "https://www.sina.com.tw/",
		"https://tw.yahoo.com/",
		"https://www.yam.com/",
	}

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
