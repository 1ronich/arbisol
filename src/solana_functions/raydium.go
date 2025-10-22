package solana_functions

import (
	solana_rpc "arbisol/src/solana_rpc"
	"arbisol/src/types"
	"context"
	"fmt"
	"log"
	"os"

	//"encoding/base64"

	"encoding/json"

	//"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	//"github.com/gagliardetto/solana-go/token"
)

//var rpcurl = os.Getenv("QUICKNODE")

func testo() {

	LPmint := "22WrmyTj8x2TRVQen3fxxi2r4Rn6JDHWoMTpsSmn8RUd"
	LPmintA := "ED5nyyWEzpPPiWimP8vYm7sD7TD3LAt3Q3gRTWHzPJBY"
	LPmintB := "So11111111111111111111111111111111111111112"

	privateKey, err := solana.PrivateKeyFromSolanaKeygenFile("../keypair.json")
	if err != nil {
		panic(err)
	}

	res := DecodeRayV4TokenAccount(solana.MustPublicKeyFromBase58(LPmint))
	swap := types.SwapInfo{
		PayerPkey:   privateKey,             //solana.MustPrivateKeyFromBase58(os.Getenv("MAIN_SOL_PKEY")),
		PayerPubkey: privateKey.PublicKey(), //solana.MustPrivateKeyFromBase58(os.Getenv("MAIN_SOL_PKEY")).PublicKey(),
		Pool: types.PoolInfo{
			Address: solana.MustPublicKeyFromBase58(LPmint), //pool address
			LPMintA: solana.MustPublicKeyFromBase58(LPmintA),
			LPMintB: solana.MustPublicKeyFromBase58(LPmintB),

			AmmProgram:    solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"),
			AmmPool:       solana.MustPublicKeyFromBase58(LPmint),
			AmmAuthority:  solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"), //getAMMAuthority(solana.MustPublicKeyFromBase58(LPmint)),
			AmmOpenOrders: solana.MustPublicKeyFromBase58(res.EncodedInfo.OpenOrders),
			AmmCoinVault:  solana.MustPublicKeyFromBase58(res.EncodedInfo.BaseVault),
			AmmPcVault:    solana.MustPublicKeyFromBase58(res.EncodedInfo.QuoteVault),

			MarketProgram:     solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketProgramId),
			Market:            solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId),
			MarketBids:        GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).Bids,
			MarketAsks:        GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).Asks,
			MarketEventQueue:  GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).EventQ,
			MarketCoinVault:   GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).CoinVault,
			MarketPcVault:     GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).PcVault,
			MarketVaultSigner: GetVaultSigner(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketProgramId), GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId)).VaultSignerNonce),
		},

		AmountIn:     uint64(5000000),
		MinAmountOut: uint64(1),
		RPC:          os.Getenv("QUICKNODE"),
	}

	fmt.Println("Starting to build instructions...")

	totalInstructions := Raydium_Swap(swap)

	/*
		for i, ix := range totalInstructions {
			fmt.Printf("Instruction #%d:\n", i+1)
			fmt.Printf("  ProgramID: %s\n", ix.ProgramID())
			fmt.Printf("  Accounts:\n")
			for j, acct := range ix.Accounts() {
				fmt.Printf("    %d: %s (Signer: %v, Writable: %v)\n", j, acct.PublicKey, acct.IsSigner, acct.IsWritable)
			}
			data, err := ix.Data()
			if err != nil {
				log.Fatalf("Could not get account data: %v\n", err)
			}
			fmt.Printf("  Data: %x\n", data)
			fmt.Println("-----")
		}
	*/

	ressa := SimTransaction(totalInstructions, swap.PayerPubkey)

	ressajson, _ := json.MarshalIndent(ressa, "  ", " ")
	fmt.Println(string(ressajson))

}

func Raydium_Swap(info types.SwapInfo) []solana.Instruction {

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintA) //a is wsol
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintB)

	brain := types.Atastruct{
		PayerPrivateKey:    info.PayerPkey,
		LPmints:            [2]solana.PublicKey{info.Pool.LPMintA, info.Pool.LPMintB},
		PoolSourceToken:    info.Pool.UserTokenSource,
		PoolDestToken:      info.Pool.UserTokenDestination,
		UserSourceTokenAcc: userSourceTokenAcc,
		UserDestTokenAcc:   userDestTokenAcc,
	}

	var total_Instructions []solana.Instruction

	compLimitInstr := SetComputeLimit(400_000)
	total_Instructions = append(total_Instructions, compLimitInstr)

	resMap := CheckForATAs(brain)
	//this for range loop checks if there are ATAs missing, and if it needs to create one, it will create one
	var index int8 = 0
	for _, value := range resMap {

		if !value {
			ATAinstr := CreateNewATA(brain, brain.LPmints[index])
			total_Instructions = append(total_Instructions, ATAinstr)
		} else {
			info.Pool.UserSourceOwner = userSourceTokenAcc
			info.Pool.UserTokenDestination = userDestTokenAcc
		}
		index++
	}

	swapInstr := BuildV2SwapInstructions(info)

	jitoTipTransferInstruction := Transfer(types.Transf{
		Sender:   info.PayerPkey,
		Receiver: solana.MustPublicKeyFromBase58("Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY"), //jitotip
	})

	for _, inner := range swapInstr {
		total_Instructions = append(total_Instructions, inner)
	}
	for _, inner := range jitoTipTransferInstruction {
		total_Instructions = append(total_Instructions, inner)
	}
	//these for range loops are necessary because you cannot pass a [][]solana.instruction to create a 'tx'

	return total_Instructions

}

func SendTransaction(total_Instructions []solana.Instruction, PayerPkey solana.PrivateKey) solana.Signature {
	recent := solana_rpc.GetLatestBlockhash()

	tx, err := solana.NewTransaction(
		total_Instructions,
		solana.MustHashFromBase58(recent.Result.Value.Blockhash),
		solana.TransactionPayer(PayerPkey.PublicKey()),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(PayerPkey.PublicKey()) {
				return &PayerPkey
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

type RPCResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  Resultt `json:"result"`
}

type Resultt struct {
	Context Context `json:"context"`
	Value   Value   `json:"value"`
}

type Context struct {
	Slot int `json:"slot"`
}

type Value struct {
	Err        interface{} `json:"err"`
	Logs       []string    `json:"logs"`
	Accounts   interface{} `json:"accounts"`
	ReturnData ReturnData  `json:"returnData"`
}

type ReturnData struct {
	ProgramID            string   `json:"programId"`
	Data                 []string `json:"data"` // Base64-encoded data + format
	ComputeUnitsConsumed int      `json:"computeUnitsConsumed"`
}
