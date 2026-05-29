package push

import (
	"context"
	"log"
	"os"
	"time"
)

func StartCron() {
	if os.Getenv("PUSH_CRON_ENABLED") != "true" {
		log.Println("push: cron disabled (set PUSH_CRON_ENABLED=true)")
		return
	}

	debugFCMCredentials()
	logFCMConfigStatus()

	go func() {
		loc := koreaLocation()
		for {
			now := time.Now().In(loc)
			next := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, loc)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}
			wait := time.Until(next)
			log.Printf("push: next re-engagement run in %s (10:00 KST)", wait.Round(time.Second))
			time.Sleep(wait)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			RunReengagement(ctx, time.Now())
			cancel()
		}
	}()
}
