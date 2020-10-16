package scraper

import (
	"log"
	"practical-crawler/config"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"practical-crawler/queue"
)

// Worker is the export interface of worker
type Worker interface {
	Run()
	Visit()
}

// worker describe the members of scrpaing worker
type worker struct {
	id        int
	broker    queue.Broker
	collector *colly.Collector
	url       string
	idle      chan bool
}

// Option is the abstract configure option
type Option interface {
	apply(*worker)
}

type optionFunc func(*worker)

func (f optionFunc) apply(c *worker) {

	f(c)
}

// IDOption is a setter of id member
func IDOption(id int) Option {
	return optionFunc(func(w *worker) {
		w.id = id
	})
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
		log.Fatal("missing broker")
	}
	if instance.collector == nil {
		instance.collector = colly.NewCollector()
	}
	if instance.idle == nil {
		instance.idle = make(chan bool, 1)
	}
	instance.hook()
	return instance
}

func (w *worker) hook() {

	w.collector.SetRequestTimeout(
		time.Duration(config.WorkerTimeout) * time.Second,
	)
	w.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		raw := string(e.Attr("href"))
		rawN := normalize(e.Request.URL.String(), raw)
		w.broker.Push(rawN)
	})
	w.collector.OnError(func(r *colly.Response, err error) {
		log.Println("Worker", w.id, "Error", err)
	})
}

func normalize(url, raw string) string {

	if len(raw) < 1 {
		return raw
	}
	if string(raw[0]) == "/" {
		raw = url + raw[1:]
	}
	if string(raw[len(raw)-1]) == "/" {
		raw = string(raw[:len(raw)-1])
	}
	raw = strings.Split(raw, "?")[0]
	raw = strings.Split(raw, "@")[0]
	raw = strings.Split(raw, "@")[0]
	return raw
}

func (w *worker) Run() {
	w.idle <- true
	for {
		select {
		case <-w.idle:
			log.Println("Worker", w.id, "Idle")
			go w.Visit()
		case <-time.After(time.Duration(config.WorkerTimeout) * time.Second):
			log.Println("Worker", w.id, "Timeout")
			go w.Visit()
		}
	}
}

func (w *worker) Visit() {

	w.url = w.broker.Pop()
	log.Println("Worker", w.id, "Visiting ", w.url)
	w.collector.Visit(w.url)
	w.idle <- true
}
