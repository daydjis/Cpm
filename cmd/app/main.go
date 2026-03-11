package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"awesomeProject3/internal/bot"
	"awesomeProject3/internal/config"
	"awesomeProject3/internal/db"
	"awesomeProject3/internal/filter"
	products "awesomeProject3/internal/parser"
	"awesomeProject3/internal/repository"
	"awesomeProject3/internal/scheduler"
	"awesomeProject3/internal/service"
	"awesomeProject3/models"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Printf("application exited with error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	config.Load()

	db.Init()
	bot.Init()

	userRepo := repository.NewUserRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	subRepo := repository.NewSubscriptionRepository(db.DB)
	notifRepo := repository.NewNotificationRepository(db.DB)

	svc := service.NewServicesFromRepos(userRepo, categoryRepo, subRepo, notifRepo)
	bot.SetServices(svc)
	bot.SetHandlersBot()

	var wg sync.WaitGroup

	startWorker := func(name string, fn func(context.Context)) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					log.Printf("worker %s panic: %v", name, r)
				}
			}()

			fn(ctx)
		}()
	}

	startWorker("telegram_bot", func(ctx context.Context) {

		bot.HandleUpdates()
	})

	job := func() {
		subs, err := svc.Subscription.ListForScheduler()
		if err != nil {
			log.Printf("scheduler: list subscriptions: %v", err)
			return
		}

		baseURL := config.ScrapeBaseURL
		if baseURL == "" {
			log.Println("scheduler: SCRAPE_BASE_URL is not set")
			return
		}

		seen := make(map[string]struct{})

		for _, s := range subs {
			p := models.FilterParams{
				SearchText: s.SearchText,
				PriceMin:   s.PriceMin,
				PriceMax:   s.PriceMax,
				RegionID:   s.RegionID,
				ProTypes:   filter.ParseProTypes(s.ProTypes),
			}
			urlStr := filter.BuildCategoryURL(baseURL, s.CategorySlug, p)

			if _, ok := seen[urlStr]; ok {
				continue
			}
			seen[urlStr] = struct{}{}

			products.FetchAndSave(urlStr, s.CategoryID, s.FilterSignature)
		}

		bot.SendNewNotifications()
	}

	sched := scheduler.New(job)

	startWorker("scheduler", func(ctx context.Context) {
		go sched.Start()
		<-ctx.Done()
		sched.Stop()
	})

	log.Println("application started")

	<-ctx.Done()
	log.Println("shutting down, stopping bot and waiting workers...")

	bot.Stop()

	wg.Wait()

	return nil
}
