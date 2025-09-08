package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type RouteInfo struct {
	BatchID     string `json:"batchID"`
	FromOrg     string `json:"fromOrg"`
	ToOrg       string `json:"toOrg"`
	Route       string `json:"route"`
	ConvoyID    string `json:"convoyId"`
	Transporter string `json:"transporter"`
	Timestamp   string `json:"timestamp"`
}

type LogisticsContract struct {
	contractapi.Contract
}

const logisticsCollection string = "pdcLogisticsHQ"

func (l *LogisticsContract) TransferBatch(ctx contractapi.TransactionContextInterface, batchID, fromOrg, toOrg string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	// if clientOrgID == "LogisticsMSP" {

	if clientOrgID == "LogisticsMSP" {
		exists, err := (&FactoryContract{}).BatchExists(ctx, batchID)
		if err != nil {
			return "", fmt.Errorf("failed to check if batch %s exists: %v", batchID, err)
		} else if exists {
			pubBytes, err := ctx.GetStub().GetState(batchID)
			if err != nil {
				return "", fmt.Errorf("failed to get Batch: %v", err)
			}

			var route RouteInfo
			transientData, err := ctx.GetStub().GetTransient()
			if err != nil {
				return "", fmt.Errorf("could not fetch transient data: %s", err)
			}
			if len(transientData) == 0 {
				return "", fmt.Errorf("please provide route info in transient data (route, convoyID, transporter)")
			}

			// Route (required)
			rout, exists := transientData["route"]
			if !exists {
				return "", fmt.Errorf("the route was not specified in transient data")
			}
			route.Route = string(rout)

			// ConvoyID (required)
			convoyID, exists := transientData["convoyID"]
			if !exists {
				return "", fmt.Errorf("the convoyID was not specified in transient data")
			}
			route.ConvoyID = string(convoyID)

			// Transporter (required)
			transporter, exists := transientData["transporter"]
			if !exists {
				return "", fmt.Errorf("the transporter was not specified in transient data")
			}
			route.Transporter = string(transporter)

			// Fill other details from params
			route.BatchID = batchID
			route.FromOrg = fromOrg
			route.ToOrg = toOrg
			route.Timestamp = time.Now().Format("2006-01-02 15:04:05")

			privBytes, err := json.Marshal(route)
			if err != nil {
				return "", fmt.Errorf("failed to marshal route info: %v", err)
			}

			// 🔹 FIX: use consistent key for both Put and Get
			err = ctx.GetStub().PutPrivateData(logisticsCollection, batchID, privBytes)
			if err != nil {
				return "", fmt.Errorf("failed to put private route data: %v", err)
			}

			// update public status
			var pub BatchPublic
			err = json.Unmarshal(pubBytes, &pub)
			if err == nil {
				pub.Status = "In-Transit"
				updatedBytes, _ := json.Marshal(pub)
				_ = ctx.GetStub().PutState(batchID, updatedBytes)
			}

		}
		return fmt.Sprintf("Batch %s transferred from %s to %s", batchID, fromOrg, toOrg), nil
	}
	return "batch transferred and status updated", nil
}


func (l *LogisticsContract) GetRouteInfo(ctx contractapi.TransactionContextInterface, batchID string) (*RouteInfo, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	fmt.Println("Client Identity: ", clientOrgID)

	if clientOrgID == "LogisticsMSP" || clientOrgID == "HQMSP" {
		// 🔹 FIX: match the key with TransferBatch
		routeBytes, err := ctx.GetStub().GetPrivateData(logisticsCollection, batchID)
		if err != nil {
			return nil, fmt.Errorf("failed to read private route info for batch %s: %s", batchID, err)
		}
		if routeBytes == nil {
			return nil, fmt.Errorf("no route info found for batch %s", batchID)
		}
		var route RouteInfo
		err = json.Unmarshal(routeBytes, &route)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal route info: %s", err)
		}
		return &route, nil
	} else {
		return nil, fmt.Errorf("user is not allowed to perform the read batch transaction")
	}
}

func (l *LogisticsContract) GetAllRouteInfo(ctx contractapi.TransactionContextInterface) ([]*RouteInfo, error) {
	// allow only LogisticsMSP and HQMSP to query
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("identity retrieval error: %v", err)
	}
	if clientOrgID != "LogisticsMSP" && clientOrgID != "HQMSP" {
		return nil, fmt.Errorf("user from %s is not allowed to read all route info", clientOrgID)
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(logisticsCollection, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get private data from %s: %v", logisticsCollection, err)
	}
	defer resultsIterator.Close()

	var routes []*RouteInfo

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate route info: %v", err)
		}

		var route RouteInfo
		err = json.Unmarshal(queryResponse.Value, &route)
		if err != nil {
			continue
		}
		routes = append(routes, &route)
	}

	return routes, nil
}
