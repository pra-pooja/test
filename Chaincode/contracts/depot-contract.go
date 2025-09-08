package contracts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Depot struct {
	BatchID    string `json:"batchID"`
	Quantity   int    `json:"quantity"`
	ReceivedBy string `json:"receivedBy"`
	ReceivedAt string `json:"receivedAt"`
}

type DepotContract struct {
	contractapi.Contract
}

const depotCollection string = "pdcDepotHQ"

func (d *DepotContract) ReceiveBatch(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	if clientOrgID == "DepotMSP" {

		//Check if batch already exists in Depot PDC
		exists, err := ctx.GetStub().GetPrivateDataHash(depotCollection, batchID)

		if err != nil {
			return "", fmt.Errorf("failed to check Depot batch existence: %v", err)
		}
		if exists != nil {
			return "", fmt.Errorf("batch %s already recorded in depot stock", batchID)
		}

		//Check if batch exists in public ledger
		publicData, err := ctx.GetStub().GetState(batchID)
		if err != nil {
			return "", fmt.Errorf("failed to read public batch: %v", err)
		}
		if publicData == nil {
			return "", fmt.Errorf("batch %s does not exist in public ledger", batchID)
		}

		var depotBatch Depot
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data: %s", err)
		}
		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide depot details in transient data (quantity)")
		}

		// quantity (required)
		qtyBytes, ok := transientData["quantity"]
		if !ok {
			return "", fmt.Errorf("the quantity was not specified in transient data")
		}

		qtyStr := string(qtyBytes) // convert []byte -> string
		qtyInt, err := strconv.Atoi(qtyStr)
		if err != nil {
			return "", fmt.Errorf("invalid quantity in transient data: %v", err)
		}

		// Use qtyInt
		depotBatch.Quantity = qtyInt

		depotBatch.BatchID = batchID
		depotBatch.ReceivedBy = clientOrgID
		depotBatch.ReceivedAt = time.Now().Format(time.RFC3339)

		depotBytes, err := json.Marshal(depotBatch)
		if err != nil {
			return "", fmt.Errorf("failed to marshal depot batch: %v", err)
		}
		err = ctx.GetStub().PutPrivateData(depotCollection, batchID, depotBytes)
		if err != nil {
			return "", fmt.Errorf("failed to put depot batch in PDC: %v", err)
		}
		var pub BatchPublic
		err = json.Unmarshal(publicData, &pub)
		if err == nil {
			pub.Status = "Received"
			updatedBytes, _ := json.Marshal(pub)
			_ = ctx.GetStub().PutState(batchID, updatedBytes)
		}
		return fmt.Sprintf("Batch %s successfully received by %s", batchID, clientOrgID), nil
	}
	return "Batch Received and Status Update to received", nil
}

func (d *DepotContract) ReadAllDepotBatches(ctx contractapi.TransactionContextInterface) ([]*Depot, error) {
	// Get caller's MSP
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client identity: %v", err)
	}

	// Only Depot and HQ can read all depot batches
	if clientOrgID != "DepotMSP" && clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("user %s not allowed to access depot data", clientOrgID)
	}

	// Query all entries in depot PDC
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(depotCollection, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read all depot batches: %v", err)
	}
	defer resultsIterator.Close()

	var depotBatches []*Depot

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error iterating depot batches: %v", err)
		}

		var depotBatch Depot
		err = json.Unmarshal(queryResponse.Value, &depotBatch)
		if err != nil {
			// skip invalid entries
			continue
		}

		depotBatches = append(depotBatches, &depotBatch)
	}

	return depotBatches, nil
}
