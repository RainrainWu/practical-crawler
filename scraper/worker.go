package scraper

import (
	"log"
	"strings"
	"time"

	"practical-crawler/config"
	"practical-crawler/queue"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

// Worker is the export interface of worker
type Worker interface {
	Run()
	Visit()
}

// worker describe the members of scrpaing worker
type worker struct {
	id        int
	logger    *zap.SugaredLogger
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

// LoggerOption is a setter of logger member
func LoggerOption(l *zap.SugaredLogger) Option {
	return optionFunc(func(w *worker) {
		w.logger = l
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
	for _, opt := range opts {
		opt.apply(instance)
	}
	if instance.logger == nil {
		log.Fatal("missing logger")
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
	instance.logger.Info("Instantiate worker instance")
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
	w.collector.OnResponse(func(r *colly.Response) {
		w.broker.Accumulate()
	})
	w.collector.OnError(func(r *colly.Response, err error) {
		w.broker.AddError()
		w.logger.Warnf("[Worker %d] Error %s", w.id, err)
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

	for _, param := range config.URLDiscardParameter {
		raw = strings.Split(raw, param)[0]
	}
	return raw
}

func (w *worker) Run() {
	w.idle <- true
	for {
		select {
		case <-w.idle:
			go w.Visit()
		}
	}
}

func (w *worker) Visit() {

	w.url = w.broker.Pop()
	w.logger.Debugf("[Worker %d] Visiting %s", w.id, w.url)
	w.collector.Visit(w.url)
	w.idle <- true
}
