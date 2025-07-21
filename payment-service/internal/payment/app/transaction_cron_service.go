package app

import (
	"fmt"
	"time"
)

// Cron job: cetak/bersihkan payment yang lebih dari 5 menit (simulasi)
func StartPaymentCleanupJob() {
	ticker := time.NewTicker(1 * time.Minute) // tiap 1 menit
	go func() {
		for range ticker.C {
			err := cleanOldPayments()
			if err != nil {
				fmt.Println("Cleanup error:", err)
			} else {
				fmt.Println("Cleanup done at", time.Now())
			}
		}
	}()
}

// Fungsi dummy: bisa diganti logic delete Mongo
func cleanOldPayments() error {
	// Simulasi: hapus data lebih dari 5 menit (dummy log)
	fmt.Println("Simulating cleanup: deleting payments older than 5 minutes...")
	// Implementasi real: tambahkan query delete by CreatedAt < now - 5 minutes
	return nil
}
