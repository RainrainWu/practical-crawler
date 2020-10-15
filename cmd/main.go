package main

import (
	"practical-crawler/queue"
	"practical-crawler/scraper"
)

func main() {

	broker := queue.NewBroker(
		queue.JobsOption(make(chan string, 32)),
	)
	broker.Push("http://go-colly.org/")

	worker := scraper.NewWorker(
		scraper.URLOption(broker.Pop()),
	)
	worker.Visit()
}
