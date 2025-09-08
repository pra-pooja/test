package main

import (
	"log"

	"supplychain/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	factoryContract := new(contracts.FactoryContract)
	logisticsContract := new(contracts.LogisticsContract)
	depotContract := new(contracts.DepotContract)
	hqContract := new(contracts.HQContract)

	// Create new chaincode
	chaincode, err := contractapi.NewChaincode(factoryContract, logisticsContract, depotContract, hqContract)
	if err != nil {
		log.Panicf("Could not create chaincode: %v", err)
	}

	// Start chaincode
	if err := chaincode.Start(); err != nil {
		log.Panicf("Failed to start chaincode: %v", err)
	}
}
