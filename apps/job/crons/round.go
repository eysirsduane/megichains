package crons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetSolanaAirdrop() {
	rpcURL := "https://api.devnet.solana.com"
	address := "9gzkywhjpa3vP9XhvPVp1WtmnfYPdbYN7hKcfvUrmJqT" // base58 地址
	lamports := uint64(2_000_000_000)                         // 1 SOL

	// JSON-RPC 请求体
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "requestAirdrop",
		"params": []any{
			address,
			lamports,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		rpcURL,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}
