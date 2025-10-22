package solana_rpc

import (
	"arbisol/src/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gagliardetto/solana-go"
)

func GetLatestBlockhash() types.BlockhashResponse {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getLatestBlockhash",
		"params": map[string]interface{}{
			"commitment":     "processed",
			"minContextSlot": 1000,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", os.Getenv("QUICKNODE"), bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending http req: %v\n", err)
	}

	fmt.Println(res.Status)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var Response types.BlockhashResponse

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		log.Fatalf("Error umarshalling http req: %v\n", erra)
	}

	return Response
}

func GetTokenAccountBalance(Vault string) types.GetBalanceResponse {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountBalance",
		"params": []interface{}{
			Vault,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", os.Getenv("QUICKNODE"), bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending http req: %v\n", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var Response types.GetBalanceResponse

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		log.Fatalf("Error umarshalling http req: %v\n", erra)
	}

	return Response
}

func GetAccountInfo(public_key any) types.SolanaAccountResponse {
	var pubkey string

	//this if allows both strings and solana.PublicKeys to be passed as parameters
	if reflect.TypeOf(public_key).Kind() == reflect.String {
		pubkey = public_key.(string)
	} else if reflect.TypeOf(public_key).Name() == "PublicKey" {
		pubkey = public_key.(solana.PublicKey).String()
	}

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getAccountInfo",
		"params": []interface{}{
			pubkey,
			map[string]interface{}{
				"encoding":   "base64",
				"commitment": "confirmed",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	rpcurl := os.Getenv("QUICKNODE")

	req, err := http.NewRequest("POST", rpcurl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var Response types.SolanaAccountResponse // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}

func GetProgramAccounts(programID string, accountSize uint64) types.GetProgramAccounts { //gets all lps for x program id
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getProgramAccounts",
		"params": []interface{}{
			programID,
			map[string]interface{}{
				"filters": []interface{}{
					map[string]interface{}{
						"dataSize": accountSize,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	rpcurl := os.Getenv("QUICKNODE")

	req, err := http.NewRequest("POST", rpcurl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var Response types.GetProgramAccounts // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	return Response
}
