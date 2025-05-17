// Package tasks contains background workers like user cleanup tasks.
package tasks

import (
	"context"
	"log"
	"time"

	"github.com/nomenarkt/lamina/internal/user"
)

// StartUserCleanupTask starts a background job to remove unconfirmed users every hour.
func StartUserCleanupTask(repo user.Repo) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
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
