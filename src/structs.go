package main

import (
	"math/big"
)

type AccountInfoResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Context struct {
			Slot       int    `json:"slot"`
			APIVersion string `json:"apiVersion"`
		} `json:"context"`
		Value struct {
			Lamports   uint64   `json:"lamports"`
			Data       []string `json:"data"` // Base64 data and encoding type
			Owner      string   `json:"owner"`
			Executable bool     `json:"executable"`
			RentEpoch  uint64   `json:"rentEpoch"`
			Space      int      `json:"space"`
		} `json:"value"`
	} `json:"result"`
}

type Whirlpool struct {
	WhirlpoolsConfig           []byte   // 32
	WhirlpoolBump              byte     // 1
	TickSpacing                uint16   // 2
	FeeTierIndexSeed           []byte   // 2
	FeeRate                    uint16   // 2
	ProtocolFeeRate            uint16   // 2
	Liquidity                  *big.Int // 16
	SqrtPrice                  *big.Int // 16
	TickCurrentIndex           int32    // 4
	ProtocolFeeOwedA           uint64   // 8
	ProtocolFeeOwedB           uint64   // 8
	TokenMintA                 []byte   // 32
	TokenVaultA                []byte   // 32
	FeeGrowthGlobalA           *big.Int // 16
	TokenMintB                 []byte   // 32
	TokenVaultB                []byte   // 32
	FeeGrowthGlobalB           *big.Int // 16
	RewardLastUpdatedTimestamp uint64   // 8
	//RewardInfos                []WhirlpoolRewardInfo // 384 total (assuming 3 x 128 bytes, needs actual layout to deserialize properly)
}

type LP struct {
	MintA string
	MintB string
}

type Ttoken struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	ChainID  int    `json:"chainId"`
	Decimals int    `json:"decimals"`
}
