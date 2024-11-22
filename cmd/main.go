package main

import (
	"fmt"
	"github.com/SamoylikV/LocaleParse/internal/config"
	"github.com/SamoylikV/LocaleParse/internal/google"
	"log"
	"sync"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	var dataRu map[string]string
	var dataEng map[string]string
	var errRu error
	var errEng error
	wg.Add(2)
	go func() {
		defer wg.Done()
		dataRu, errRu = google.Parse(cfg, cfg.RuReadRange)
	}()
	go func() {
		defer wg.Done()
		dataEng, errEng = google.Parse(cfg, cfg.EngReadRange)
	}()
	wg.Wait()
	if errRu != nil {
		log.Fatal(errRu)
	}
	if errEng != nil {
		log.Fatal(errEng)
	}
	fmt.Println("Russian Data:", dataRu)
	fmt.Println("English Data:", dataEng)
}
