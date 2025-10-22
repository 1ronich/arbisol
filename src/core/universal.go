package core

import (
	"arbisol/src/types"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func UnloadRaydiumList() types.Pupsik {
	//var ray [][]string

	rayList, _ := UnloadJSONList(fmt.Sprintf("./results/%s", "raydiumPubs.json"))

	return rayList
}

func AppendToJSONFile(content any, filename string) {
	// Convert []string to map[string]string

	data := map[string]interface{}{
		"keys": content,
	}

	updatedData, _ := json.MarshalIndent(data, "", "  ")

	// Step 4: Write back to file (overwrite)
	_ = os.WriteFile(filename, updatedData, 0644)

}

func UnloadJSONList(filename string) (types.Pupsik, error) {

	file, err := os.Open(filename) //fmt.Sprintf("./pubs/%s", filename))
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var person types.Pupsik
	err = json.Unmarshal(bytes, &person)
	if err != nil {
		log.Fatalf("error unmarshaling raydium JSON: %v", err)
	}

	return person, nil
}

func GetPrice(token_mint string) float64 {

	url := fmt.Sprintf("https://solana-gateway.moralis.io/token/mainnet/%s/price", token_mint)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", os.Getenv("MORALISAPI"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var Response TokenInfo
	err := json.Unmarshal(body, &Response)
	if err != nil {
		panic(err)
	}

	//jsonResp, _ := json.MarshalIndent(Response, " ", "")
	//fmt.Println(string(jsonResp))

	return Response.UsdPrice
}

type TokenInfo struct {
	TokenAddress              string          `json:"tokenAddress"`
	PairAddress               string          `json:"pairAddress"`
	ExchangeName              string          `json:"exchangeName"`
	ExchangeAddress           string          `json:"exchangeAddress"`
	NativePrice               NativePriceInfo `json:"nativePrice"`
	UsdPrice                  float64         `json:"usdPrice"`
	UsdPrice24h               float64         `json:"usdPrice24h"`
	UsdPrice24hrUsdChange     float64         `json:"usdPrice24hrUsdChange"`
	UsdPrice24hrPercentChange float64         `json:"usdPrice24hrPercentChange"`
	Logo                      string          `json:"logo"`
	Name                      string          `json:"name"`
	Symbol                    string          `json:"symbol"`
	IsVerifiedContract        bool            `json:"isVerifiedContract"`
}
type NativePriceInfo struct {
	Value    string `json:"value"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
}
