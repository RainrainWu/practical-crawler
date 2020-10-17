package queue

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"regexp"
	"strings"

	"practical-crawler/config"
	"practical-crawler/db"

	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
)

// Broker is the export interface of Broker
type Broker interface {
	Push(url string)
	Pop() string
	AddError()
	Accumulate()
	GetLeft() int
	GetErrorCount() int
	GetAccumulate() int
}

// broker describe the members of jobs broker
type broker struct {
	logger     *zap.SugaredLogger
	dbHandler  db.Handler
	pattern    *regexp.Regexp
	cache      *lru.Cache
	jobs       chan string
	errorCount int
	accumulate int
}

// Option is the abstract configure option
type Option interface {
	apply(*broker)
}

type optionFunc func(*broker)

func (f optionFunc) apply(c *broker) {

	f(c)
}

// LoggerOption is a setter of logger member
func LoggerOption(l *zap.SugaredLogger) Option {
	return optionFunc(func(b *broker) {
		b.logger = l
	})
}

// DBHandlerOption is a setter of dbHandler member
func DBHandlerOption(h db.Handler) Option {
	return optionFunc(func(b *broker) {
		b.dbHandler = h
	})
}

// JobsOption is a setter of jobs member
func JobsOption(c chan string) Option {
	return optionFunc(func(b *broker) {
		b.jobs = c
	})
}

// NewBroker instantiate a new broker
func NewBroker(opts ...Option) Broker {

	instance := &broker{
		errorCount: 0,
		accumulate: 0,
	}
	for _, opt := range opts {
		opt.apply(instance)
	}
	if instance.logger == nil {
		log.Fatal("missing logger")
	}
	if instance.dbHandler == nil {
		log.Fatal("missing dbHandler")
	}
	if instance.pattern == nil {
		pattern, err := regexp.Compile(config.URLPattern)
		if err != nil {
			log.Fatal(err)
		}
		instance.pattern = pattern
	}
	if instance.cache == nil {
		cache, err := lru.New(config.BrokerCacheSize)
		if err != nil {
			log.Fatal(err)
		}
		instance.cache = cache
	}
	if instance.jobs == nil {
		instance.jobs = make(chan string, config.BrokerJobsMaximum)
	}
	instance.logger.Info("Instantiate broker instance")
	return instance
}

// Push will push a new url into jobs queue
func (b *broker) Push(url string) {

	data := []byte(url)
	digest := md5.Sum(data)
	urlHash := hex.EncodeToString(digest[:])
	if !b.pattern.MatchString(url) {
		// b.logger.Debugf("Invalid URL %s", url)
	} else if excludePostfix(url) {
		// b.logger.Debugf("Exclude postfix %s", url)
	} else if b.cache.Contains(urlHash) {
		// b.logger.Debugf("Duplicate URL %s", url)
	} else if b.dbHandler.Search(urlHash) {
		// b.logger.Debugf("Duplicate URL %s", url)
	} else {
		select {
		case b.jobs <- url:
			b.cache.Add(urlHash, true)
			b.dbHandler.Push(urlHash)
			b.logger.Debugf("Pushed %s", url)
		default:
			// b.logger.Debugf("Channel full, discard %s", url)
		}
	}
}

func excludePostfix(url string) bool {

	urlSplit := strings.Split(url, ".")
	target := urlSplit[len(urlSplit)-1]
	for _, postfix := range config.URLExcludePostfix {
		if postfix == target {
			return true
		}
	}
	return false
}

// Pop will pop out a url from jobs queue
func (b *broker) Pop() string {

	url := <-b.jobs
	return url
}

func (b *broker) AddError() {
	b.errorCount++
}

func (b *broker) Accumulate() {
	b.accumulate++
}

// GetLeft get the left amount of jobs
func (b *broker) GetLeft() int {
	return len(b.jobs)
}

// GetErrorCount get the accumulate amount of encountered errors
func (b *broker) GetErrorCount() int {
	return b.errorCount
}

// GetAccumulate get the accumulate amount of executed jobs
func (b *broker) GetAccumulate() int {
	return b.accumulate
}
