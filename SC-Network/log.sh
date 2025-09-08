#!/bin/bash

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=supplychannel
export CORE_PEER_LOCALMSPID=LogisticsMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/users/Admin@logistics.supplychain.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem
export FACTORY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export LOGISTICS_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export DEPOT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export HQ_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt

sleep 2

echo "—---------------Join Logistics peer to the channel—-------------"

echo ${FABRIC_CFG_PATH}
sleep 2
peer channel join -b ${PWD}/channel-artifacts/${CHANNEL_NAME}.block
sleep 3

echo "-----channel List----"
peer channel list

echo "—-------------Logistics anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.LogisticsMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.logistics.supplychain.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile $ORDERER_CA
sleep 1

echo "—---------------package chaincode—-------------"

peer lifecycle chaincode package supplychain.tar.gz --path ${PWD}/../Chaincode/ --lang golang --label supplychain_1.0
sleep 1

echo "—---------------install chaincode in Logistics peer—-------------"

peer lifecycle chaincode install supplychain.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid supplychain.tar.gz)

echo "—---------------Approve chaincode in Logistics peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2

