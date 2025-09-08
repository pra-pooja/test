package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type BatchPublic struct {
	BatchID         string `json:"batchID"`
	Type            string `json:"type"`
	Quantity        int    `json:"quantity"`
	ManufactureDate string `json:"manufactureDate"`
	ExpiryDate      string `json:"expiryDate"`
	Status          string `json:"status"`
}

type BatchPrivate struct {
	BatchID     string `json:"batchID"`
	Composition string `json:"composition"`
	Inspection  string `json:"inspection"`
	Serials     string `json:"serials"`
}

// =============================
// Factory Smart Contract
// =============================

const collectionName string = "pdcFactoryHQ"

type FactoryContract struct {
	contractapi.Contract
}

func (c *FactoryContract) BatchExists(ctx contractapi.TransactionContextInterface, BatchId string) (bool, error) {
	data, err := ctx.GetStub().GetState(BatchId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

func (f *FactoryContract) CreateBatch(ctx contractapi.TransactionContextInterface, BatchId string, Typ string, Qty int, dateOfManufacture string, dateofexpiry string, status string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	// if clientOrgID == "FactoryMSP" {
	if clientOrgID == "FactoryMSP" {
		exists, err := f.BatchExists(ctx, BatchId)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the batch, %s already exists", BatchId)
		}

		batchpublic := BatchPublic{
			BatchID:         BatchId,
			Type:            Typ,
			Quantity:        Qty,
			ManufactureDate: dateOfManufacture,
			ExpiryDate:      dateofexpiry,
			Status:          status,
		}

		pubbytes, _ := json.Marshal(batchpublic)
		err = ctx.GetStub().PutState(BatchId, pubbytes)
		if err != nil {
			return "", err
		}

		fmt.Println("Created Batch ======= ", batchpublic)

		// --- Private Data Handling ---
		var batchPriv BatchPrivate
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}
		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of Composition, Inspection, Serials")
		}

		comp, exists := transientData["composition"]
		if !exists {
			return "", fmt.Errorf("the composition was not specified in transient data. Please try again")
		}
		batchPriv.Composition = string(comp)

		inspect, exists := transientData["inspection"]
		if !exists {
			return "", fmt.Errorf("the inspection was not specified in transient data. Please try again")
		}
		batchPriv.Inspection = string(inspect)

		serial, exists := transientData["serials"]
		if !exists {
			return "", fmt.Errorf("the serials was not specified in transient data. Please try again")
		}
		batchPriv.Serials = string(serial)

		batchPriv.BatchID = BatchId
		privbytes, err := json.Marshal(batchPriv)
		if err != nil {
			return "", fmt.Errorf("failed to marshal private batch data: %s", err)
		}

		err = ctx.GetStub().PutPrivateData(collectionName, BatchId, privbytes)
		if err != nil {
			return "", fmt.Errorf("could not write batchPrivate data: %s", err)
		}

		return fmt.Sprintf("successfully added batch %v with private data", BatchId), nil
	}

	return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
}

func (f *FactoryContract) ReadBatchBoth(ctx contractapi.TransactionContextInterface, BatchId string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "FactoryMSP" && clientOrgID != "HQMSP" {
		return "", fmt.Errorf("user from %s is not allowed to read batch data", clientOrgID)
	}

	// ---- Read Public Batch ----
	pubBytes, err := ctx.GetStub().GetState(BatchId)
	if err != nil {
		return "", fmt.Errorf("failed to read public batch: %v", err)
	}
	if pubBytes == nil {
		return "", fmt.Errorf("the batch %s does not exist in public state", BatchId)
	}

	// ---- Read Private Batch ----
	privBytes, err := ctx.GetStub().GetPrivateData(collectionName, BatchId)
	if err != nil {
		// If no access, still return public
		return string(pubBytes), nil
	}

	// ---- Combine into one JSON string ----
	combined := map[string]json.RawMessage{
		"public":  pubBytes,
		"private": privBytes,
	}
	out, err := json.Marshal(combined)
	if err != nil {
		return "", fmt.Errorf("failed to combine batch data: %v", err)
	}

	return string(out), nil
}

func (f *FactoryContract) ReadAllBatches(ctx contractapi.TransactionContextInterface) ([]map[string]interface{}, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "FactoryMSP" && clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("user from %s is not allowed to read batch data", clientOrgID)
	}
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}
	defer resultsIterator.Close()

	var allBatches []map[string]interface{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate state: %v", err)
		}

		// ---------- Public Batch ----------
		var publicBatch BatchPublic
		err = json.Unmarshal(queryResponse.Value, &publicBatch)
		if err != nil {
			// skip non-batch entries
			continue
		}

		// ---------- Private Batch ----------
		privateBytes, err := ctx.GetStub().GetPrivateData(collectionName, publicBatch.BatchID)
		var privateBatch BatchPrivate
		if privateBytes != nil && err == nil {
			_ = json.Unmarshal(privateBytes, &privateBatch)
		}

		// Merge both into one combined response
		combined := map[string]interface{}{
			"public":  publicBatch,
			"private": privateBatch,
		}

		allBatches = append(allBatches, combined)
	}

	return allBatches, nil
}
