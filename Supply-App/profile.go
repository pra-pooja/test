package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"factory": {
		CertPath:     "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/users/User1@factory.supplychain.com/msp/signcerts/cert.pem",
		KeyDirectory: "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/users/User1@factory.supplychain.com/msp/keystore/",
		TLSCertPath:  "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.factory.supplychain.com",
		MSPID:        "FactoryMSP",
	},

	"logistics": {
		CertPath:     "../SC-Network/organizations/peerOrganizations/logistics.supplychain.com/users/User1@logistics.supplychain.com/msp/signcerts/cert.pem",
		KeyDirectory: "../SC-Network/organizations/peerOrganizations/logistics.supplychain.com/users/User1@logistics.supplychain.com/msp/keystore/",
		TLSCertPath:  "../SC-Network/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.logistics.supplychain.com",
		MSPID:        "LogisticsMSP",
	},

	"depot": {
		CertPath:     "../SC-Network/organizations/peerOrganizations/depot.supplychain.com/users/User1@depot.supplychain.com/msp/signcerts/cert.pem",
		KeyDirectory: "../SC-Network/organizations/peerOrganizations/depot.supplychain.com/users/User1@depot.supplychain.com/msp/keystore/",
		TLSCertPath:  "../SC-Network/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.depot.supplychain.com",
		MSPID:        "DepotMSP",
	},

	"HQ": {
		CertPath:     "../SC-Network/organizations/peerOrganizations/HQ.supplychain.com/users/User1@HQ.supplychain.com/msp/signcerts/cert.pem",
		KeyDirectory: "../SC-Network/organizations/peerOrganizations/HQ.supplychain.com/users/User1@HQ.supplychain.com/msp/keystore/",
		TLSCertPath:  "../SC-Network/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:12051",
		GatewayPeer:  "peer0.HQ.supplychain.com",
		MSPID:        "HQMSP",
	},

	"factory2": {
		CertPath:     "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/users/User2@factory.supplychain.com/msp/signcerts/cert.pem",
		KeyDirectory: "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/users/User2@factory.supplychain.com/msp/keystore/",
		TLSCertPath:  "../SC-Network/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.factory.supplychain.com",
		MSPID:        "FactoryMSP",
	},

	"minifab-factory": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/factory.supplychain.com/users/Admin@factory.supplychain.com/msp/signcerts/Admin@factory.supplychain.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/factory.supplychain.com/users/Admin@factory.supplychain.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/factory.supplychain.com/peers/peer1.factory.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7003",
		GatewayPeer:  "peer1.factory.supplychain.com",
		MSPID:        "factory-supplychain-com",
	},

	"minifab-logistics": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/logistics.supplychain.com/users/Admin@logistics.supplychain.com/msp/signcerts/Admin@logistics.supplychain.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/logistics.supplychain.com/users/Admin@logistics.supplychain.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/logistics.supplychain.com/peers/peer1.logistics.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7004",
		GatewayPeer:  "peer0.logistics.supplychain.com",
		MSPID:        "logistics-supplychain-com",
	},

	"minifab-depot": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/depot.supplychain.com/users/Admin@depot.supplychain.com/msp/signcerts/Admin@depot.supplychain.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/depot.supplychain.com/users/Admin@depot.supplychain.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/depot.supplychain.com/peers/peer1.depot.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7005",
		GatewayPeer:  "peer0.depot.supplychain.com",
		MSPID:        "depot-supplychain-com",
	},

	"minifab-HQ": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/HQ.supplychain.com/users/Admin@HQ.supplychain.com/msp/signcerts/Admin@HQ.supplychain.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/HQ.supplychain.com/users/Admin@HQ.supplychain.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/HQ.supplychain.com/peers/peer1.HQ.supplychain.com/tls/ca.crt",
		PeerEndpoint: "localhost:7005",
		GatewayPeer:  "peer0.HQ.supplychain.com",
		MSPID:        "HQ-supplychain-com",
	},
}
