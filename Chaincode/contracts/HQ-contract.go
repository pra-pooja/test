package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// HQContract provides auditing functions for HQ
type HQContract struct {
	contractapi.Contract
}

type HistoryQueryResult struct {
	Record    *BatchPublic `json:"record"`
	TxId      string       `json:"txId"`
	Timestamp string       `json:"timestamp"`
	IsDelete  bool         `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records []*BatchPublic `json:"records"`
	//FetchedRecordsCount int32          `json:"fetchedRecordsCount"`
	Bookmark string `json:"bookmark"`
}

const (
	factoryHQCollection   = "pdcFactoryHQ"
	logisticsHQCollection = "pdcLogisticsHQ"
	depotHQCollection     = "pdcDepotHQ"
)

func (h *HQContract) GetBatchesWithPagination(ctx contractapi.TransactionContextInterface, field string, value string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s"}}`, field, value)

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get paginated query results: %v", err)
	}
	defer resultsIterator.Close()

	var results []*BatchPublic
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		var batch BatchPublic
		if err := json.Unmarshal(queryResponse.Value, &batch); err == nil {
			results = append(results, &batch)
		}
	}

	return &PaginatedQueryResult{
		Records:  results,
		Bookmark: responseMetadata.Bookmark,
	}, nil
}

func (h *HQContract) QueryFactoryBatchesByStatus(ctx contractapi.TransactionContextInterface, status string) ([]*BatchPublic, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	query := fmt.Sprintf(`{"selector":{"status":"%s"}}`, status)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer resultsIterator.Close()

	var results []*BatchPublic
	for resultsIterator.HasNext() {
		resp, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var detail BatchPublic
		if err := json.Unmarshal(resp.Value, &detail); err == nil {
			results = append(results, &detail)
		}
	}
	return results, nil
}
func (h *HQContract) QueryBatchesByBatchID(ctx contractapi.TransactionContextInterface, batchID string) ([]map[string]interface{}, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	query := fmt.Sprintf(`{"selector":{"batchID":"%s"}}`, batchID)

	var results []map[string]interface{}

	// ---------- 1. Public Data ----------
	pubIter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("public query failed: %v", err)
	}
	defer pubIter.Close()

	for pubIter.HasNext() {
		resp, err := pubIter.Next()
		if err != nil {
			return nil, err
		}
		var batch BatchPublic
		if err := json.Unmarshal(resp.Value, &batch); err == nil {
			results = append(results, map[string]interface{}{
				"source": "public",
				"batch":  batch,
			})
		}
	}

	// ---------- 2. Private Data (FactoryHQ) ----------
	privIter, err := ctx.GetStub().GetPrivateDataQueryResult(factoryHQCollection, query)
	if err != nil {
		return nil, fmt.Errorf("private query failed: %v", err)
	}
	defer privIter.Close()

	for privIter.HasNext() {
		resp, err := privIter.Next()
		if err != nil {
			return nil, err
		}
		var batch BatchPrivate
		if err := json.Unmarshal(resp.Value, &batch); err == nil {
			results = append(results, map[string]interface{}{
				"source": "private",
				"batch":  batch,
			})
		}
	}

	return results, nil
}

func (h *HQContract) QueryAllPublicBatches(ctx contractapi.TransactionContextInterface) ([]*BatchPublic, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	results := []*BatchPublic{}

	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("range query failed: %v", err)
	}
	defer iter.Close()

	for iter.HasNext() {
		resp, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var batch BatchPublic
		if err := json.Unmarshal(resp.Value, &batch); err == nil {
			results = append(results, &batch)
		}
	}
	return results, nil
}

// QueryPublicBatchesByRange fetches public batches within a key range
func (h *HQContract) QueryPublicBatchesByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*BatchPublic, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	results := []*BatchPublic{}

	// Pass user-specified range instead of ("", "")
	iter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("range query failed: %v", err)
	}
	defer iter.Close()

	for iter.HasNext() {
		resp, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var batch BatchPublic
		if err := json.Unmarshal(resp.Value, &batch); err == nil {
			results = append(results, &batch)
		}
	}
	return results, nil
}

func (h *HQContract) GetBatchHistoryPublic(ctx contractapi.TransactionContextInterface, batchID string) ([]map[string]interface{}, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("access denied to %s ", clientOrgID)
	}
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for batch %s: %v", batchID, err)
	}
	defer resultsIterator.Close()

	var records []map[string]interface{}
	for resultsIterator.HasNext() {
		response, _ := resultsIterator.Next()
		var value interface{}
		if response.Value != nil {
			_ = json.Unmarshal(response.Value, &value)
		}
		record := map[string]interface{}{
			"TxID":      response.TxId,
			"Timestamp": response.Timestamp,
			"Value":     value,
			"IsDelete":  response.IsDelete,
		}
		records = append(records, record)
	}
	return records, nil
}
