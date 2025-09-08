##Blockchain-Based Ammunition Supply Chain Traceability

A secure and transparent system for tracking ammunition batches across the supply chain using Hyperledger Fabric. The solution ensures traceability, prevents counterfeiting, and enables HQ oversight without exposing sensitive data.

##Steps to run the project:
1. Execute the follwoing commands in the terminal:

#"-----------Register the ca admin for each organization—----------------
docker compose -f docker/docker-compose-ca.yaml up -d

sudo chmod -R 777 organizations/

#------------Register and enroll the users for each organization—-----------
chmod +x registerEnroll.sh
./registerEnroll.sh

#—-------------Build the infrastructure—-----------------

docker compose -f docker/docker-compose-4org.yaml up -d

#-------------Generate the genesis block—-------------------------------

export FABRIC_CFG_PATH=${PWD}/config

export CHANNEL_NAME=supplychannel

configtxgen -profile SupplyChainChannel -outputBlock ${PWD}/channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME

#----- Create the application channel----
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

**************** peer0_Factory terminal ********************
***Build the core.yaml in peercfg folder

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=supplychannel
export CORE_PEER_LOCALMSPID=FactoryMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/factory.supplychain.com/users/Admin@factory.supplychain.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem
export FACTORY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export LOGISTICS_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export DEPOT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export HQ_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt

—---------------Join peer to the channel—-------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

#—-------------Factory anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.FactoryMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.factory.supplychain.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile $ORDERER_CA

#—---------------package chaincode—-------------

peer lifecycle chaincode package supplychain.tar.gz --path ${PWD}/../Chaincode/ --lang golang --label supplychain_1.0

#—---------------install chaincode in Factory peer—-------------

peer lifecycle chaincode install supplychain.tar.gz

peer lifecycle chaincode queryinstalled

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid supplychain.tar.gz)

#---------------Approve chaincode in Factory peer—-------------

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

**************** peer0_Logistics terminal *****************

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=supplychannel
export CORE_PEER_LOCALMSPID=LogisticsMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/users/Admin@logistics.supplychain.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem
export FACTORY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export LOGISTICS_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export DEPOT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export HQ_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt

—------------ Join peer to the channel ---------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

#—-------------Logistics anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.LogisticsMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.logistics.supplychain.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile $ORDERER_CA

#—---------------install chaincode in Logistics peer—-------------"

peer lifecycle chaincode install supplychain.tar.gz

peer lifecycle chaincode queryinstalled

—---------------Approve chaincode in Logistics peer—-------------"

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid supplychain.tar.gz)

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent


**************** peer0_Depot terminal ******************

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=supplychannel 
export CORE_PEER_LOCALMSPID=DepotMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/depot.supplychain.com/users/Admin@depot.supplychain.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem
export FACTORY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export LOGISTICS_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export DEPOT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export HQ_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt

—-------------- Join peer to the channel ----------------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
 "—-------------Depot anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
 

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.DepotMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.depot.supplychain.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile $ORDERER_CA
 

peer channel getinfo -c $CHANNEL_NAME

 "—---------------install chaincode in Depot peer—-------------"

peer lifecycle chaincode install supplychain.tar.gz
 

peer lifecycle chaincode queryinstalled

 "—---------------Approve chaincode in Depot peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

**************** peer0_HQ terminal ******************

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=supplychannel 
export CORE_PEER_LOCALMSPID=HQMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:12051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/users/Admin@HQ.supplychain.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem
export FACTORY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/factory.supplychain.com/peers/peer0.factory.supplychain.com/tls/ca.crt
export LOGISTICS_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/logistics.supplychain.com/peers/peer0.logistics.supplychain.com/tls/ca.crt
export DEPOT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/depot.supplychain.com/peers/peer0.depot.supplychain.com/tls/ca.crt
export HQ_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/HQ.supplychain.com/peers/peer0.HQ.supplychain.com/tls/ca.crt

—-------------- Join peer to the channel ----------------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

"—-------------HQ anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
 

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.HQMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.HQ.supplychain.com","port": 12051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile $ORDERER_CA
 

peer channel getinfo -c $CHANNEL_NAME

"—---------------install chaincode in HQ peer—-------------"

peer lifecycle chaincode install supplychain.tar.gz
 

peer lifecycle chaincode queryinstalled

"—---------------Approve chaincode in HQ peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
 

"—---------------Commit chaincode in HQ peer—-------------"

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name supplychain --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --channelID $CHANNEL_NAME --name supplychain --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT

peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name supplychain --cafile $ORDERER_CA


##Now, the network is ready with the chaincode installed and deployed over the network. All the four organizations are running in four terminals.. Now execute the following commands to invoke and query the chaincode.

**************** Factory terminal ********************
xport COMPOSITION=$(echo -n "Powder,Explosives" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Passed" | base64 | tr -d \\n)
export SERIALS=$(echo -n "1001-1009" | base64 | tr -d \\n)


peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["BATCH-001","5.62mm Mortar Shells","1000","2025-06-15","2035-06-15","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"

export COMPOSITION=$(echo -n "Liquid,Chemicals" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Passed" | base64 | tr -d \\n)
export SERIALS=$(echo -n "2001-2010" | base64 | tr -d \\n)

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["Batch-002","Industrial Solvent","500","2025-07-01","2030-07-01","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"

export COMPOSITION=$(echo -n "Solid,Metals" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Failed" | base64 | tr -d \\n)
export SERIALS=$(echo -n "3001-3015" | base64 | tr -d \\n)

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["Batch-003","Steel Rods","1500","2025-08-10","2040-08-10","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"


export COMPOSITION=$(echo -n "Gas,Flammable" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Passed" | base64 | tr -d \\n)
export SERIALS=$(echo -n "4001-4020" | base64 | tr -d \\n)

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["Batch-004","Propane Canisters","800","2025-05-20","2030-07-01","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"
  
export COMPOSITION=$(echo -n "Powder,Pharmaceuticals" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Passed" | base64 | tr -d \\n)
export SERIALS=$(echo -n "5001-5010" | base64 | tr -d \\n)

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["Batch-005","Painkillers","2000","2025-09-05","2028-09-05","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"

export COMPOSITION=$(echo -n "Granules,Agricultural" | base64 | tr -d \\n)
export INSPECTION=$(echo -n "Pending" | base64 | tr -d \\n)
export SERIALS=$(echo -n "6001-6025" | base64 | tr -d \\n)

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls --cafile $ORDERER_CA \
  -C $CHANNEL_NAME \
  -n supplychain \
  --peerAddresses localhost:7051 --tlsRootCertFiles $FACTORY_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $LOGISTICS_PEER_TLSROOTCERT \
  --peerAddresses localhost:11051 --tlsRootCertFiles $DEPOT_PEER_TLSROOTCERT \
  --peerAddresses localhost:12051 --tlsRootCertFiles $HQ_PEER_TLSROOTCERT \
  -c '{"function":"CreateBatch","Args":["Batch-006","Fertilizer","3000","2025-04-15","2030-07-01","CREATED"]}' \
  --transient "{\"composition\":\"$COMPOSITION\",\"inspection\":\"$INSPECTION\",\"serials\":\"$SERIALS\"}"

-------------ReadBatchBoth----------
peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"ReadBatchBoth","Args":["Batch-001"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"ReadBatchBoth","Args":["Batch-002"]}'

---------------ReadAllBatches---------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"Args":["ReadAllBatches"]}'


**************** Logistics terminal ********************

# ---------- Transfer Batch-001 ----------
export ROUTE=$(echo -n "Route-A" | base64 | tr -d \\n)
export CONVOY=$(echo -n "CONVOY-101" | base64 | tr -d \\n)
export TRANSPORTER=$(echo -n "Transporter Alpha" | base64 | tr -d \\n)


peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"LogisticsContract:TransferBatch","Args":["Batch-001","FactoryMSP","DepotMSP"]}'   --transient "{\"route\":\"$ROUTE\",\"convoyID\":\"$CONVOY\",\"transporter\":\"$TRANSPORTER\"}"

# ---------- Transfer Batch-002 ----------
export ROUTE=$(echo -n "Route-B" | base64 | tr -d \\n)
export CONVOY=$(echo -n "CONVOY-102" | base64 | tr -d \\n)
export TRANSPORTER=$(echo -n "Transporter Beta" | base64 | tr -d \\n)

peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"LogisticsContract:TransferBatch","Args":["Batch-002","FactoryMSP","DepotMSP"]}'   --transient "{\"route\":\"$ROUTE\",\"convoyID\":\"$CONVOY\",\"transporter\":\"$TRANSPORTER\"}"

# ---------- Transfer Batch-003 ----------
export ROUTE=$(echo -n "Route-C" | base64 | tr -d \\n)
export CONVOY=$(echo -n "CONVOY-103" | base64 | tr -d \\n)
export TRANSPORTER=$(echo -n "Transporter Alpha" | base64 | tr -d \\n)

peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"LogisticsContract:TransferBatch","Args":["Batch-003","FactoryMSP","DepotMSP"]}'   --transient "{\"route\":\"$ROUTE\",\"convoyID\":\"$CONVOY\",\"transporter\":\"$TRANSPORTER\"}"


#--------------GetRouteInfo----------------
peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"LogisticsContract:GetRouteInfo","Args":["Batch-001"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"LogisticsContract:GetRouteInfo","Args":["Batch-002"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"LogisticsContract:GetRouteInfo","Args":["Batch-003"]}'

#--------------GetAlRouteInfo----------------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"Args":["LogisticsContract:GetAllRouteInfo"]}'


**************** Depot terminal ********************
export QUANTITY=$(echo -n "500" | base64 | tr -d \\n)

peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"DepotContract:ReceiveBatch","Args":["Batch-001"]}'   --transient "{\"quantity\":\"$QUANTITY\"}"


export QUANTITY=$(echo -n "5000" | base64 | tr -d \\n)

peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"DepotContract:ReceiveBatch","Args":["Batch-002"]}'   --transient "{\"quantity\":\"$QUANTITY\"}"

export QUANTITY=$(echo -n "500" | base64 | tr -d \\n)

peer chaincode invoke   -o localhost:7050   --ordererTLSHostnameOverride orderer.supplychain.com   --tls   --cafile "$ORDERER_CA"   -C "$CHANNEL_NAME"   -n supplychain   --peerAddresses localhost:7051 --tlsRootCertFiles "$FACTORY_PEER_TLSROOTCERT"   --peerAddresses localhost:9051 --tlsRootCertFiles "$LOGISTICS_PEER_TLSROOTCERT"   --peerAddresses localhost:11051 --tlsRootCertFiles "$DEPOT_PEER_TLSROOTCERT"   --peerAddresses localhost:12051 --tlsRootCertFiles "$HQ_PEER_TLSROOTCERT"   -c '{"function":"DepotContract:ReceiveBatch","Args":["Batch-003"]}'   --transient "{\"quantity\":\"$QUANTITY\"}"

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"Args":["DepotContract:ReadAllDepotBatches"]}'

**************** HQ terminal ********************

----------------1.QueryFactoryBatchesByStatus--------------------
peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryFactoryBatchesByStatus","Args":["CREATED"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryFactoryBatchesByStatus","Args":["In-Transit"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryFactoryBatchesByStatus","Args":["Received"]}'

----------------2.QueryBatchesByBatchID--------------------
peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryBatchesByBatchID","Args":["Batch-001"]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryBatchesByBatchID","Args":["Batch-002"]}'

----------------3.QueryAllPublicBatches--------------------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"Args":["HQContract:QueryAllPublicBatches"]}'

----------------4.GetBatchHistoryPublic--------------------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:GetBatchHistoryPublic","Args":["Batch-002"]}'

----------------5.GetBatchesWithPagination (Quantity can't be used as selector)------------------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:GetBatchesWithPagination","Args":["status","CREATED","5",""]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:GetBatchesWithPagination","Args":["type","Steel Rods","5",""]}'

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:GetBatchesWithPagination","Args":["expiryDate","2030-07-01","5",""]}'

----------------6.QueryPublicBatchesByRange (Only Show Batch Public Details)--------------------

peer chaincode query   -C $CHANNEL_NAME   -n supplychain   -c '{"function":"HQContract:QueryPublicBatchesByRange","Args":["BATCH-001","BATCH-006"]}'

## TO stop the network:
./stopSupplyNetwork.sh



##UI component is exist for only Factory organization .i.e. only createBatch and readBatchBoth functionalities are working in UI. It will be completed in next version


#To run the project:
1. Run the following in the ubuntu terminal ( from SC_Network)
./startSupplyNetwork.sh

2. Go inside the Supply-App folder and run the following command in terminal-
go run .

## Open the link : localhost:3001 in browser

3. To stop the network, run the following:
./stopSupplyNetwork.sh




