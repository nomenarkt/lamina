// Package tasks contains background workers like user cleanup tasks.
package tasks

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nomenarkt/lamina/internal/user"
)

// StartUserCleanupTask starts a background job to remove unconfirmed users at a custom interval.
func StartUserCleanupTask(repo user.Repo) {
	intervalStr := os.Getenv("CLEANUP_INTERVAL")
	if intervalStr == "" {
		intervalStr = "1h" // default fallback
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("❌ Invalid CLEANUP_INTERVAL=%q: %v. Using default 1h", intervalStr, err)
		interval = 1 * time.Hour
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := repo.DeleteExpiredPendingUsers(context.Background()); err != nil {
				log.Printf("❌ Failed to clean up expired users: %v", err)
			} else {
				log.Println("🧹 Expired pending users cleaned up")
			}
		}
	}()
}
