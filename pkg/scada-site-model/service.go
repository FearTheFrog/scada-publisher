package scadasitemodel

import (
	"fmt"

	"github.com/google/uuid"
)

func Run() {
	fmt.Println("Running Oil Scada Site Data Model Service...")
}

type FlowMeterRecord struct {
	clientID    uuid.UUID `json:"clientID"`
	Time        string    `json:"time"`
	StartTime   string    `json:"start-time"`
	EndTime     string    `json:"end-time"`
	Density     float64   `json:"density"`
	Temperature float64   `json:"temperature"`
	NetBarrels  float64   `json:"net-barrels"`
}
