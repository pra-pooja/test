package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func submitTxnFn(
	organization string,
	channelName string,
	chaincodeName string,
	contractName string,
	txnType string,
	privateData map[string][]byte,
	endorsingOrgs []string,
	txnName string,
	args ...string,
) string {
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return fmt.Sprintf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	fmt.Printf("\n--> Submitting Transaction: %s\n", txnName)

	switch txnType {
	case "invoke":
		result, err := contract.Submit(
			txnName,
			client.WithArguments(args...),
		)
		if err != nil {
			return fmt.Sprintf("Failed to submit transaction: %v", err)
		}
		return fmt.Sprintf("*** Transaction submitted successfully: %s\n", string(result))

	case "query":
		evaluateResult, err := contract.EvaluateTransaction(txnName, args...)
		if err != nil {
			return fmt.Sprintf("Failed to evaluate transaction: %v", err)
		}
		if isByteSliceEmpty(evaluateResult) {
			return string(evaluateResult)
		}
		return formatJSON(evaluateResult)

	case "private":
		//fmt.Printf("DEBUG: In private case, endorsingOrgs = %v, len = %d\n", endorsingOrgs, len(endorsingOrgs))
		// MUST use WithEndorsingOrganizations for private data
		if len(endorsingOrgs) > 0 {
			result, err := contract.Submit(
				txnName,
				client.WithArguments(args...),
				client.WithTransient(privateData),
				client.WithEndorsingOrganizations(endorsingOrgs...),
			)
			if err != nil {
				return fmt.Sprintf("Failed to submit private transaction: %v", err)
			}
			return fmt.Sprintf("*** Private transaction committed successfully\nResult: %s\n", string(result))
		} else {
			return "Error: Private data requires endorsing organizations"
		}
	}

	return ""
}
