package config

const (

	// RootURL indicates the url start to scrape
	RootURL string = "https://shopping.pchome.com.tw/"

	// URLPattern indicates the regular expression of valid URL
	URLPattern string = "^http[s]?://[a-z0-9-]+(.[a-z0-9-]+)+"

	// JobsMaximum indicates the maximums of queueing jobs
	JobsMaximum int = 256

	// WorkerTimeout indicate the time limit of a request
	WorkerTimeout int = 2

	// WorkerAmount indicates the size of workers pool
	WorkerAmount int = 8
)
