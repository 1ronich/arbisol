package core

import (
	"arbisol/src/solana_functions"
	"arbisol/src/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	//"strings"
	"time"
)

func FetchAllRaydiumPools() {
	var result [][]types.DexLP
	seen := make(map[string]bool)   // to track seen LPs
	middle := make(map[string]bool) // to track seen LPs

	pageSize := 1000
	page := 1
outerLoop:
	for {

		resp := fetchRaydiumPools(strconv.Itoa(pageSize), strconv.Itoa(page))
		fmt.Println(page)

		fmt.Println(len(resp.Data.Data))
		if len(resp.Data.Data) == 0 || len(result) > 100 {
			fmt.Println(result)
			AppendToJSONFile(result, "./results/raydiumPubs.json")

			fmt.Println("DATA 0")
			break outerLoop
		}

		for _, lp1 := range resp.Data.Data {

			if (lp1.LpPrice * lp1.LpAmount) < 200000 {
				continue
			}

			if lp1.MintA.Address == lp1.MintB.Address || lp1.MintA.Address == types.WSOL || lp1.MintB.Address == types.WSOL {
				continue
			}

			middleKey := strings.TrimSpace(lp1.MintA.Name) + "/" + strings.TrimSpace(lp1.MintB.Name)
			if middle[middleKey] {
				continue // Skip if this combination has already been seen
			} else {
				middle[middleKey] = true // Mark this combination as seen
			}

			for _, lp2 := range resp.Data.Data {

				for _, lp3 := range resp.Data.Data {

					if lp1.MintA.Address == lp3.MintB.Address &&
						lp1.MintB.Address == lp2.MintA.Address &&
						lp2.MintB.Address == lp3.MintA.Address {

						/*if lp1.ProgramID != RAYDIUMCONCENTRATED && lp1.ProgramID != RAYDIUMCONSTANTMM && lp1.ProgramID != RAYDIUMV4POOL {
							continue
						}
						if lp2.ProgramID != RAYDIUMCONCENTRATED && lp2.ProgramID != RAYDIUMCONSTANTMM && lp2.ProgramID != RAYDIUMV4POOL {
							continue
						}
						if lp3.ProgramID != RAYDIUMCONCENTRATED && lp3.ProgramID != RAYDIUMCONSTANTMM && lp3.ProgramID != RAYDIUMV4POOL {
							continue
						}*/

						seenKey := strings.TrimSpace(lp1.MintA.Name) + "/" + strings.TrimSpace(lp1.MintB.Name) + "-" + strings.TrimSpace(lp2.MintA.Name) + "/" + strings.TrimSpace(lp2.MintB.Name) + "-" + strings.TrimSpace(lp3.MintA.Name) + "/" + strings.TrimSpace(lp3.MintB.Name)
						if seen[seenKey] {
							continue // Skip if this combination has already been seen
						} else {
							seen[seenKey] = true // Mark this combination as seen
						}
						fmt.Println("MATCH FOUND")
						fmt.Println("LP1 MintA:", lp1.MintA.Name, "LP1 MintB:", lp1.MintB.Name)

						result = append(result, []types.DexLP{
							{
								ProgramId: lp1.ProgramID,
								ID:        lp1.ID,
								MintA:     lp1.MintA.Address,
								//MintASymbol:    lp1.MintA.Symbol,
								//MintAName:      lp1.MintA.Name,
								TokenADecimals: lp1.MintA.Decimals,
								MintB:          lp1.MintB.Address,
								//MintBSymbol:    lp1.MintB.Symbol,
								//MintBName:      lp1.MintB.Name,
								TokenBDecimals: lp1.MintB.Decimals,
							}, {
								ProgramId: lp2.ProgramID,
								ID:        lp2.ID,
								MintA:     lp2.MintA.Address,
								//MintASymbol:    lp2.MintA.Symbol,
								//MintAName:      lp2.MintA.Name,
								TokenADecimals: lp2.MintA.Decimals,
								MintB:          lp2.MintB.Address,
								//MintBSymbol:    lp2.MintB.Symbol,
								//MintBName:      lp2.MintB.Name,
								TokenBDecimals: lp2.MintB.Decimals,
							}, {
								ProgramId: lp3.ProgramID,
								ID:        lp3.ID,
								MintA:     lp3.MintA.Address,
								//MintASymbol:    lp3.MintA.Symbol,
								//MintAName:      lp3.MintA.Name,
								TokenADecimals: lp3.MintA.Decimals,
								MintB:          lp3.MintB.Address,
								//MintBSymbol:    lp3.MintB.Symbol,
								//MintBName:      lp3.MintB.Name,
								TokenBDecimals: lp3.MintB.Decimals,
							},
						})
					}
				}
			}

		}
		page++
	}
	fmt.Println("finished big mf")
	fmt.Println(seen)
	fmt.Println(middle)
	fmt.Println(len(seen))
	fmt.Println(len(middle))

}

/*
	var result []types.RaydiumLP

	pageSize := 1000
	page := 1

outerLoop:

	for {
		resp := fetchRaydiumPools(strconv.Itoa(pageSize), strconv.Itoa(page))
		fmt.Println(page)

		for _, lp := range resp.Data.Data {

			if len(result) == 100 {
				AppendToJSONFile(result, "./results/raydiumPubs.json")
				fmt.Println(result)
				fmt.Println("Appended to file")
				break outerLoop
			}

			if len(resp.Data.Data) == 0 {
				break
			}
			if (lp.LpPrice*lp.LpAmount) < 500000 || (lp.LpPrice*lp.LpAmount) > 2000000 {
				continue // Skip to next pool if current pool is not in range
			}
			innerresult := types.RaydiumLP{
				ID:             lp.ID,
				MintA:          lp.MintA.Address,
				MintASymbol:    lp.MintA.Symbol,
				TokenADecimals: lp.MintA.Decimals,
				MintB:          lp.MintB.Address,
				MintBSymbol:    lp.MintB.Symbol,
				TokenBDecimals: lp.MintB.Decimals,
			}

			//if len(resp.Data.Data) < pageSize {
			//	innerresult = append(innerresult, lp.ID, lp.MintA.Address, lp.MintA.Symbol, lp.MintB.Address, lp.MintB.Symbol)
			//	result = append(result, innerresult)

			//	break // Last page
			//}

			result = append(result, innerresult)
		}

		time.Sleep(50 * time.Millisecond)
		page++
	}
	fmt.Println("Finished big guy")
*/
func fetchRaydiumPools(amount, page string) types.RaydiumAPI {
	apiurl := fmt.Sprintf("https://api-v3.raydium.io/pools/info/list?poolType=standard&poolSortField=liquidity&sortType=desc&pageSize=%s&page=%s", amount, page)
	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	// Read and print the response body for debugging
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var response types.RaydiumAPI

	// Unmarshal the JSON data into the Go struct
	erra := json.Unmarshal(responseBody, &response)
	if erra != nil {
		log.Fatalf("Error unmarshaling JSON: %v", erra)
	}

	//rspns, _ := json.MarshalIndent(response, "", "  ")
	//fmt.Println("Response Body:", string(rspns))

	return response
}

func FetchOrcaPools() []string {
	var pubkeys []string

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getProgramAccounts",
		"params": []interface{}{
			"whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc",
			map[string]interface{}{
				"encoding": "base64",
				"filters": []map[string]interface{}{
					{
						"dataSize": 653,
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

	var Response types.Response // FOR getAccountInfo

	erra := json.Unmarshal(body, &Response)
	if erra != nil {
		panic(erra)
	}

	for _, result := range Response.Result {
		pubkeys = append(pubkeys, result.Pubkey)
	}

	return pubkeys
}

func GetOrcaPools() {
	var orcaTokens [][]string
	LPs := FetchOrcaPools()

	result := solana_functions.FullWhirlpoolDeserialize(LPs)
	orcaTokens = append(orcaTokens, result)

	fmt.Println("Finished big loop")
	AppendToJSONFile(orcaTokens, "./pubs/orcaPubs.json")
}

// checks common LP addresses in the Orca list and the Raydium list
func Checker(list1 []string, list2 []string) []string {
	var result []string

	for listIndex1, orcaKey := range list1 {
		for listIndex2, raykey := range list2 {
			if orcaKey == raykey && list2[listIndex1] == list2[listIndex2] { //checks that both tokens in the LP match
				result = append(result, orcaKey, list2[listIndex1], raykey, list2[listIndex2]) //appends both pairs

			}
		}

	}

	return result
}
