package main

import (
	"practical-crawler/scraper"
)

func main() {

	worker := scraper.NewWorker(
		scraper.URLOption("http://go-colly.org/"),
	)
	worker.Visit()
}
