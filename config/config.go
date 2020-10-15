package config

const (

	// RootURL indicates the url start to scrape
	RootURL string = "http://go-colly.org"

	// URLPattern indicates the regular expression of valid URL
	URLPattern string = "^http[s]?://"

	// JobsMaximum indicates the maximums of queueing jobs
	JobsMaximum int = 256

	// WorkerAmount indicates the size of workers pool
	WorkerAmount int = 8
)
