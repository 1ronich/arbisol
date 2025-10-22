package main

import (
	"arbisol/src/core"
	"arbisol/src/solana_functions"

	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/joho/godotenv"
)

var orcafile = "orcaPubs.json"
var rayfile = "raydiumPubs.json"

var rpcclient *rpc.Client
var wsclient *ws.Client

func init() {
	rpcclient = rpc.New(os.Getenv("QUICKNODE"))
	wsclient, _ = ws.Connect(context.Background(), os.Getenv("QUICKNODEWS"))
}

func main() {
	start := time.Now()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//getOrcaPools()
	//getCommonPublicKeys()
	//fetchAllRaydiumPools()
	//compareAll()

	//res := core.UnloadRaydiumList()
	//core.CalculateOpportunity(res)

	//core.FetchAllRaydiumPools()
	//core.GetOrcaPools()

	//res := solana_functions.BuildOrcaSwapData(SwapParams)
	//res := solana_functions.OrcaSwapInstruction(info)
	core := solana_functions.Brains{
		Payer:  solana.MustPrivateKeyFromBase58(os.Getenv("SOL_3enQ")),
		LpInfo: solana_functions.DecodeWhirlpool("FwewVm8u6tFPGewAyHmWAqad9hmF7mvqxK4mJ7iNqqGC"),
	}

	//tx := core.Orca_Swap(solana.MustPublicKeyFromBase58("FwewVm8u6tFPGewAyHmWAqad9hmF7mvqxK4mJ7iNqqGC"), 0.01, "So11111111111111111111111111111111111111112", "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB")
	//fmt.Println(tx)
	//solana_functions.GetTokenDecimal("6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN")
	//	res := solana_rpc.GetAccountInfo("6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN")
	//	fmt.Println(res.Result.Value.Data[0])

	//	bytes, _ := base64.StdEncoding.DecodeString(res.Result.Value.Data[0])

	//	fmt.Println(hex.EncodeToString(bytes))
	res := core.Orca_Sim_Swap(solana.MustPublicKeyFromBase58("FwewVm8u6tFPGewAyHmWAqad9hmF7mvqxK4mJ7iNqqGC"), 0.01, "So11111111111111111111111111111111111111112", "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB")
	fmt.Println(res)
	end := time.Since(start)
	fmt.Printf("Main unction took %s\n", end)
}

/*
func getCommonPublicKeys() { //gets common public keys of both .json lists

	file1 := orcafile
	file2 := rayfile

	filenames := []string{
		file1,
		file2,
	}

	list := core.UnloadAllJSONlists(filenames)
	if len(list) == 0 {
		fmt.Println("List length == 0")
	}

	fmt.Println("Starting checker")

	commonPubs := solana_functions.Checker(list[0][0], list[1][0])
	if len(commonPubs) == 0 {
		fmt.Println("Common pubs length == 0")
		return
	}

	//core.AppendToJSONFile(commonPubs, "./pubs/commonPubs.json")

}*/

func compare(raydiumLP, orcaLP LP) {
	//maxRayFeePercentage := 0.25
	//maxOrcaFeePercentage := 1

	//rayFee := (amount / 100) * maxRayFeePercentage
	//orcaFee := (amount / 100) * maxOrcaFeePercentage

	priceRaydium := core.GetPrice(raydiumLP.MintA)
	priceOrca := core.GetPrice(orcaLP.MintA)
	fmt.Printf("Raydium price: %f\n", priceRaydium)
	fmt.Printf("Orca price: %f\n\n", priceOrca)

	diff := math.Abs(priceRaydium - priceOrca)
	percentDiff := (diff / math.Min(priceRaydium, priceOrca)) * 100

	if percentDiff > 1.0 /*0.55*/ {

		if priceRaydium > priceOrca {
			fmt.Println("----Swapping token--------------") //swap()
		} else if priceOrca < priceRaydium {
			fmt.Println("----Swapping token--------------") //swap()
		}
	}
}

/*
func compareAll() {
	mag, _ := core.UnloadJSONList("./pubs/commonPubs.json")

	for i, lp := range mag.Keys {
		if i == len(mag.Keys)-3 {
			return
		}
		fmt.Println("I:", i)
		orcaLP := LP{
			MintA: lp,
			MintB: mag.Keys[i+1],
		}
		fmt.Printf("Orca:\n%s\n%s\n", orcaLP.MintA, orcaLP.MintB)
		rayLP := LP{
			MintA: mag.Keys[i+2],
			MintB: mag.Keys[i+3],
		}
		fmt.Printf("Raydium:\n%s\n%s\n", rayLP.MintA, rayLP.MintB)

		compare(orcaLP, rayLP)
	}
}
*/
