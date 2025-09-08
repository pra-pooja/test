package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Batch struct {
	BatchID         string `json:"BatchID"`
	Type            string `json:"Type"`
	Quantity        int    `json:"Quantity"`
	ManufactureDate string `json:"ManufactureDate"`
	ExpiryDate      string `json:"ExpiryDate"`
	Status          string `json:"Status"`
	Composition     string `json:"Composition"`
	Inspection      string `json:"Inspection"`
	Serials         string `json:"Serials"`
}

func main() {
	router := gin.Default()

	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	// UI page
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Ammunition Supply Chain Dashboard",
		})
	})
	
	router.GET("/factory", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "factory.html", gin.H{
			"title": "Factory Dashboard",
		})
	})

	// POST: Create Batch (private data)
	router.POST("/api/batch", func(ctx *gin.Context) {
		var req Batch
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		fmt.Println("request", req)

		// Convert Quantity to string
		qty := strconv.Itoa(req.Quantity)

		// Prepare transient data for private fields
		privateData := map[string][]byte{
			"composition": []byte(req.Composition),
			"inspection":  []byte(req.Inspection),
			"serials":     []byte(req.Serials),
		}
		endorsingOrgs := []string{"FactoryMSP", "LogisticsMSP", "DepotMSP", "HQMSP"}

		result := submitTxnFn(
			"factory",
			"supplychannel",
			"supplychain",
			"FactoryContract",
			"private",
			privateData,
			endorsingOrgs, // <-- must include all orgs in PDC
			"CreateBatch",
			req.BatchID,
			req.Type,
			qty,
			req.ManufactureDate,
			req.ExpiryDate,
			req.Status,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"message": result,
			"batchID": req.BatchID,
		})
	})

	// GET: Query Batch
	router.GET("/api/batch/:id", func(ctx *gin.Context) {
		batchID := ctx.Param("id")

		result := submitTxnFn(
			"factory",
			"supplychannel",
			"supplychain",
			"FactoryContract",
			"query",
			nil,             // no private data
			nil,             // discovery for query
			"ReadBatchBoth", // chaincode query function
			batchID,
		)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// Run server
	router.Run("localhost:3001")
}

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type Batch struct {
// 	BatchID         string `json:"BatchID"`
// 	Type            string `json:"Type"`
// 	Quantity        int    `json:"Quantity"`
// 	ManufactureDate string `json:"ManufactureDate"`
// 	ExpiryDate      string `json:"ExpiryDate"`
// 	Status          string `json:"Status"`
// 	Composition     string `json:"Composition"`
// 	Inspection      string `json:"Inspection"`
// 	Serials         string `json:"Serials"`
// }

// func main() {
// 	router := gin.Default()

// 	router.Static("/public", "./public")
// 	router.LoadHTMLGlob("templates/*")

// 	router.GET("/", func(ctx *gin.Context) {
// 		ctx.HTML(http.StatusOK, "factory.html", gin.H{
// 			"title": "Factory Dashboard",
// 		})
// 	})

// 	router.POST("/api/batch", func(ctx *gin.Context) {
// 		var req Batch
// 		if err := ctx.BindJSON(&req); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
// 			return
// 		}
// 		fmt.Println("request", req)
// 		// Convert Quantity to string
// 		qty := strconv.Itoa(req.Quantity)

// 		// Prepare transient data for private fields
// 		privateData := map[string][]byte{
// 			"composition": []byte(req.Composition),
// 			"inspection":  []byte(req.Inspection),
// 			"serials":     []byte(req.Serials),
// 		}

// 		endorsingOrgs := []string{"FactoryMSP", "LogisticsMSP", "DepotMSP", "HQMSP"}

// 		result := submitTxnFn(
// 			"factory",
// 			"supplychannel",
// 			"supplychain",
// 			"FactoryContract",
// 			"private",
// 			privateData,
// 			endorsingOrgs, // All required orgs
// 			"CreateBatch",
// 			req.BatchID,
// 			req.Type,
// 			qty,
// 			req.ManufactureDate,
// 			req.ExpiryDate,
// 			req.Status,
// 		)

// 		ctx.JSON(http.StatusOK, gin.H{
// 			"message": result,
// 			"batchID": req.BatchID,
// 		})
// 	})
// 	router.GET("/api/batch/:id", func(ctx *gin.Context) {
// 		batchID := ctx.Param("id")

// 		result := submitTxnFn(
// 			"factory",
// 			"supplychannel",
// 			"supplychain",
// 			"FactoryContract",
// 			"query",
// 			nil,
// 			nil,
// 			"ReadBatchBoth",
// 			batchID,
// 		)

// 		ctx.JSON(http.StatusOK, gin.H{"data": result})
// 	})

// 	router.Run("localhost:3001")
// }

//router.GET("/", func(ctx *gin.Context) {
// 	result := submitTxnFn("factory", "supplychannel", "SupplyChain", "FactoryContract", "invoke", make(map[string][]byte), "CreateBatch")

// 	var batch1 []BatchPublic
// 	//var batch2 []BatchPrivate
// 	if len(result) > 0 {
// 		// Unmarshal the JSON array string into the cars slice
// 		if err := json.Unmarshal([]byte(result), &batch1); err != nil {
// 			fmt.Println("Error:", err)
// 			return
// 		}
// 	}

// 	ctx.HTML(http.StatusOK, "index.html", gin.H{
// 		"title": "Supply App", "batchList": batch1,
// 	})
// })
// type Depot struct {
// 	BatchID    string `json:"BatchID"`
// 	Quantity   int    `json:"Quantity"`
// 	ReceivedBy string `json:"ReceivedBy"`
// 	ReceivedAt string `json:"ReceivedAt"`
// }
// type RouteInfo struct {
// 	BatchID     string `json:"BatchID"`
// 	FromOrg     string `json:"FromOrg"`
// 	ToOrg       string `json:"ToOrg"`
// 	Route       string `json:"Route"`
// 	ConvoyID    string `json:"ConvoyId"`
// 	Transporter string `json:"Transporter"`
// 	Timestamp   string `json:"Timestamp"`
// }

// type BatchHistory struct {
// 	Record    *BatchPublic `json:"Record"`
// 	TxId      string       `json:"TxId"`
// 	Timestamp string       `json:"Timestamp"`
// 	IsDelete  bool         `json:"IsDelete"`
// }
