package app

import (
	"fmt"
	"time"

	"shopping-service/internal/shopping/infra"

	"github.com/robfig/cron/v3"
)

// StartTransactionCron menjalankan cron job harian
// untuk menghapus transaksi dengan status "failed" yang usianya lebih dari 24 jam
func StartTransactionCron(repo infra.TransactionRepository) {
	c := cron.New()

	// jalankan setiap hari pukul 01:00
	spec := "0 1 * * *"

	c.AddFunc(spec, func() {
		err := repo.DeleteFailedOlderThan(24 * time.Hour)
		if err != nil {
			fmt.Printf("[CRON ERROR] Gagal hapus transaksi failed: %v\n", err)
			return
		}
		fmt.Println("[CRON INFO] Transaksi failed > 24 jam berhasil dihapus")
	})

	c.Start()
}
