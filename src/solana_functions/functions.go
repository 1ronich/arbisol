package solana_functions

//THIS FILE KEEPS ALL THE BASIC FUNCTIONS

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"

	solana_rpc "arbisol/src/solana_rpc"
	"arbisol/src/types"
)

//---------------------------------------------------RAYDIUM---------------------------------------------------\\

func BuildRaydiumSwapData(amountIn uint64, minAmountOut uint64) []byte {
	swap := types.SimpleSwapInstruction{
		Discriminator:    9,
		AmountIn:         amountIn,
		MinimumAmountOut: minAmountOut,
	}

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, swap.Discriminator); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, swap.AmountIn); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, swap.MinimumAmountOut); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}

	return buf.Bytes()
}

func BuildV2SwapInstructions(info types.SwapInfo) []solana.Instruction {
	raydiumProgramID := solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	instrData := BuildRaydiumSwapData(info.AmountIn, info.MinAmountOut)

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintA)
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintB)

	accounts := []*solana.AccountMeta{

		solana.NewAccountMeta(solana.TokenProgramID, false, false),
		solana.NewAccountMeta(info.Pool.AmmPool, true, false),
		solana.NewAccountMeta(info.Pool.AmmAuthority, false, false),

		solana.NewAccountMeta(info.Pool.AmmOpenOrders, true, false),

		solana.NewAccountMeta(info.Pool.AmmCoinVault, true, false),
		solana.NewAccountMeta(info.Pool.AmmPcVault, true, false),

		solana.NewAccountMeta(info.Pool.MarketProgram, false, false),
		solana.NewAccountMeta(info.Pool.Market, true, false),
		solana.NewAccountMeta(info.Pool.MarketBids, true, false),
		solana.NewAccountMeta(info.Pool.MarketAsks, true, false),
		solana.NewAccountMeta(info.Pool.MarketEventQueue, true, false),

		solana.NewAccountMeta(info.Pool.MarketCoinVault, true, false),
		solana.NewAccountMeta(info.Pool.MarketPcVault, true, false),

		solana.NewAccountMeta(info.Pool.MarketVaultSigner, false, false),

		solana.NewAccountMeta(userSourceTokenAcc, true, false),
		solana.NewAccountMeta(userDestTokenAcc, true, false),

		solana.NewAccountMeta(info.PayerPubkey, true, true),
	}
	instructions := []solana.Instruction{}
	instructions = append(instructions, solana.NewInstruction(raydiumProgramID, accounts, instrData))

	return instructions
}

func GetVaultSigner(marketID, programID solana.PublicKey, nonce uint64) solana.PublicKey {
	seed := [][]byte{
		marketID.Bytes(),
		Uint64ToLittleEndianBytes(nonce),
	}
	pda, err := solana.CreateProgramAddress(seed, programID)
	if err != nil {
		log.Fatalf("error creating vault signer PDA: %v", err)
	}
	return pda
}

func GetAMMAuthority(ammID solana.PublicKey) solana.PublicKey {
	seed := [][]byte{
		[]byte("amm authority"),
		ammID.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seed, ammID)
	if err != nil {
		log.Fatalf("error getting amm auth: %v\n", err)
	}

	return pda
}

func Uint64ToLittleEndianBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

func DecodeOpenBookData(data string) types.MarketState {
	bin, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalf("error decoding openbook data")
	}

	fmt.Println("B64 data:", data)

	var m types.MarketState
	buf := bytes.NewReader(bin[5:])

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &m); err != nil {
		panic(err)
	}
	return m
}

func GetMarketData(marketId solana.PublicKey) types.MarketState {
	response := solana_rpc.GetAccountInfo(marketId.String())

	res := DecodeOpenBookData(response.Result.Value.Data[0])

	return res
}

// ---------------------------------------------------ORCA---------------------------------------------------\\

//---------------------------------------------------GLOBAL---------------------------------------------------\\

func (s Brains) SendTransaction(total_Instructions []solana.Instruction) solana.Signature {
	recent := solana_rpc.GetLatestBlockhash()

	tx, err := solana.NewTransaction(
		total_Instructions,
		solana.MustHashFromBase58(recent.Result.Value.Blockhash),
		solana.TransactionPayer(s.Payer.PublicKey()),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(s.Payer.PublicKey()) {
				return &s.Payer
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("error signing the tx: %v", err)
	}

	client := rpc.New(os.Getenv("QUICKNODE"))
	res, err := client.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{
		PreflightCommitment: rpc.CommitmentConfirmed, // Match Solscan's commitment
	})
	if err != nil {
		log.Fatalf("could not send transaction successfully: %v", err)
	}

	return res
}

func SetComputeLimit(amount uint32) *computebudget.Instruction {
	computBudget := computebudget.NewSetComputeUnitLimitInstruction(
		amount,
	).Build()
	fmt.Printf("ProgramID: %s\n", computBudget.ProgramID())
	byts, _ := computBudget.Data()
	fmt.Printf("Data: %x\n", byts)

	return computBudget
}

func Transfer(info types.Transf) []solana.Instruction {
	transfInstr := []solana.Instruction{
		system.NewTransferInstruction(
			info.Amount,
			info.Sender.PublicKey(),
			info.Receiver,
		).Build(),
	}
	return transfInstr
}

func CreateNewATA(brain types.Atastruct, TokenMintAddress solana.PublicKey) solana.Instruction {
	//this function checks if the ATA that was not found is either "PoolSourcetoken" or "PoolDestToken" and then based on the ATA
	//it appends to 'poolToken' the address from the 'brain' variable, and then a CreateATAinstruction is created
	//var poolToken solana.PublicKey

	/*if PoolToken == "UserSourceTokenAcc" {
		poolToken = brain.UserSourceTokenAcc
		fmt.Println("FOund UserSourceTokenAcc:", brain.UserSourceTokenAcc)
	} else if PoolToken == "UserDestTokenAcc" {
		poolToken = brain.UserDestTokenAcc
		fmt.Println("FOund UserSourceTokenAcc:", brain.UserSourceTokenAcc)

	}*/

	ATA := associatedtokenaccount.NewCreateInstruction(
		brain.PayerPrivateKey.PublicKey(), //brain.PayerPrivateKey.PublicKey(),
		brain.PayerPrivateKey.PublicKey(), //brain.PayerPrivateKey.PublicKey(),
		TokenMintAddress,
	).Build()
	fmt.Printf("ATA Payer (signer & funder): %s\n", brain.PayerPrivateKey.PublicKey())
	fmt.Printf("ATA Owner: %s\n", brain.PayerPrivateKey.PublicKey())
	fmt.Printf("Token Mint: %s\n", TokenMintAddress)

	return ATA
}

func CheckForATAs(brain types.Atastruct) map[string]bool {
	//this function checks if there are any existing ATA accounts by calling the getAccountInfo method
	//and checking the length of the data returned; if length=0 --> no existing account, so I need to create an account
	//that means sending an instruction to the tx so it knows that it needs to create an ATA

	state := map[string]bool{
		"UserSourceTokenAcc": true,
		"UserDestTokenAcc":   true,
	}

	resp := solana_rpc.GetAccountInfo(brain.UserSourceTokenAcc.String())

	if resp.Result.Value.Data[0] == "" {
		state["UserSourceTokenAcc"] = false //ATA not existing
	}
	respa := solana_rpc.GetAccountInfo(brain.UserDestTokenAcc.String())
	if respa.Result.Value.Data[0] == "" {
		state["UserDestTokenAcc"] = false //ATA not existing
	}

	return state
}

func SimTransaction(total_Instructions []solana.Instruction, PayerPubkey solana.PublicKey) *rpc.SimulateTransactionResponse {

	recent := solana_rpc.GetLatestBlockhash()

	tx, err := solana.NewTransaction(
		total_Instructions,
		solana.MustHashFromBase58(recent.Result.Value.Blockhash),
		solana.TransactionPayer(PayerPubkey),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Sign
	/*_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(PayerPkey.PublicKey()) {
				return &PayerPkey
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("error signing the tx: %v", err)
	}*/

	client := rpc.New(os.Getenv("QUICKNODE"))
	resp, err := client.SimulateTransactionWithOpts(context.Background(), tx, &rpc.SimulateTransactionOpts{
		Commitment: rpc.CommitmentConfirmed, // Match Solscan's commitment
	})

	if err != nil {
		log.Fatal("err simulating tx:", err)
	}

	return resp

}

func GetTokenPrice(vault types.DexLP) float64 {
	var balance_responses []float64
	var account_balances []float64

	//var Mint string
	var Vault string

	for i := 0; i < 2; i++ {

		if vault.ProgramId == types.RAYDIUMCONCENTRATED {
			if i == 0 {
				Vault = DecodeCAMM(vault.ID).TokenVault0
			} else {
				Vault = DecodeCAMM(vault.ID).TokenVault1
			}
		} else if vault.ProgramId == types.RAYDIUMCONSTANTMM {
			if i == 0 {
				Vault = DecodeCPMM(vault.ID).Token0Vault
			} else {
				Vault = DecodeCPMM(vault.ID).Token1Vault
			}
		} else if vault.ProgramId == types.RAYDIUMV4POOL {
			if i == 0 {
				Vault = DecodeRayV4TokenAccount(vault.ID).EncodedInfo.BaseVault
			} else {
				Vault = DecodeRayV4TokenAccount(vault.ID).EncodedInfo.QuoteVault
			}
		} else if vault.ProgramId == types.WHIRLPOOLPROGRAMID {
			if i == 0 {
				Vault = DecodeWhirlpool(vault.ID).TokenVaultA.String()
			} else {
				Vault = DecodeWhirlpool(vault.ID).TokenVaultB.String()
			}
		}

		Response := solana_rpc.GetTokenAccountBalance(Vault)

		balanceFloat, _ := strconv.ParseFloat(Response.Result.Value.Amount, 64)
		balance_responses = append(balance_responses, balanceFloat)

	}

	account_balances = append(account_balances, balance_responses[0]/float64(vault.TokenADecimals))
	account_balances = append(account_balances, balance_responses[1]/float64(vault.TokenBDecimals))

	var price float64        //var to return
	var WhichIsStable string //to know which is the stable coin either VaultA(vault[0]) or VaultB(vault[1])

	switch vault.MintA { //A
	case types.USDC:
		//fmt.Println("MintA: USDC")
		WhichIsStable = "A"
	case types.USDT:
		//fmt.Println("MintA: USDT")
		WhichIsStable = "A"
	case types.WSOL:
		//fmt.Println("MintA: WSOL")
		WhichIsStable = "A"
	default:
		//fmt.Printf("MintA: %s\n", vaults.MintASymbol)
		WhichIsStable = "B"
	}

	switch vault.MintA { //B
	case types.USDC:
		//fmt.Println("MintB: USDC")
		WhichIsStable = "B"
	case types.USDT:
		//fmt.Println("MintB: USDT")
		WhichIsStable = "B"
	case types.WSOL:
		//fmt.Println("MintB: WSOL")
		WhichIsStable = "B"
	default:
		//fmt.Printf("MintB: %s\n", vaults.MintBSymbol)
		WhichIsStable = "A"
	}
	//fmt.Println("Which Is Stable:", WhichIsStable)

	if WhichIsStable == "A" { //price in stablecoin = AmountOfStableCoin / AmountOfToken
		price = account_balances[0] / account_balances[1] //in this case the stablecoin is the first vault
	} else if WhichIsStable == "B" {
		price = account_balances[1] / account_balances[0] //in this case the stablecoin is the second vault
	}

	return price
}

type Mint struct {
	MintAuthorityOption   uint32           // Whether mintAuthority is set
	MintAuthority         solana.PublicKey // PublicKey (32 bytes)
	Supply                uint64           // Total supply of tokens
	Decimals              uint8            // Number of decimals
	IsInitialized         bool             // Whether the mint is initialized
	FreezeAuthorityOption uint32           // Whether freezeAuthority is set
	FreezeAuthority       solana.PublicKey // PublicKey (32 bytes)
}

func GetTokenDecimal(token_mint string) uint8 {
	res := solana_rpc.GetAccountInfo(token_mint)

	if len(res.Result.Value.Data[0]) == 0 {
		fmt.Println("data not long enough")
		return 0
	}

	raw, err := base64.StdEncoding.DecodeString(res.Result.Value.Data[0])
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(raw)
	var data Mint

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		panic(err)
	}

	/*ress := Mint{
		MintAuthorityOption:   data.MintAuthorityOption,
		MintAuthority:         solana.PublicKeyFromBytes(data.MintAuthority[:]),
		Supply:                data.Supply,
		Decimals:              data.Decimals,
		IsInitialized:         data.IsInitialized,
		FreezeAuthorityOption: data.FreezeAuthorityOption,
		FreezeAuthority:       solana.PublicKeyFromBytes(data.FreezeAuthority[:]),
	}*/

	return data.Decimals
}
