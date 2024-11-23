package main

import (
	"github.com/SamoylikV/LocaleParse/internal/config"
	"github.com/SamoylikV/LocaleParse/internal/redis"
	"github.com/SamoylikV/LocaleParse/internal/updater"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	localeUpdater := updater.NewUpdater(redisClient, cfg)

	localeUpdater.StartAutoUpdate(24 * time.Hour)

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received manual update request")
		if err := localeUpdater.UpdateLocales(); err != nil {
			log.Printf("Manual update failed: %v", err)
			http.Error(w, "Failed to update locales", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Locales updated successfully"))
	})

	log.Println("Starting server on :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
