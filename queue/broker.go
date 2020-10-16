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
)

// Broker is the export interface of Broker
type Broker interface {
	Push(url string)
	Pop() string
	GetLeft() int
	GetAccumulate() int
}

// broker describe the members of jobs broker
type broker struct {
	dbHandler  db.Handler
	pattern    *regexp.Regexp
	cache      *lru.Cache
	jobs       chan string
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
		accumulate: 0,
	}
	log.Println("Instantiate broker instance")
	for _, opt := range opts {
		opt.apply(instance)
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
	return instance
}

// Push will push a new url into jobs queue
func (b *broker) Push(url string) {

	data := []byte(url)
	digest := md5.Sum(data)
	urlHash := hex.EncodeToString(digest[:])
	if !b.pattern.MatchString(url) {
		log.Println("Invalid URL", url)
	} else if excludePostfix(url) {
		log.Println("Exclude postfix", url)
	} else if b.cache.Contains(urlHash) {
		log.Println("Duplicate URL", url, "conflict hash", urlHash)
	} else if b.dbHandler.Search(urlHash) {
		log.Println("Duplicate URL", url, "conflict hash", urlHash)
	} else {
		select {
		case b.jobs <- url:
			b.cache.Add(urlHash, true)
			b.dbHandler.Push(urlHash)
			log.Println("Pushed", url)
		default:
			log.Println("Channel full, discard", url)
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
	b.accumulate++
	return url
}

// GetLeft get the left amount of jobs
func (b *broker) GetLeft() int {
	return len(b.jobs)
}

// GetAccumulate get the accumulate amount of executed jobs
func (b *broker) GetAccumulate() int {
	return b.accumulate
}
