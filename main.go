package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"github.com/ting-app/ting-task-voa/ting"
	"log"
	"os"
	"time"
)

func main() {
	enableSentry := os.Getenv("ENABLE_SENTRY") == "true"
	sentryDsn := os.Getenv("SENTRY_DSN")

	if enableSentry {
		if sentryDsn == "" {
			log.Fatal("sentryDsn is required")
		}

		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		})

		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		} else {
			log.Println("sentry enabled")
		}
	}

	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).Day().At("23:59").Do(run(enableSentry))

	if err != nil {
		log.Fatalf("Failed to schedule task, %v", err)
	}

	log.Println("Task scheduled")

	scheduler.StartBlocking()
}

func run(enableSentry bool) func() {
	return func() {
		channels := []ting.Channel{
			{Url: "https://learningenglish.voanews.com/api/zpyp_e-rm_", Tag: ting.Tag{Id: 1, Name: "Arts & Culture"}},
			{Url: "https://learningenglish.voanews.com/api/ztmp_eibp_", Tag: ting.Tag{Id: 2, Name: "Education"}},
			{Url: "https://learningenglish.voanews.com/api/zmmpqeb-po", Tag: ting.Tag{Id: 3, Name: "Health & Lifestyle"}},
			{Url: "https://learningenglish.voanews.com/api/zmg_pebmyp", Tag: ting.Tag{Id: 4, Name: "Science & Technology"}},
		}

		for _, channel := range channels {
			err := ting.RunTask(channel)

			if err != nil {
				log.Printf("Run task error %v, channel=%s\n", err, channel)

				if enableSentry {
					sentry.CaptureException(err)
				}
			}
		}
	}
}
