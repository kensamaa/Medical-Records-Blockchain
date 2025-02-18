package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kensamaa/blockchain-medical-records/chaincode/smartcontract"
)

func main() {
	// Create a new chaincode instance using our SmartContract struct.
	chaincode, err := contractapi.NewChaincode(new(smartcontract.SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s\n", err.Error())
		os.Exit(1)
	}

	// Start the chaincode and listen for incoming requests.
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err.Error())
		os.Exit(1)
	}
}
