# Blockchain-Based Ammunition Supply Chain Traceability

A secure and transparent system for tracking ammunition batches across the supply chain using Hyperledger Fabric. The solution ensures traceability, prevents counterfeiting, and enables HQ oversight without exposing sensitive data.

## 🚀 Features
- End-to-end traceability of ammunition batches
- Immutable ledger to prevent tampering
- Controlled access using Private Data Collections (PDCs)
- Status updates and audit capabilities
- HQ monitoring without data exposure

## 🏗 System Architecture
The project involves four organizations:
1. **Factory (Org1)** – Creates batches
2. **Logistics Division (Org2)** – Transports batches
3. **Army Depot (Org3)** – Stores and manages stock
4. **Defense HQ (Org4)** – Audits and monitors the entire supply chain

## 🛠 Technologies Used
- Hyperledger Fabric
- Smart Contracts (Chaincode)
- Private Data Collections (PDCs)
- Docker

## This directory contains 4 folders:
1. Chaincode- It contains the business logic for the developed application
2. SC-Network - It contains the files for bringing up the production network 
3. Supply-App - It contains the UI component for running the application
4. SupplyChain-Minifab- The project can be run using minifab from this directory.


## Chaincode and its functions
- Factory Contract - createBatch, ReadBatchBoth, ReadAllBatches
- Logistics Contract - TransferBatch, GetRouteInfo, GetAllRouteInfo
- Depot Contract- ReceiveBatch, ReadAllDepotBatches
- HQ Contract - GetBatchesWithPagination, QueryFactoryBatchesByStatus, QueryBatchesByBatchID, QueryBatchesByBatchID, QueryPublicBatchesByRange, GetBatchHistoryPublic

## There are two ways to run this project: 
1. Using Minifab (Move to SupplyChain_Minifab)
2. Using CLI (Move to SC-Network)