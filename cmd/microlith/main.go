package main

import (
	"fmt"

	clientmetadata "github.com/eden-advisory/mcf-publisher/v2/pkg/client-metadata"
	scadadataminr "github.com/eden-advisory/mcf-publisher/v2/pkg/scada-data-minr"
	scadasitemodel "github.com/eden-advisory/mcf-publisher/v2/pkg/scada-site-model"
)

func main() {

	fmt.Println("Initializing MBL Publisher..")
	clientmetadata.Run()
	scadasitemodel.Run()
	clientDetails := clientmetadata.GetClientDetails()

	scadadataminr.Run(clientDetails)
}
