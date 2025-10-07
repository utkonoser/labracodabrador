package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type NetworkInfo struct {
	ChainID      uint64 `json:"chainId"`
	LatestBlock  uint64 `json:"latestBlock"`
	GasPrice     string `json:"gasPrice"`
	GasPriceGwei string `json:"gasPriceGwei"`
}

type BlockInfo struct {
	Number           uint64    `json:"number"`
	Hash             string    `json:"hash"`
	Timestamp        time.Time `json:"timestamp"`
	Miner            string    `json:"miner"`
	TransactionCount int       `json:"transactionCount"`
	GasUsed          uint64    `json:"gasUsed"`
}

type AddressInfo struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	BalanceEth string `json:"balanceEth"`
	TxCount    uint64 `json:"txCount"`
	IsContract bool   `json:"isContract"`
}

var client *ethclient.Client

func main() {
	// Подключение к Ethereum RPC
	var err error
	client, err = ethclient.Dial("http://nginx:8545")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// Создание роутера
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if req.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, req)
		})
	})

	// API Routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Network info
	api.HandleFunc("/network", getNetworkInfo).Methods("GET")

	// Blocks
	api.HandleFunc("/blocks/latest", getLatestBlocks).Methods("GET")
	api.HandleFunc("/blocks/{number}", getBlock).Methods("GET")

	// Addresses
	api.HandleFunc("/address/{address}", getAddressInfo).Methods("GET")

	// Transactions
	api.HandleFunc("/transactions/{hash}", getTransaction).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{
			Success: true,
			Data:    map[string]string{"status": "healthy"},
		})
	}).Methods("GET")

	log.Println("🚀 Labracodabrador API Server starting on :8081")
	log.Println("📡 Connected to Ethereum RPC at nginx:8545")
	log.Println("🔗 API Endpoints:")
	log.Println("   GET /api/v1/network - Network information")
	log.Println("   GET /api/v1/blocks/latest - Latest blocks")
	log.Println("   GET /api/v1/blocks/{number} - Get block by number")
	log.Println("   GET /api/v1/address/{address} - Address information")
	log.Println("   GET /api/v1/transactions/{hash} - Transaction details")
	log.Println("   GET /health - Health check")

	log.Fatal(http.ListenAndServe(":8081", r))
}

func getNetworkInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	// Получаем информацию о сети
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		sendError(w, "Failed to get chain ID", err)
		return
	}

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		sendError(w, "Failed to get latest block", err)
		return
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		sendError(w, "Failed to get gas price", err)
		return
	}

	// Конвертируем gas price в Gwei
	gasPriceGwei := float64(gasPrice.Int64()) / 1e9

	networkInfo := NetworkInfo{
		ChainID:      chainID.Uint64(),
		LatestBlock:  blockNumber,
		GasPrice:     gasPrice.String(),
		GasPriceGwei: fmt.Sprintf("%.2f", gasPriceGwei),
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    networkInfo,
	})
}

func getLatestBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		sendError(w, "Failed to get latest block number", err)
		return
	}

	// Получаем последние 6 блоков
	var blocks []BlockInfo
	for i := uint64(0); i < 6 && i <= latestBlock; i++ {
		blockNum := latestBlock - i

		block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			continue // Пропускаем блоки, которые не удалось получить
		}

		blockInfo := BlockInfo{
			Number:           block.Number().Uint64(),
			Hash:             block.Hash().Hex(),
			Timestamp:        time.Unix(int64(block.Time()), 0),
			Miner:            block.Coinbase().Hex(),
			TransactionCount: len(block.Transactions()),
			GasUsed:          block.GasUsed(),
		}

		blocks = append(blocks, blockInfo)
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    blocks,
	})
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	vars := mux.Vars(r)
	blockNumStr := vars["number"]

	blockNum, err := strconv.ParseInt(blockNumStr, 10, 64)
	if err != nil {
		sendError(w, "Invalid block number", err)
		return
	}

	block, err := client.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		sendError(w, "Failed to get block", err)
		return
	}

	blockInfo := BlockInfo{
		Number:           block.Number().Uint64(),
		Hash:             block.Hash().Hex(),
		Timestamp:        time.Unix(int64(block.Time()), 0),
		Miner:            block.Coinbase().Hex(),
		TransactionCount: len(block.Transactions()),
		GasUsed:          block.GasUsed(),
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    blockInfo,
	})
}

func getAddressInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	vars := mux.Vars(r)
	addressStr := vars["address"]

	address := common.HexToAddress(addressStr)

	// Получаем баланс
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		sendError(w, "Failed to get balance", err)
		return
	}

	// Получаем количество транзакций
	txCount, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		sendError(w, "Failed to get transaction count", err)
		return
	}

	// Проверяем, является ли адрес контрактом
	code, err := client.CodeAt(ctx, address, nil)
	isContract := len(code) > 0

	// Конвертируем баланс в ETH (используем big.Float для точности)
	balanceFloat := new(big.Float).SetInt(balance)
	ethFloat, _ := new(big.Float).SetString("1000000000000000000") // 1e18
	balanceEthFloat := new(big.Float).Quo(balanceFloat, ethFloat)

	// Проверяем, что результат положительный
	var balanceEth float64
	if balanceEthFloat.Sign() < 0 {
		balanceEth = 0.0
	} else {
		balanceEth, _ = balanceEthFloat.Float64()
	}

	addressInfo := AddressInfo{
		Address:    addressStr,
		Balance:    balance.String(),
		BalanceEth: fmt.Sprintf("%.6f", balanceEth),
		TxCount:    txCount,
		IsContract: isContract,
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    addressInfo,
	})
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	vars := mux.Vars(r)
	txHashStr := vars["hash"]

	txHash := common.HexToHash(txHashStr)

	// Получаем транзакцию
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		sendError(w, "Failed to get transaction", err)
		return
	}

	// Простая информация о транзакции
	txInfo := map[string]interface{}{
		"hash":     tx.Hash().Hex(),
		"pending":  isPending,
		"gas":      tx.Gas(),
		"gasPrice": tx.GasPrice().String(),
		"nonce":    tx.Nonce(),
		"value":    tx.Value().String(),
	}

	if !isPending {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			txInfo["blockNumber"] = receipt.BlockNumber.Uint64()
			txInfo["status"] = receipt.Status
			txInfo["gasUsed"] = receipt.GasUsed
		}
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    txInfo,
	})
}

func sendError(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   fmt.Sprintf("%s: %v", message, err),
	})
}
