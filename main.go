package main

import (
	"github.com/go-co-op/gocron"
	"github.com/ting-app/ting-task-voa/ting"
	"log"
	"time"
)

func main() {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).Day().At("23:59").Do(func() {
		err := ting.RunTask("https://learningenglish.voanews.com/api/zpyp_e-rm_")

		if err != nil {
			log.Printf("Run task error %v\n", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule task, %v", err)
	}

	log.Println("Task scheduled")

	scheduler.StartBlocking()
}
