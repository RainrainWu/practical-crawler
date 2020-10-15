package scraper

import (
	"log"

	"github.com/gocolly/colly"
)

// Worker is the export interface of worker
type Worker interface {
	Visit()
}

// worker describe the members of scrpaing worker
type worker struct {
	collector *colly.Collector
	url       string
}

// Option is the abstract configure option
type Option interface {
	apply(*worker)
}

type optionFunc func(*worker)

func (f optionFunc) apply(c *worker) {

	f(c)
}

// CollectorOption is a setter of collector member
func CollectorOption(c *colly.Collector) Option {
	return optionFunc(func(w *worker) {
		w.collector = c
	})
}

// URLOption is a setter of url member
func URLOption(url string) Option {
	return optionFunc(func(w *worker) {
		w.url = url
	})
}

// NewWorker instantiate a new worker
func NewWorker(opts ...Option) Worker {

	instance := &worker{}
	log.Println("Instantiate worker instance")
	for _, opt := range opts {
		opt.apply(instance)
	}
	if instance.collector == nil {
		instance.collector = colly.NewCollector()
	}
	instance.hook()
	return instance
}

func (w *worker) hook() {

	w.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		log.Println(e.Attr("href"))
	})
	w.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)
	})
}

func (w *worker) Visit() {
	log.Println("Visiting ", w.url)
}
