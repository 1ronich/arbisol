package solana_functions

import (
	"arbisol/src/solana_rpc"
	"arbisol/src/types"
	"arbisol/src/utils"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type OrcaSwap struct {
	PayerPkey   solana.PrivateKey
	PayerPubkey solana.PublicKey

	LPmint solana.PublicKey

	MintA solana.PublicKey
	MintB solana.PublicKey

	//Pool PoolInfo
	Accounts   OrcaSwapAccounts
	SwapParams SwapParams
	//	Ticks      Tick_array_addresses

	AmountIn uint64
	//MinAmountOut uint64
	RateAforB float64
	Slippage  uint64
	AToB      bool
}

type Brains struct {
	Payer  solana.PrivateKey
	LpInfo DecodedWhirlpoolData
}

type FullSwapData struct {
	LP_Account_address solana.PublicKey

	Compute_Limit uint64
	Swap_Params   SwapParams
}

func (s Brains) Orca_Sim_Swap(LP_address solana.PublicKey, amount float64, token_A, token_B string) *rpc.SimulateTransactionResponse {
	var TokenA string
	//var TokenB string
	var A_to_B bool

	//LpInfo := DecodeWhirlpool(LP_address)

	tokenPrice := GetTokenPrice(types.DexLP{
		ProgramId: types.WHIRLPOOLPROGRAMID,
		ID:        LP_address,
		MintA:     s.LpInfo.TokenMintA,
		MintB:     s.LpInfo.TokenMintB,
	})
	sqrtPrice := math.Sqrt(tokenPrice * 0.90)

	if token_A == s.LpInfo.TokenMintA.String() {
		TokenA = token_A
		//TokenB = token_B
		A_to_B = true
	} else if token_A == s.LpInfo.TokenMintB.String() {
		TokenA = token_B
		//TokenB = token_A
		A_to_B = false
	} else {
		fmt.Println("booooooooomboclat")
	}

	fmt.Println("A_to_B", A_to_B)

	amount_in := amount * float64(GetTokenDecimal(TokenA))
	min_amount_out := float64(amount_in) * 0.95
	Swap := FullSwapData{
		LP_Account_address: LP_address,
		Compute_Limit:      400_000,
		Swap_Params: SwapParams{
			Amount:                 uint64(amount_in),
			OtherAmountThreshold:   uint64(min_amount_out),
			SqrtPriceLimit:         big.NewInt(int64(sqrtPrice)),
			AmountSpecifiedIsInput: true,
			AToB:                   A_to_B,
		},
	}

	res := s.FullInstructionsBuild(Swap)
	tx := SimTransaction(res, s.Payer.PublicKey())

	return tx
}

func (s Brains) Orca_Swap(LP_address solana.PublicKey, amount float64, token_A, token_B string) solana.Signature {
	var TokenA string
	//var TokenB string
	var A_to_B bool

	//LpInfo := DecodeWhirlpool(LP_address)

	tokenPrice := GetTokenPrice(types.DexLP{
		ProgramId: types.WHIRLPOOLPROGRAMID,
		ID:        LP_address,
		MintA:     s.LpInfo.TokenMintA,
		MintB:     s.LpInfo.TokenMintB,
	})
	sqrtPrice := math.Sqrt(tokenPrice * 0.90)

	if token_A == s.LpInfo.TokenMintA.String() {
		TokenA = token_A
		//TokenB = token_B
		A_to_B = true
	} else if token_A == s.LpInfo.TokenMintB.String() {
		TokenA = token_B
		//TokenB = token_A
		A_to_B = false
	} else {
		fmt.Println("booooooooomboclat")
	}

	amount_in := amount * float64(GetTokenDecimal(TokenA))
	min_amount_out := float64(amount_in) * 0.95
	Swap := FullSwapData{
		LP_Account_address: LP_address,
		Compute_Limit:      400_000,
		Swap_Params: SwapParams{
			Amount:                 uint64(amount_in),
			OtherAmountThreshold:   uint64(min_amount_out),
			SqrtPriceLimit:         big.NewInt(int64(sqrtPrice)),
			AmountSpecifiedIsInput: true,
			AToB:                   A_to_B,
		},
	}

	res := s.FullInstructionsBuild(Swap)
	tx := s.SendTransaction(res)

	return tx
}

func (s Brains) FullInstructionsBuild(Swap FullSwapData) []solana.Instruction {
	//var Is_A_to_B bool
	var instructions []solana.Instruction

	//lpInfo := DecodeWhirlpool(Swap.LP_Account_address)

	//---CompLimit---
	fmt.Println("COMPUTE LIMTI INSTR:----------------")

	computeLimInstr := SetComputeLimit(uint32(Swap.Compute_Limit))
	instructions = append(instructions, computeLimInstr)

	/////
	//---ATAs---
	fmt.Println("\nCREATE NEW ATA:----------------")

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(s.Payer.PublicKey(), s.LpInfo.TokenMintA) //a is wsol
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(s.Payer.PublicKey(), s.LpInfo.TokenMintB)
	brain := types.Atastruct{
		PayerPrivateKey:    s.Payer,
		LPmints:            [2]solana.PublicKey{s.LpInfo.TokenMintA, s.LpInfo.TokenMintB},
		UserSourceTokenAcc: userSourceTokenAcc,
		UserDestTokenAcc:   userDestTokenAcc,
	}

	ATAmap := CheckForATAs(brain)
	var index int8 = 0
	for _, value := range ATAmap {
		if !value {
			ATAinstr := CreateNewATA(brain, brain.LPmints[index])

			instructions = append(instructions, ATAinstr)
		}
		index++
	}
	/////
	//---Ticks---
	fmt.Println("\nCREATE ACCOUNTS:----------------\n")

	tickAddresses, PDAinstructions := s.GetTickArrayPDAs(Swap.LP_Account_address.String(), s.LpInfo, 3, Swap.Swap_Params.AToB)

	if len(PDAinstructions) > 0 {
		for _, instr := range PDAinstructions {
			instructions = append(instructions, instr)
		}
	}

	/////
	//---Swap---
	fmt.Println("\nSWAP:----------------\n")

	swapInstr := s.OrcaSwapInstruction(OrcaSwap{
		MintA: s.LpInfo.TokenMintA,
		MintB: s.LpInfo.TokenMintB,

		Accounts: OrcaSwapAccounts{
			TokenAuthority: solana.MustPublicKeyFromBase58(s.Payer.PublicKey().String()),
			Whirlpool:      solana.MustPublicKeyFromBase58(Swap.LP_Account_address.String()),
			TokenVaultA:    s.LpInfo.TokenVaultA,
			TokenVaultB:    s.LpInfo.TokenVaultB,
			TickArray0:     tickAddresses[0],
			TickArray1:     tickAddresses[1],
			TickArray2:     tickAddresses[2],
			Oracle:         getOracle(Swap.LP_Account_address.String()),
		},
		SwapParams: Swap.Swap_Params,
	}, Swap.Swap_Params.AToB)

	instructions = append(instructions, swapInstr...)

	return instructions

}

func GetStartTickIndex(tickCurrentIndex int, tickSpacing int, offset int32) int32 {
	TICK_ARRAY_SIZE := 88

	currentArrayIndex := ((tickCurrentIndex) / (tickSpacing) * (TICK_ARRAY_SIZE)) * (TICK_ARRAY_SIZE)

	//Reverse(int32(Index1), "FwewVm8u6tFPGewAyHmWAqad9hmF7mvqxK4mJ7iNqqGC")

	return (int32(currentArrayIndex) + (offset * int32(tickSpacing) * int32(TICK_ARRAY_SIZE)))
}

func (s Brains) GetTickArrayPDAs(LpMint string, lpInfo DecodedWhirlpoolData, numOfTickArrays int, a_To_b bool) ([]solana.PublicKey, []solana.Instruction) {
	var res []solana.PublicKey
	var PDAinstructions []solana.Instruction

	currentTick := lpInfo.TickCurrentIndex
	tickSpacing := lpInfo.TickSpacing

	for i := 0; i < numOfTickArrays; i++ {
		var offset int32
		if a_To_b {
			offset = int32(-1)
		} else {
			offset = int32(1)
		}

		fmt.Println("Arraty index list:", lpInfo.TickCurrentIndex)

		startTick := GetStartTickIndex(int(currentTick), int(tickSpacing), offset)

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, startTick) // int32 = 4 bytes

		seed := [][]byte{
			[]byte("tick_array"),
			solana.MustPublicKeyFromBase58(LpMint).Bytes(),
			buf.Bytes(),
		}

		tickAddress, _, err := solana.FindProgramAddress(seed, solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"))
		if err != nil {
			log.Fatalf("error getting PDA: %v", err)
		}

		if !AccountExists(tickAddress.String()) {
			instruction := initTickArray(startTick, LpMint, tickAddress, s.Payer.PublicKey())
			PDAinstructions = append(PDAinstructions, instruction...)
		}

		res = append(res, tickAddress)
	}

	return res, PDAinstructions
}

func AccountExists(account string) bool {
	if len(solana_rpc.GetAccountInfo(account).Result.Value.Data[0]) < 0 {
		return false
	} else {
		return true
	}
}

func NegativeRound(num float64) float64 {
	if num < 0 {
		return math.Floor(num)
	}
	return num
}

type Tick_data struct {
	Tick_Spacing       uint16
	Current_tick_index int32
	Start_index1       int32
	Start_index2       int32
	Start_index3       int32
}

func getTickAddresses(tickData Tick_data, LpMint string) []solana.PublicKey {
	var tickAddresses []solana.PublicKey

	startIndexes := []int32{tickData.Start_index1, tickData.Start_index2, tickData.Start_index3}
	for _, startIndex := range startIndexes {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, startIndex) // int32 = 4 bytes

		oracleSeed := [][]byte{
			[]byte("tick_array"),
			solana.MustPublicKeyFromBase58(LpMint).Bytes(),
			buf.Bytes(),
		}

		oraclePDA, _, err := solana.FindProgramAddress(oracleSeed, solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"))
		if err != nil {
			log.Fatalf("error getting PDA Oracle: %v", err)
		}

		tickAddresses = append(tickAddresses, oraclePDA)
	}
	return tickAddresses
}

func initTickArray(startIndex int32, LPmint string, tickArray, funder solana.PublicKey) []solana.Instruction {
	hexa, _ := hex.DecodeString("0bbcc1d68d5b95b8") // orca tick array denominator
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, startIndex)
	binary.Write(data, binary.LittleEndian, hexa)

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(LPmint), true, false),
		solana.NewAccountMeta(funder, true, true),
		solana.NewAccountMeta(tickArray, true, false),
		solana.NewAccountMeta(solana.SystemProgramID, false, false),
	}
	instruction := []solana.Instruction{}

	fmt.Printf("LP mint address: %s\n", LPmint)
	fmt.Printf("Funder (payer/signer): %s\n", funder)
	fmt.Printf("Tick array account: %s\n", tickArray)
	fmt.Printf("System program ID: %s\n", solana.SystemProgramID)

	instruction = append(instruction, solana.NewInstruction(types.WHIRLPOOLPROGRAMID, accounts, data.Bytes()))
	return instruction
}

func getOracle(LPmint string) solana.PublicKey {
	oracleSeed := [][]byte{
		[]byte("oracle"),
		solana.MustPublicKeyFromBase58(LPmint).Bytes(), // poolAddress is a solana.PublicKey
	}

	oraclePDA, _, err := solana.FindProgramAddress(oracleSeed, types.WHIRLPOOLPROGRAMID)
	if err != nil {
		log.Fatalf("error getting PDA Oracle: %v", err)
	}
	return oraclePDA
}

func calculatePriceLimit(rate float64, slippage uint64) *big.Int {
	priceLim := rate * ((100 - float64(slippage)) / 100) //expected amount out with slippage
	sqrtPriceLim := math.Sqrt(priceLim)
	power := math.Pow(priceLim, 64)
	res := power * sqrtPriceLim
	result, _ := new(big.Float).SetFloat64(res).Int(nil)

	return result
}

type OrcaSwapAccounts struct {
	TokenProgram   solana.PublicKey // Must match token::ID
	TokenAuthority solana.PublicKey // Signer
	Whirlpool      solana.PublicKey // Account<Whirlpool>, mut

	TokenOwnerAccountA solana.PublicKey // mut, must match whirlpool.token_mint_a
	TokenVaultA        solana.PublicKey // mut, must match whirlpool.token_vault_a

	TokenOwnerAccountB solana.PublicKey // mut, must match whirlpool.token_mint_b
	TokenVaultB        solana.PublicKey // mut, must match whirlpool.token_vault_b

	TickArrays021 []solana.PublicKey

	TickArray0 solana.PublicKey // UncheckedAccount, mut
	TickArray1 solana.PublicKey // UncheckedAccount, mut
	TickArray2 solana.PublicKey // UncheckedAccount, mut

	Oracle solana.PublicKey // UncheckedAccount, derived with seeds: ["oracle", whirlpool.key()]
}

func GetTokenOwnerAccount_A_B(whirlpool string) []solana.PublicKey {

	fmt.Println(whirlpool)
	res := DecodeWhirlpool(whirlpool)

	return []solana.PublicKey{(res.TokenMintA), (res.TokenMintB)}

}

func (s Brains) OrcaSwapInstruction(info OrcaSwap, Is_A_to_B bool) []solana.Instruction {

	instrData := BuildOrcaSwapData(info.SwapParams)

	var mintA solana.PublicKey
	var mintB solana.PublicKey

	if Is_A_to_B {
		mintA = info.MintA
		mintB = info.MintB
	} else if !Is_A_to_B {
		mintA = info.MintB
		mintB = info.MintA
	}

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(s.Payer.PublicKey(), mintA)
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(s.Payer.PublicKey(), mintB)

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(solana.TokenProgramID, false, false),
		solana.NewAccountMeta(info.Accounts.TokenAuthority, true, true),
		solana.NewAccountMeta(info.Accounts.Whirlpool, true, false),

		solana.NewAccountMeta(userSourceTokenAcc, true, false),
		solana.NewAccountMeta(info.Accounts.TokenVaultA, true, false),

		solana.NewAccountMeta(userDestTokenAcc, true, false),
		solana.NewAccountMeta(info.Accounts.TokenVaultB, true, false),

		solana.NewAccountMeta(info.Accounts.TickArray0, true, false),
		solana.NewAccountMeta(info.Accounts.TickArray1, true, false),
		solana.NewAccountMeta(info.Accounts.TickArray2, true, false),
		solana.NewAccountMeta(info.Accounts.Oracle, false, false),
	}
	fmt.Printf("Token Program ID: %s\n", solana.TokenProgramID)
	fmt.Printf("Authority (signer): %s\n", info.Accounts.TokenAuthority)
	fmt.Printf("Whirlpool account: %s\n", info.Accounts.Whirlpool)

	fmt.Printf("User source token account: %s\n", userSourceTokenAcc)
	fmt.Printf("Vault A (pool side A token vault): %s\n", info.Accounts.TokenVaultA)

	fmt.Printf("User destination token account: %s\n", userDestTokenAcc)
	fmt.Printf("Vault B (pool side B token vault): %s\n", info.Accounts.TokenVaultB)

	fmt.Printf("Tick Array 0: %s\n", info.Accounts.TickArray0)
	fmt.Printf("Tick Array 1: %s\n", info.Accounts.TickArray1)
	fmt.Printf("Tick Array 2: %s\n", info.Accounts.TickArray2)

	fmt.Printf("Oracle account: %s\n", info.Accounts.Oracle)

	instruction := []solana.Instruction{}
	instruction = append(instruction, solana.NewInstruction(types.WHIRLPOOLPROGRAMID, accounts, instrData))
	return instruction
}

type SwapParams struct {
	Amount                 uint64   `json:"amount"`               // u64
	OtherAmountThreshold   uint64   `json:"otherAmountThreshold"` // u64
	SqrtPriceLimit         *big.Int `json:"sqrtPriceLimit"`       // u128
	AmountSpecifiedIsInput bool     `json:"amountSpecifiedIsInput"`
	AToB                   bool     `json:"aToB"`
}

func BuildOrcaSwapData(params SwapParams) []byte {
	whirlpoolDiscriminator, _ := hex.DecodeString("f8c69e91e17587c8")
	var buf bytes.Buffer

	// Serialize uint64 fields
	if err := binary.Write(&buf, binary.LittleEndian, whirlpoolDiscriminator); err != nil {
		log.Fatalf("error binary write while building orca swap data: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, params.Amount); err != nil {
		log.Fatalf("error binary write while building orca swap data: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, params.OtherAmountThreshold); err != nil {
		log.Fatalf("error binary write while building orca swap data: %v", err)
	}

	// Serialize *big.Int (SqrtPriceLimit)
	// We want to serialize it as a 128-bit (16 bytes) value.
	sqrtPriceLimitBytes := utils.BigIntToLEBytesRightPad(params.SqrtPriceLimit, 16)
	// Ensure it is exactly 16 bytes (zero-padded if necessary)
	if len(sqrtPriceLimitBytes) < 16 {
		padding := make([]byte, 16-len(sqrtPriceLimitBytes))
		sqrtPriceLimitBytes = append(padding, sqrtPriceLimitBytes...)
	}
	buf.Write(sqrtPriceLimitBytes)

	// Serialize bool fields
	if err := binary.Write(&buf, binary.LittleEndian, params.AmountSpecifiedIsInput); err != nil {
		log.Fatalf("error binary write while building orca swap data: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, params.AToB); err != nil {
		log.Fatalf("error binary write while building orca swap data: %v", err)
	}

	return buf.Bytes()
}
