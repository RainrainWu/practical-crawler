package main

import (
	"context"
	"practical-crawler/config"
	"practical-crawler/queue"
	"practical-crawler/scraper"

	"go.uber.org/fx"
)

// ProvideBroker provide a broker instance
func ProvideBroker() queue.Broker {
	return queue.NewBroker(
		queue.JobsOption(make(chan string, config.JobsMaximum)),
	)
}

// ProvideWorkers provide a set of workers instances
func ProvideWorkers(b queue.Broker) []scraper.Worker {
	workers := []scraper.Worker{}
	for i := 0; i < config.WorkerAmount; i++ {
		w := scraper.NewWorker(
			scraper.BrokerOption(b),
		)
		workers = append(workers, w)
	}
	return workers
}

func register(lifecycle fx.Lifecycle, b queue.Broker, ws []scraper.Worker) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				b.Push(config.RootURL)
				for _, worker := range ws {
					go worker.Visit()
				}
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}

func main() {
	fx.New(
		fx.Provide(
			ProvideBroker,
			ProvideWorkers,
		),
		fx.Invoke(register),
	).Run()
}
