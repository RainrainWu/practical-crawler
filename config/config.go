package config

import (
	"log"

	"github.com/joho/godotenv"
)

const (

	// URLPattern indicates the regular expression of valid URL
	URLPattern string = "^http[s]?://[a-z0-9-]+(.[a-z0-9-]+)+"

	// URLLevelUpperbound indicates the maximum depth of route.
	URLLevelUpperbound int = 4

	// DBHandlerMaxConn indicates the maximum connections to postgres db
	DBHandlerMaxConn int = 100

	// BrokerJobsMaximum indicates the maximums of queueing jobs
	BrokerJobsMaximum int = 4096

	// BrokerCacheSize indicates the maximums size of broker lru cache
	BrokerCacheSize int = 2048

	// WorkerTimeout indicate the time limit of a request
	WorkerTimeout int = 2

	// WorkerAmount indicates the size of workers pool
	WorkerAmount int = 256

	// BenchmarkInterval indicates the time interval for benchmark
	BenchmarkInterval int = 5

	// BenchmarkDuration indicates the time duration for benchmark
	BenchmarkDuration int = 60
)

var (

	// URLSeed indicates the urls start to scrape
	URLSeed []string = []string{

		// Search-engine
		"https://tw.yahoo.com/",
		"https://www.yam.com/",

		// E-commerce
		"https://shopee.tw/",
		"https://24h.pchome.com.tw/",
		"https://www.momoshop.com.tw/main/Main.jsp",
		// "https://www.ruten.com.tw/",
		// "https://www.books.com.tw/",
		"https://www.rakuten.com.tw/",
		"https://www.buy123.com.tw/",
		"https://www.pcone.com.tw/",
		"https://www.etmall.com.tw/",
		"https://tw.carousell.com/",
		"https://www.ebay.com/",
		"https://www.amazon.com/",

		// Media
		"https://udn.com/news/index",
		"https://www.ltn.com.tw/",
		"https://www.chinatimes.com/?chdtv",
		// "https://news.ebc.net.tw/",
		// "https://www.setn.com/",
		"https://www.nownews.com/",

		// Forum
		"https://www.pixnet.net/",
		"https://www.gamer.com.tw/",
		"https://www.mobile01.com/",
		"https://www.dcard.tw/f",
	}

	// URLExcludePostfix sort out all postfix to be discarded
	URLExcludePostfix []string = []string{"jpg", "png", "pdf", "asp"}

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
