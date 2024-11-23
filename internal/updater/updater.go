package updater

import (
	"log"
	"sync"
	"time"

	"github.com/SamoylikV/LocaleParse/internal/config"
	"github.com/SamoylikV/LocaleParse/internal/google"
	"github.com/SamoylikV/LocaleParse/internal/redis"
)

type Updater struct {
	redisClient *redis.Client
	cfg         *config.Config
}

func NewUpdater(redisClient *redis.Client, cfg *config.Config) *Updater {
	return &Updater{
		redisClient: redisClient,
		cfg:         cfg,
	}
}

func (u *Updater) UpdateLocales() error {
	var wg sync.WaitGroup
	var dataRu, dataEng map[string]string
	var errRu, errEng error

	wg.Add(2)
	go func() {
		defer wg.Done()
		dataRu, errRu = google.Parse(u.cfg, u.cfg.RuReadRange)
	}()
	go func() {
		defer wg.Done()
		dataEng, errEng = google.Parse(u.cfg, u.cfg.EngReadRange)
	}()
	wg.Wait()

	if errRu != nil {
		log.Printf("Error loading Russian data: %v", errRu)
		return errRu
	}
	if errEng != nil {
		log.Printf("Error loading English data: %v", errEng)
		return errEng
	}

	locales := map[string]map[string]string{
		"ru":  dataRu,
		"eng": dataEng,
	}

	err := u.redisClient.SetLocaleData("locales", locales, 24*time.Hour)
	if err != nil {
		log.Printf("Failed to save locales to Redis: %v", err)
		return err
	}

	log.Println("Locales successfully updated in Redis")
	return nil
}

func (u *Updater) StartAutoUpdate(interval time.Duration) {
	go func() {
		for {
			log.Println("Starting automatic locales update")
			if err := u.UpdateLocales(); err != nil {
				log.Printf("Error during automatic update: %v", err)
			}
			time.Sleep(interval)
		}
	}()
}
