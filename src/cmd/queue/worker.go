package main

import (
	"flag"
	"fmt"
	"logispro/internal/worker"
	"os"
)

func main() {
	workerType := flag.String("worker", "", "")
	flag.Parse()
	switch *workerType {
	case "building_embd":
		worker.StartBuildingEmbdGenerationWorker()
	case "contact_embd":
		worker.StartContactEmbdGeneration()
	case "email_invoice":
		worker.StartEmailPaymentInvoicesWorker()
	default:
		fmt.Println("Usage: ./worker --worker={type}")
		os.Exit(1)
	}
}
