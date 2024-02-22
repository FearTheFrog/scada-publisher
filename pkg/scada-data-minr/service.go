package scadadataminr

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	shared "github.com/eden-advisory/mcf-publisher/v2/pkg/shared"

	clientmetadata "github.com/eden-advisory/mcf-publisher/v2/pkg/client-metadata"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
)

var clientDetails clientmetadata.ClientDetails

func Run(clientDetails clientmetadata.ClientDetails) {
	csvFilePath := loadCSVFilePath(clientDetails)
	scheduleCSVRead(csvFilePath)

}

func loadCSVFilePath(clientDetails clientmetadata.ClientDetails) string {
	if clientDetails.CSVFilePath == "" {
		log.Fatalf("CSV file path is empty in client details.")
	}
	return clientDetails.CSVFilePath
}

func readCSV(csvFilePath string) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to open the CSV file: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read the CSV file: %s", err)
	}

	// Connect to NATS server
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()
	subj := shared.TEXAS_LNG
	fmt.Println(subj)
	// Just an example: print the CSV contents
	for _, record := range records {
		// Convert map to JSON
		jsonMsg, err := json.Marshal(record)
		if err != nil {
			log.Fatalf("failed to marshal JSON: %w", err)
		}
		nc.Publish(subj, jsonMsg)
		fmt.Println(string(jsonMsg))
	}
}

func scheduleCSVRead(csvFilePath string) {
	c := cron.New(cron.WithSeconds())
	fmt.Println(">Scheduled read every 10s..")
	c.AddFunc("@every 10s", func() { readCSV(csvFilePath) })

	c.Start()

	// Keep the program running
	select {}
}

// Connect to NATS server and publish a JSON message
func publishToNATS(msg []string) error {
	// Connect to NATS server
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	// Convert map to JSON
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Publish message
	if err := nc.Publish(shared.TEXAS_LNG+clientDetails.ID.String(), jsonMsg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Optional: Flush connection (waits for up to 1 second for ACK from NATS)
	// If you don't care about this, you can remove the Flush line.
	nc.Flush()

	return nil
}
