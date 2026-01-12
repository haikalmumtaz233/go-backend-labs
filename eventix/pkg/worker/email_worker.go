package worker

import (
	"log"
	"time"

	"eventix/internal/entity"
)

type EmailJob struct {
	Order entity.Order
	Email string
}

func StartEmailWorker(jobChan <-chan EmailJob) {
	go func() {
		log.Println("[EmailWorker] Started and listening for jobs...")
		for job := range jobChan {
			processEmailJob(job)
		}
		log.Println("[EmailWorker] Channel closed, worker stopped")
	}()
}

func processEmailJob(job EmailJob) {
	log.Printf("[EmailWorker] Simulating email sending for Order ID %d to %s...", job.Order.ID, job.Email)
	log.Printf("[EmailWorker] Order Details - Event: %s, Quantity: %d, Total: $%.2f",
		job.Order.Event.Title, job.Order.Quantity, job.Order.TotalAmount)

	time.Sleep(2 * time.Second)

	log.Printf("[EmailWorker] Email sent successfully for Order ID %d", job.Order.ID)
}
