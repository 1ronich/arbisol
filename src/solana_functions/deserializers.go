package solana_functions

import (
	solana_rpc "arbisol/src/solana_rpc"
	"arbisol/src/types"
	"arbisol/src/utils"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func DeserializeRayV4(data_b64 string) types.DeserializedRayV4 { //675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8
	raw, err := base64.StdEncoding.DecodeString(data_b64)
	if err != nil {
		panic(err)
	}

	var data types.PoolData
	buf := bytes.NewReader(raw)

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		panic(err)
	}

	baseVault := base58.Encode(data.BaseVault[:])
	quoteVault := base58.Encode(data.QuoteVault[:])
	baseMint := base58.Encode(data.BaseMint[:])
	quoteMint := base58.Encode(data.QuoteMint[:])
	lpMint := base58.Encode(data.LpMint[:])
	openOrders := base58.Encode(data.OpenOrders[:])
	marketId := base58.Encode(data.MarketId[:])
	marketProgramId := base58.Encode(data.MarketProgramId[:])
	targetOrders := base58.Encode(data.TargetOrders[:])
	withdrawQueue := base58.Encode(data.WithdrawQueue[:])
	lpVault := base58.Encode(data.LpVault[:])
	owner := base58.Encode(data.Owner[:])

	deserializedAddresses := types.DeserializedAddresses{
		BaseVault:       baseVault,
		QuoteVault:      quoteVault,
		BaseMint:        baseMint,
		QuoteMint:       quoteMint,
		LpMint:          lpMint,
		OpenOrders:      openOrders,
		MarketId:        marketId,
		MarketProgramId: marketProgramId,
		TargetOrders:    targetOrders,
		WithdrawQueue:   withdrawQueue,
		LpVault:         lpVault,
		Owner:           owner,
	}

	fullDeserializedData := types.DeserializedRayV4{
		Status:                 data.Status,
		Nonce:                  data.Nonce,
		MaxOrder:               data.MaxOrder,
		Depth:                  data.Depth,
		BaseDecimal:            data.BaseDecimal,
		QuoteDecimal:           data.QuoteDecimal,
		State:                  data.State,
		ResetFlag:              data.ResetFlag,
		MinSize:                data.MinSize,
		VolMaxCutRatio:         data.VolMaxCutRatio,
		AmountWaveRatio:        data.AmountWaveRatio,
		BaseLotSize:            data.BaseLotSize,
		QuoteLotSize:           data.QuoteLotSize,
		MinPriceMultiplier:     data.MinPriceMultiplier,
		MaxPriceMultiplier:     data.MaxPriceMultiplier,
		SystemDecimalValue:     data.SystemDecimalValue,
		MinSeparateNumerator:   data.MinSeparateNumerator,
		MinSeparateDenominator: data.MinSeparateDenominator,
		TradeFeeNumerator:      data.TradeFeeNumerator,
		TradeFeeDenominator:    data.TradeFeeDenominator,
		PnlNumerator:           data.PnlNumerator,
		PnlDenominator:         data.PnlDenominator,
		SwapFeeNumerator:       data.SwapFeeNumerator,
		SwapFeeDenominator:     data.SwapFeeDenominator,
		BaseNeedTakePnl:        data.BaseNeedTakePnl,
		QuoteNeedTakePnl:       data.QuoteNeedTakePnl,
		QuoteTotalPnl:          data.QuoteTotalPnl,
		BaseTotalPnl:           data.BaseTotalPnl,
		PoolOpenTime:           data.PoolOpenTime,
		PunishPcAmount:         data.PunishPcAmount,
		PunishCoinAmount:       data.PunishCoinAmount,
		OrderbookToInitTime:    data.OrderbookToInitTime,
		SwapBaseInAmount:       data.SwapBaseInAmount,
		SwapQuoteOutAmount:     data.SwapQuoteOutAmount,
		SwapBase2QuoteFee:      data.SwapBase2QuoteFee,
		SwapQuoteInAmount:      data.SwapQuoteInAmount,
		SwapBaseOutAmount:      data.SwapBaseOutAmount,
		SwapQuote2BaseFee:      data.SwapQuote2BaseFee,
		EncodedInfo:            deserializedAddresses,
		LpReserve:              data.LpReserve,
		Padding:                data.Padding,
	}

	return fullDeserializedData
}
func DecodeRayV4TokenAccount(account solana.PublicKey) types.DeserializedRayV4 {
	respa := solana_rpc.GetAccountInfo(account.String())
	ress := DeserializeRayV4(respa.Result.Value.Data[0])

	return ress
}

func DeserializeCPMMaccount(data_b64 string) types.DecodedCPMM { //CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C
	raw, err := base64.StdEncoding.DecodeString(data_b64)
	if err != nil {
		panic(err)
	}
	hexHExx := hex.EncodeToString(raw)
	fmt.Println("RAW:", data_b64)
	fmt.Println("hex:", hexHExx)

	buf := bytes.NewReader(raw)

	buf.Seek(8, io.SeekStart)
	AmmConfig := utils.ReadPubkey(buf)
	PoolCreator := utils.ReadPubkey(buf)
	Token0Vault := utils.ReadPubkey(buf)
	Token1Vault := utils.ReadPubkey(buf)
	LpMint := utils.ReadPubkey(buf)
	Token0Mint := utils.ReadPubkey(buf)
	Token1Mint := utils.ReadPubkey(buf)
	Token0Program := utils.ReadPubkey(buf)
	Token1Program := utils.ReadPubkey(buf)
	ObservationKey := utils.ReadPubkey(buf)
	//buf.Seek(3, io.SeekStart)
	//Mint0Decimals := utils.ReadU8(buf)
	//Mint1Decimals := utils.ReadU8(buf)

	return types.DecodedCPMM{
		AmmConfig:      AmmConfig,
		PoolCreator:    PoolCreator,
		Token0Vault:    Token0Vault,
		Token1Vault:    Token1Vault,
		LpMint:         LpMint,
		Token0Mint:     Token0Mint,
		Token1Mint:     Token1Mint,
		Token0Program:  Token0Program,
		Token1Program:  Token1Program,
		ObservationKey: ObservationKey,
		//Mint0Decimals:  Mint0Decimals,
		//Mint1Decimals:  Mint1Decimals,
	}
}
func DecodeCPMM(account solana.PublicKey) types.DecodedCPMM {
	respa := solana_rpc.GetAccountInfo(account.String())
	ress := DeserializeCPMMaccount(respa.Result.Value.Data[0])

	return ress
}

func DeserializeCAMM(data_b64 string) types.DecodedCAMM {
	raw, err := base64.StdEncoding.DecodeString(data_b64)
	if err != nil {
		panic(err)
	}

	var data types.DeserializedCAMM
	buf := bytes.NewReader(raw)

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		panic(err)
	}

	Bump := data.Bump
	AmmConfig := solana.PublicKeyFromBytes(data.AmmConfig[:]).String()
	Owner := solana.PublicKeyFromBytes(data.Owner[:]).String()
	TokenMint0 := solana.PublicKeyFromBytes(data.TokenMint0[:]).String()
	TokenMint1 := solana.PublicKeyFromBytes(data.TokenMint1[:]).String()
	TokenVault0 := solana.PublicKeyFromBytes(data.TokenVault0[:]).String()
	TokenVault1 := solana.PublicKeyFromBytes(data.TokenVault1[:]).String()
	ObservationKey := solana.PublicKeyFromBytes(data.ObservationKey[:]).String()

	MintDecimals0 := data.MintDecimals0
	MintDecimals1 := data.MintDecimals1
	TickSpacing := data.TickSpacing
	//Liquidity := new(big.Int).SetBytes(data.Liquidity[:]).String()
	//SqrtPriceX64 := new(big.Int).SetBytes(data.SqrtPriceX64[:]).String()
	//TickCurrent := data.TickCurrent
	//_ padding3
	//_ padding4
	//FeeGrowthGlobal0X64 := new(big.Int).SetBytes(data.FeeGrowthGlobal0X64[:]).String()
	//FeeGrowthGlobal1X64 := new(big.Int).SetBytes(data.FeeGrowthGlobal1X64[:]).String()
	//ProtocolFeesToken0 := data.ProtocolFeesToken0
	//ProtocolFeesToken1 := data.ProtocolFeesToken1
	//SwapInAmountToken0 := new(big.Int).SetBytes(data.SwapInAmountToken0[:]).String()
	//SwapOutAmountToken1 := new(big.Int).SetBytes(data.SwapOutAmountToken1[:]).String()
	//SwapInAmountToken1 := new(big.Int).SetBytes(data.SwapInAmountToken1[:]).String()
	//SwapOutAmountToken0 := new(big.Int).SetBytes(data.SwapOutAmountToken0[:]).String()
	//Status := data.Status
	//_ padding 7byte
	//var RewardInfo []string
	//TickArrayBitMap
	//TotalFeesToken0 := data.TotalFeesToken0
	//TotalFeesClaimedToken0 := data.TotalFeesClaimedToken0
	//TotalFeesToken1 := data.TotalFeesToken1
	//TotalFeesClaimedToken1 := data.TotalFeesClaimedToken1
	//FundFeesToken0 := data.FundFeesToken0
	//FundFeesToken1 := data.FundFeesToken1
	//OpenTime := data.OpenTime
	//RecentEpoch := data.RecentEpoch
	//_padding1 24byte
	//_padding2 32byte

	//for i, _ := range data.RewardInfos {
	//	rewardInfo := base58.Encode(data.RewardInfos[i].Mint[:])
	//	RewardInfo = append(RewardInfo, rewardInfo)
	//}

	//for i := 0; i < 16; i++ {
	//	b := make([]byte, 8)
	//	binary.LittleEndian.PutUint64(b, data.TickArrayBitmap[i])
	//	fmt.Printf("TickArrayBitmap[%d]: %s\n", i, base58.Encode(b))
	//}

	//for i := 0; i < 24; i++ {
	//	b := make([]byte, 8)
	//	binary.LittleEndian.PutUint64(b, data.Padding1[i])
	//	fmt.Printf("Padding1[%d]: %s\n", i, base58.Encode(b))
	//}

	return types.DecodedCAMM{
		Bump:           Bump,
		AmmConfig:      AmmConfig,
		Owner:          Owner,
		TokenMint0:     TokenMint0,
		TokenMint1:     TokenMint1,
		TokenVault0:    TokenVault0,
		TokenVault1:    TokenVault1,
		ObservationKey: ObservationKey,
		MintDecimals0:  MintDecimals0,
		MintDecimals1:  MintDecimals1,
		TickSpacing:    TickSpacing,
	}
}
func DecodeCAMM(account solana.PublicKey) types.DecodedCAMM {
	respa := solana_rpc.GetAccountInfo(account.String())
	ress := DeserializeCAMM(respa.Result.Value.Data[0])

	return ress
}

type WhirlpoolData struct {
	WhirlpoolsConfig           [32]byte // publicKey
	WhirlpoolBump              uint8    // u8,1
	TickSpacing                uint16   // u16
	TickSpacingSeed            [2]uint8 // u8,2 (array of 2 u8s)
	FeeRate                    uint16   // u16
	ProtocolFeeRate            uint16   // u16
	Liquidity                  [16]byte // u128
	SqrtPrice                  [16]byte // u128
	TickCurrentIndex           int32    // i32
	ProtocolFeeOwedA           uint64   // u64
	ProtocolFeeOwedB           uint64   // u64
	TokenMintA                 [32]byte // publicKey
	TokenVaultA                [32]byte // publicKey
	FeeGrowthGlobalA           [16]byte // u128
	TokenMintB                 [32]byte // publicKey
	TokenVaultB                [32]byte // publicKey
	FeeGrowthGlobalB           [16]byte // u128
	RewardLastUpdatedTimestamp uint64   // u64
	//RewardInfos                  [3]RewardInfo    // [object Object],3
}

type DecodedWhirlpoolData struct {
	WhirlpoolsConfig           solana.PublicKey // publicKey
	WhirlpoolBump              uint8            // u8,1
	TickSpacing                uint16           // u16
	TickSpacingSeed            [2]uint8         // u8,2 (array of 2 u8s)
	FeeRate                    uint16           // u16
	ProtocolFeeRate            uint16           // u16
	Liquidity                  [16]byte         // u128
	SqrtPrice                  [16]byte         // u128
	TickCurrentIndex           int32            // i32
	ProtocolFeeOwedA           uint64           // u64
	ProtocolFeeOwedB           uint64           // u64
	TokenMintA                 solana.PublicKey // publicKey
	TokenVaultA                solana.PublicKey // publicKey
	FeeGrowthGlobalA           [16]byte         // u128
	TokenMintB                 solana.PublicKey // publicKey
	TokenVaultB                solana.PublicKey // publicKey
	FeeGrowthGlobalB           [16]byte         // u128
	RewardLastUpdatedTimestamp uint64           // u64
}

func FullWhirlpoolDeserialize(pubKeys []string) []string {
	var pubkeys []string

	for _, pubKey := range pubKeys {
		res := solana_rpc.GetAccountInfo(pubKey)
		DeserializeWhirlpool(res.Result.Value.Data[0])

	}

	return pubkeys
}

func DeserializeWhirlpool(dataB64 string) DecodedWhirlpoolData {
	raw, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(raw[8:])
	var data WhirlpoolData

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		panic(err)
	}

	whirlpoolsConfig := solana.PublicKeyFromBytes(data.WhirlpoolsConfig[:])
	whirlpoolBump := data.WhirlpoolBump
	tickSpacing := data.TickSpacing
	tickSpacingSeed := data.TickSpacingSeed
	feeRate := data.FeeRate
	protocolFeeRate := data.ProtocolFeeRate
	liquidity := data.Liquidity
	sqrtPrice := data.SqrtPrice
	tickCurrentIndex := data.TickCurrentIndex
	protocolFeeOwedA := data.ProtocolFeeOwedA
	protocolFeeOwedB := data.ProtocolFeeOwedB
	tokenMintA := solana.PublicKeyFromBytes(data.TokenMintA[:]).String()
	tokenVaultA := solana.PublicKeyFromBytes(data.TokenVaultA[:]).String()
	feeGrowthGlobalA := data.FeeGrowthGlobalA
	tokenMintB := solana.PublicKeyFromBytes(data.TokenMintB[:]).String()
	tokenVaultB := solana.PublicKeyFromBytes(data.TokenVaultB[:]).String()
	feeGrowthGlobalB := data.FeeGrowthGlobalB
	rewardLastUpdatedTimestamp := data.RewardLastUpdatedTimestamp

	return DecodedWhirlpoolData{
		WhirlpoolsConfig:           whirlpoolsConfig,
		WhirlpoolBump:              whirlpoolBump,
		TickSpacing:                tickSpacing,
		TickSpacingSeed:            tickSpacingSeed,
		FeeRate:                    feeRate,
		ProtocolFeeRate:            protocolFeeRate,
		Liquidity:                  liquidity,
		SqrtPrice:                  sqrtPrice,
		TickCurrentIndex:           tickCurrentIndex,
		ProtocolFeeOwedA:           protocolFeeOwedA,
		ProtocolFeeOwedB:           protocolFeeOwedB,
		TokenMintA:                 solana.MustPublicKeyFromBase58(tokenMintA),
		TokenVaultA:                solana.MustPublicKeyFromBase58(tokenVaultA),
		FeeGrowthGlobalA:           feeGrowthGlobalA,
		TokenMintB:                 solana.MustPublicKeyFromBase58(tokenMintB),
		TokenVaultB:                solana.MustPublicKeyFromBase58(tokenVaultB),
		FeeGrowthGlobalB:           feeGrowthGlobalB,
		RewardLastUpdatedTimestamp: rewardLastUpdatedTimestamp,
	}
}

func DecodeWhirlpool(pubkey any) DecodedWhirlpoolData {
	res := solana_rpc.GetAccountInfo(pubkey)
	result := DeserializeWhirlpool(res.Result.Value.Data[0])
	return result
}
