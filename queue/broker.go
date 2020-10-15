package queue

import "log"

// Broker is the export interface of Broker
type Broker interface {
	Push(url string)
	Pop() string
}

// broker describe the members of jobs broker
type broker struct {
	jobs chan string
}

// Option is the abstract configure option
type Option interface {
	apply(*broker)
}

type optionFunc func(*broker)

func (f optionFunc) apply(c *broker) {

	f(c)
}

// JobsOption is a setter of jobs member
func JobsOption(c chan string) Option {
	return optionFunc(func(b *broker) {
		b.jobs = c
	})
}

// NewBroker instantiate a new broker
func NewBroker(opts ...Option) Broker {

	instance := &broker{}
	log.Println("Instantiate broker instance")
	for _, opt := range opts {
		opt.apply(instance)
	}
	if instance.jobs == nil {
		instance.jobs = make(chan string, 256)
	}
	return instance
}

// Push will push a new url into jobs queue
func (b *broker) Push(url string) {

	b.jobs <- url
}

// Pop will pop out a url from jobs queue
func (b *broker) Pop() string {

	return <-b.jobs
}
