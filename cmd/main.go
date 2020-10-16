package main

import (
	"context"
	"os"
	"time"

	"practical-crawler/config"
	"practical-crawler/db"
	"practical-crawler/queue"
	"practical-crawler/scraper"

	"go.uber.org/fx"
)

// ProvideDBHandler provide a DBHandler instance
func ProvideDBHandler() db.Handler {
	return db.NewHandler(
		db.DropOption(true),
	)
}

// ProvideBroker provide a broker instance
func ProvideBroker(h db.Handler) queue.Broker {
	return queue.NewBroker(
		queue.DBHandlerOption(h),
	)
}

// ProvideWorkers provide a set of workers instances
func ProvideWorkers(b queue.Broker) []scraper.Worker {
	workers := []scraper.Worker{}
	for i := 0; i < config.WorkerAmount; i++ {
		w := scraper.NewWorker(
			scraper.IDOption(i),
			scraper.BrokerOption(b),
		)
		workers = append(workers, w)
	}
	return workers
}

func register(lifecycle fx.Lifecycle, h db.Handler, b queue.Broker, ws []scraper.Worker) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				b.Push(config.RootURL)
				for _, worker := range ws {
					go worker.Run()
				}
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}

func timer() {
	select {
	case <-time.After(time.Duration(30) * time.Second):
		os.Exit(0)
	}
}

func main() {
	go timer()
	fx.New(
		fx.Provide(
			ProvideDBHandler,
			ProvideBroker,
			ProvideWorkers,
		),
		fx.Invoke(register),
	).Run()
}
