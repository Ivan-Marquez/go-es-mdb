package main

import (
	"log"

	"github.com/robfig/cron"
)

func main() {
	log.Println("Create new cron")
	c := cron.New()

	err := c.AddFunc("@every 1m", func() {
		// TODO: query es_failed_updates and run ES update
		// if successful, delete from es_failed_updates
		log.Println("[Job 1]Every minute job")
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
	select {}
}
