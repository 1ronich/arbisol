package core

import (
	"arbisol/src/solana_functions"
	"arbisol/src/types"
	"fmt"
	"time"
)

type Pair struct {
	ID string

	VaultA         string // address of the vault that holds tokenA
	TokenSymbolA   string
	TokenADecimals int //decimals of tokenA

	VaultB         string //address of the vault that holds tokenB
	TokenSymbolB   string
	TokenBDecimals int //decimals of tokenB

}

func CalculateOpportunity(triangles types.Pupsik) {

	for _, combination := range triangles.Keys {
		var combinationPrice []float64
		time.Sleep(500 * time.Millisecond)
		for _, pair := range combination {
			price := solana_functions.GetTokenPrice(types.DexLP{

				ID:    pair.ID,
				MintA: pair.MintA,
				//MintASymbol:    pair.MintASymbol,
				//MintAName:      pair.MintAName,
				TokenADecimals: pair.TokenADecimals,
				MintB:          pair.MintB,
				//MintBSymbol:    pair.MintBSymbol,
				//MintBName:      pair.MintBName,
				TokenBDecimals: pair.TokenBDecimals,
			})
			combinationPrice = append(combinationPrice, price)

			pricedif := (1 * combinationPrice[0]) * combinationPrice[1] * combinationPrice[2]
			if pricedif > 1.0 {
				fmt.Println("Opportunity found!")
				//fmt.Printf("Combination: %s/%s -> %s/%s -> %s/%s\n", combination[0].MintASymbol, combination[0].MintBSymbol, combination[1].MintASymbol, combination[1].MintBSymbol, combination[2].MintASymbol, combination[2].MintBSymbol)
			}
		}
	}
}
