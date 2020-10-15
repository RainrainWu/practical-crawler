package scraper

import (
	"log"

	"github.com/gocolly/colly"

	"practical-crawler/queue"
)

// Worker is the export interface of worker
type Worker interface {
	Visit()
}

// worker describe the members of scrpaing worker
type worker struct {
	broker    queue.Broker
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

// BrokerOption is a setter of broker member
func BrokerOption(b queue.Broker) Option {
	return optionFunc(func(w *worker) {
		w.broker = b
	})
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
	if instance.broker == nil {
		log.Fatal("missing broker.")
	}
	if instance.collector == nil {
		instance.collector = colly.NewCollector()
	}
	instance.hook()
	return instance
}

func (w *worker) hook() {

	w.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		raw := string(e.Attr("href"))
		if len(raw) > 1 {
			if string(raw[0]) == "/" {
				raw = e.Request.URL.String() + raw
			}
			if string(raw[len(raw)-1]) == "/" {
				raw = raw[:len(raw)-1]
			}
		}
		w.broker.Push(raw)
	})
	w.collector.OnScraped(func(r *colly.Response) {
		go w.Visit()
	})
	w.collector.OnError(func(r *colly.Response, err error) {
		log.Fatal(err)
	})
}

func (w *worker) Visit() {

	w.url = w.broker.Pop()
	log.Println("Visiting ", w.url)
	w.collector.Visit(w.url)
}
