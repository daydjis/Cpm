package main

import (
	"awesomeProject3/internal/bot"
	"awesomeProject3/internal/config"
	"awesomeProject3/internal/db"
	"awesomeProject3/internal/filter"
	products "awesomeProject3/internal/parser"
	"awesomeProject3/internal/repository"
	"awesomeProject3/internal/scheduler"
	"awesomeProject3/internal/service"
	"awesomeProject3/models"
	"log"
)

func main() {
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

	go bot.HandleUpdates()

	job := func() {
		subs, err := svc.Subscription.ListForScheduler()
		if err != nil {
			log.Println("scheduler: list subscriptions:", err)
			return
		}
		baseURL := config.ScrapeBaseURL
		if baseURL == "" {
			log.Println("scheduler: SCRAPE_BASE_URL is not set")
			return
		}
		// Уникальные URL для парсинга (один URL — один раз)
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
	go sched.Start()

	log.Println("App running")
	select {}
}
