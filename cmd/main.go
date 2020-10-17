package main

import (
	"context"
	"log"
	"os"
	"time"

	"practical-crawler/config"
	"practical-crawler/db"
	"practical-crawler/queue"
	"practical-crawler/scraper"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProvideLogger provide a Logger instance
func ProvideLogger() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.OutputPaths = []string{
		"./logs/crawler.log",
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger.Sugar()
}

// ProvideDBHandler provide a DBHandler instance
func ProvideDBHandler(l *zap.SugaredLogger) db.Handler {
	return db.NewHandler(
		db.LoggerOption(l),
		db.DropOption(true),
	)
}

// ProvideBroker provide a broker instance
func ProvideBroker(l *zap.SugaredLogger, h db.Handler) queue.Broker {
	return queue.NewBroker(
		queue.LoggerOption(l),
		queue.DBHandlerOption(h),
	)
}

// ProvideWorkers provide a set of workers instances
func ProvideWorkers(l *zap.SugaredLogger, b queue.Broker) []scraper.Worker {
	workers := []scraper.Worker{}
	for i := 0; i < config.WorkerAmount; i++ {
		w := scraper.NewWorker(
			scraper.IDOption(i),
			scraper.LoggerOption(l),
			scraper.BrokerOption(b),
		)
		workers = append(workers, w)
	}
	return workers
}

func register(
	lifecycle fx.Lifecycle,
	l *zap.SugaredLogger,
	h db.Handler,
	b queue.Broker,
	ws []scraper.Worker,
) {

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go timer(b, h)
				for _, url := range config.URLSeed {
					b.Push(url)
				}
				for _, worker := range ws {
					go worker.Run()
				}
				return nil
			},
			OnStop: func(context.Context) error {
				l.Sync()
				return nil
			},
		},
	)
}

func timer(b queue.Broker, h db.Handler) {
	select {
	case <-time.After(time.Duration(config.BenchmarkDuration) * time.Second):
		log.Println(
			"\n",
			"Benchmark", config.BenchmarkDuration, "seconds\n",
			"Left", b.GetLeft(), "jobs\n",
			"Encounter", b.GetErrorCount(), "errors\n",
			"Recieve", b.GetAccumulate(), "responses\n",
			"Total", h.Count(), "records",
		)
		os.Exit(0)
	}
}

func main() {
	fx.New(
		fx.Provide(
			ProvideLogger,
			ProvideDBHandler,
			ProvideBroker,
			ProvideWorkers,
		),
		fx.Invoke(register),
	).Run()
}
