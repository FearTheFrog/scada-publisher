package clientmetadata

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

const configFileName = "config.json"

type ClientDetails struct {
	ID               uuid.UUID `json:"ID"`
	ClientName       string    `json:"client_name"`
	OrganizationName string    `json:"org_name"`
	ContactEmail     string    `json:"contact_email"`
	CSVFilePath      string    `json:"csv_file_path"`
}

func Run() {
	client := GetClientDetails()
	if client.CSVFilePath == "" {
		client.CSVFilePath = getCSVFilePath()
		saveConfig(&client)
	}
	displayClientDetails(client)
}

func GetClientDetails() ClientDetails {
	// Try to load existing config first
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load the config: %s", err)
	}

	// If client details are already in config, return them
	if config.ID != uuid.Nil {
		return *config
	}

	// Otherwise, gather client details
	fmt.Println("######################################")
	fmt.Println("Collecting Client Metadata...")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter client name:")
	clientName, _ := reader.ReadString('\n')

	fmt.Println("Enter organization name (leave blank if same as client name):")
	organizationName, _ := reader.ReadString('\n')
	organizationName = strings.TrimSpace(organizationName)
	if organizationName == "" {
		organizationName = strings.TrimSpace(clientName)
	}

	fmt.Println("Enter contact email:")
	contactEmail, _ := reader.ReadString('\n')

	clientDetails := ClientDetails{
		ID:               uuid.New(),
		ClientName:       strings.TrimSpace(clientName),
		OrganizationName: organizationName,
		ContactEmail:     strings.TrimSpace(contactEmail),
	}

	return clientDetails
}

func getCSVFilePath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the filepath to your CSV file:")

	path, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read the filepath: %s", err)
	}

	return strings.TrimSpace(path)
}

// loadConfig attempts to load the configuration from the config file.
// If the file doesn't exist, it returns an empty ClientDetails.
func loadConfig() (*ClientDetails, error) {
	fileData, err := os.ReadFile(configFileName)
	if os.IsNotExist(err) {
		return &ClientDetails{}, nil
	} else if err != nil {
		return nil, err
	}

	var clientDetails ClientDetails
	if err := json.Unmarshal(fileData, &clientDetails); err != nil {
		return nil, err
	}

	return &clientDetails, nil
}

// saveConfig writes the given client details to the config file in JSON format.
func saveConfig(details *ClientDetails) error {
	fileData, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFileName, fileData, 0644)
}

func displayClientDetails(client ClientDetails) {
	fmt.Println("######################################")
	fmt.Println("Captured Details:")
	fmt.Printf("UUID: %s\n", client.ID)
	fmt.Printf("Client Name: %s\n", client.ClientName)
	fmt.Printf("Organization Name: %s\n", client.OrganizationName)
	fmt.Printf("Contact Email: %s\n", client.ContactEmail)
	fmt.Printf("CSV File Path: %s\n", client.CSVFilePath)
}
