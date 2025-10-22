package types

import "github.com/gagliardetto/solana-go"

type Uint128 struct {
	Low  uint64
	High uint64
}

type Pupsik struct {
	Keys [][]DexLP `json:"keys"`
}

type DexLP struct {
	ProgramId solana.PublicKey
	ID        solana.PublicKey `json:"id"`
	MintA     solana.PublicKey `json:"mintA"`
	//MintASymbol    string           `json:"mintASymbol"`
	//MintAName      string           `json:"mintAName"`
	TokenADecimals int              `json:"tokenADecimals"`
	MintB          solana.PublicKey `json:"mintB"`
	//MintBSymbol    string           `json:"mintBSymbol"`
	//MintBName      string           `json:"mintAName"`
	TokenBDecimals int `json:"tokenBDecimals"`
}

type Transf struct {
	Sender    solana.PrivateKey
	Receiver  solana.PublicKey
	Authority solana.PublicKey
	Amount    uint64
}

type SimpleSwapInstruction struct {
	Discriminator    uint8  // "u8" with value 9
	AmountIn         uint64 // "u64" with value 4487772414
	MinimumAmountOut uint64 // "u64" with value 1
}

type SwapInfo struct {
	PayerPkey   solana.PrivateKey
	PayerPubkey solana.PublicKey

	Pool PoolInfo

	AmountIn     uint64
	MinAmountOut uint64
	Slippage     uint64

	RPC string
}

type PoolInfo struct {
	Address solana.PublicKey
	LPMintA solana.PublicKey
	LPMintB solana.PublicKey

	AmmProgram    solana.PublicKey // AMM protocol program ID (e.g., Raydium)
	AmmPool       solana.PublicKey // AMM pool state account
	AmmAuthority  solana.PublicKey // AMM pool authority
	AmmOpenOrders solana.PublicKey // AMM pool open orders account (on Serum)
	AmmCoinVault  solana.PublicKey // AMM pool's base token vault
	AmmPcVault    solana.PublicKey // AMM pool's quote token vault

	MarketProgram     solana.PublicKey // Serum DEX program ID
	Market            solana.PublicKey // Serum market account
	MarketBids        solana.PublicKey // Serum market bids orderbook
	MarketAsks        solana.PublicKey // Serum market asks orderbook
	MarketEventQueue  solana.PublicKey // Serum market event queue
	MarketCoinVault   solana.PublicKey // Serum market's base token vault
	MarketPcVault     solana.PublicKey // Serum market's quote token vault
	MarketVaultSigner solana.PublicKey // Serum market vault authority

	UserTokenSource      solana.PublicKey // User's source token account
	UserTokenDestination solana.PublicKey // User's destination token account
	UserSourceOwner      solana.PublicKey // User's wallet (signer)
}

type Atastruct struct {
	PayerPrivateKey    solana.PrivateKey
	LPmints            [2]solana.PublicKey
	PoolSourceToken    solana.PublicKey
	PoolDestToken      solana.PublicKey
	UserSourceTokenAcc solana.PublicKey
	UserDestTokenAcc   solana.PublicKey
}

type DeserializedRayV4 struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       U128
	SwapQuoteOutAmount     U128
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      U128
	SwapBaseOutAmount      U128
	SwapQuote2BaseFee      uint64
	EncodedInfo            DeserializedAddresses
	LpReserve              uint64    `json:"lpReserve"`
	Padding                [3]uint64 `json:"padding"`
}

type DeserializedAddresses struct {
	BaseVault       string
	QuoteVault      string
	BaseMint        string
	QuoteMint       string
	LpMint          string
	OpenOrders      string
	MarketId        string
	MarketProgramId string
	TargetOrders    string
	WithdrawQueue   string
	LpVault         string
	Owner           string
}

// RPC METHOD STRUCTS
type SolanaAccountResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			ApiVersion string `json:"apiVersion"`
			Slot       uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Data       [2]string `json:"data"` // ["", "base64"]
			Executable bool      `json:"executable"`
			Lamports   uint64    `json:"lamports"`
			Owner      string    `json:"owner"`
			RentEpoch  uint64    `json:"rentEpoch"`
			Space      uint64    `json:"space"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}
type BlockhashResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Blockhash            string `json:"blockhash"`
			LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}
type GetBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			//ApiVersion string `json:"apiVersion"`
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Amount         string  `json:"amount"`
			Decimals       uint8   `json:"decimals"`
			UiAmount       float64 `json:"uiAmount"`
			UiAmountString string  `json:"uiAmountString"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

type GetProgramAccounts struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`

	Result struct {
		Pubkey  string `json:"pubkey"`
		Account struct {
			Lamports   uint64   `json:"lamports"`
			Owner      string   `json:"owner"`
			Data       []string `json:"data"` // Represents ["<encoded data>", "<encoding>"]
			Executable bool     `json:"executable"`
			RentEpoch  uint64   `json:"rentEpoch"`
			Space      uint64   `json:"space"`
		} `json:"account"`
	} `json:"result"`
}

//RAYDIUM STRUCTS

type U128 struct {
	Lo uint64
	Hi uint64
}

type PoolData struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64

	SwapBaseInAmount   U128
	SwapQuoteOutAmount U128
	SwapBase2QuoteFee  uint64
	SwapQuoteInAmount  U128
	SwapBaseOutAmount  U128
	SwapQuote2BaseFee  uint64
	BaseVault          [32]byte  `json:"baseVault"`
	QuoteVault         [32]byte  `json:"quoteVault"`
	BaseMint           [32]byte  `json:"baseMint"`
	QuoteMint          [32]byte  `json:"quoteMint"`
	LpMint             [32]byte  `json:"lpMint"`
	OpenOrders         [32]byte  `json:"openOrders"`
	MarketId           [32]byte  `json:"marketId"`
	MarketProgramId    [32]byte  `json:"marketProgramId"`
	TargetOrders       [32]byte  `json:"targetOrders"`
	WithdrawQueue      [32]byte  `json:"withdrawQueue"`
	LpVault            [32]byte  `json:"lpVault"`
	Owner              [32]byte  `json:"owner"`
	LpReserve          uint64    `json:"lpReserve"`
	Padding            [3]uint64 `json:"padding"`
}

type MarketState struct {
	AccountFlags           uint64
	OwnAddress             solana.PublicKey
	VaultSignerNonce       uint64
	CoinMint               solana.PublicKey
	PcMint                 solana.PublicKey
	CoinVault              solana.PublicKey
	CoinDepositsTotal      uint64
	CoinFeesAccrued        uint64
	PcVault                solana.PublicKey
	PcDepositsTotal        uint64
	PcFeesAccrued          uint64
	PcDustThreshold        uint64
	ReqQ                   solana.PublicKey
	EventQ                 solana.PublicKey
	Bids                   solana.PublicKey
	Asks                   solana.PublicKey
	CoinLotSize            uint64
	PcLotSize              uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
}

// LISTS structs
type Response struct {
	JSONRPC string   `json:"jsonrpc"`
	Result  []Result `json:"result"`
	ID      int      `json:"id"`
}
type Result struct {
	Pubkey  string  `json:"pubkey"`
	Account Account `json:"account"`
}
type Account struct {
	Data       []string `json:"data"`
	Executable bool     `json:"executable"`
	Lamports   uint64   `json:"lamports"`
	Owner      string   `json:"owner"`
	RentEpoch  uint64   `json:"rentEpoch"`
	Space      uint64   `json:"space"`
}

type RaydiumPriceAPI struct { //struct used to fetch the LPs from the Raydium API
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Data    []Pool `json:"data"`
}

type RaydiumAPI struct { //struct used to fetch the LPs from the Raydium API
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Data    struct {
		Count int    `json:"count"`
		Data  []Pool `json:"data"`
	} `json:"data"`
}
type Pool struct {
	//Type string `json:"type"`
	ProgramID   solana.PublicKey `json:"programId"`
	ID          solana.PublicKey `json:"id"`
	MintA       Token            `json:"mintA"`
	MintB       Token            `json:"mintB"`
	Price       float64          `json:"price"`
	MintAmountA float64          `json:"mintAmountA"`
	MintAmountB float64          `json:"mintAmountB"`
	//FeeRate     float64 `json:"feeRate"`
	//OpenTime    string  `json:"openTime"`
	//TVL         float64 `json:"tvl"`
	//Day                    Stats               `json:"day"`
	//Week                   Stats               `json:"week"`
	//Month                  Stats               `json:"month"`
	//PoolType []string `json:"pooltype"`
	//RewardDefaultPoolInfos string              `json:"rewardDefaultPoolInfos,omitempty"`
	//RewardDefaultInfos     []RewardDefaultInfo `json:"rewardDefaultInfos"`
	//FarmUpcomingCount int     `json:"farmUpcomingCount"`
	//FarmOngoingCount  int     `json:"farmOngoingCount"`
	//FarmFinishedCount int     `json:"farmFinishedCount"`
	//MarketID          string  `json:"marketId"`
	//LpMint            Token   `json:"lpMint"`
	LpPrice  float64 `json:"lpPrice"`
	LpAmount float64 `json:"lpAmount"`
	//BurnPercent       float64 `json:"burnPercent"`
	//LaunchMigratePool bool    `json:"launchMigratePool"`
}
type Token struct {
	//ChainID   int    `json:"chainId"`
	Address   solana.PublicKey `json:"address"`
	ProgramID solana.PublicKey `json:"programId"`
	//LogoURI    string            `json:"logoURI"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	//Tags       []string          `json:"tags"`
	//Extensions map[string]string `json:"extensions"`
}
type Stats struct {
	Volume      float64   `json:"volume"`
	VolumeQuote float64   `json:"volumeQuote"`
	VolumeFee   float64   `json:"volumeFee"`
	APR         float64   `json:"apr"`
	FeeAPR      float64   `json:"feeApr"`
	PriceMin    float64   `json:"priceMin"`
	PriceMax    float64   `json:"priceMax"`
	RewardAPR   []float64 `json:"rewardApr"`
}
type RewardDefaultInfo struct {
	Mint      Token  `json:"mint"`
	PerSecond string `json:"perSecond"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type RewardInfo struct {
	Mint               [32]byte // publicKey
	Vault              [32]byte // publicKey
	Authority          [32]byte // publicKey
	EmissionsPerSecond uint64   // u128
	GrowthGlobal       [16]byte // u128
	OpenTime           uint64   // u64
	LastUpdateTime     uint64   // u64
}

type SwapRequest struct {
	TokenAmountIn TokenAmount `json:"tokenAmountIn"`
	TokenOut      Token       `json:"tokenOut"`
	To            string      `json:"to"`
	From          string      `json:"from"`
	Slippage      int         `json:"slippage"` // in basis points (e.g., 50 = 0.5%)
}
type TokenAmount struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	ChainID  int    `json:"chainId"`
	Decimals int    `json:"decimals"`
	Amount   string `json:"amount"` // Using string to support large integers
}

type Vault struct {
	Vault    string
	Decimals int
	Symbol   string
}

type DeserializedCPMM struct {
	//Paddeng            uint64
	AmmConfig          [32]byte
	PoolCreator        [32]byte
	Token0Vault        [32]byte
	Token1Vault        [32]byte
	LpMint             [32]byte
	Token0Mint         [32]byte
	Token1Mint         [32]byte
	Token0Program      [32]byte
	Token1Program      [32]byte
	ObservationKey     [32]byte
	AuthBump           uint8
	Status             uint8
	LpMintDecimals     uint8
	Mint0Decimals      uint8
	Mint1Decimals      uint8
	_                  [3]byte // padding to align to 8 bytes (since 5 uint8 values before the next u64)
	LpSupply           uint64
	ProtocolFeesToken0 uint64
	ProtocolFeesToken1 uint64
	FundFeesToken0     uint64
	FundFeesToken1     uint64
	OpenTime           uint64
	Padding            [32]uint64
}
type DecodedCPMM struct {
	AmmConfig      string
	PoolCreator    string
	Token0Vault    string
	Token1Vault    string
	LpMint         string
	Token0Mint     string
	Token1Mint     string
	Token0Program  string
	Token1Program  string
	ObservationKey string

	Mint0Decimals uint8
	Mint1Decimals uint8
}

type CAMMRewardInfo struct {
	// You need to define this based on the actual structure of rewardInfos[0]
	// Here’s a placeholder example:
	Mint      solana.PublicKey
	Vault     solana.PublicKey
	Emissions uint64
	Padding   [8]byte
}

type DecodedCAMM struct {
	Bump           uint8
	AmmConfig      string
	Owner          string
	TokenMint0     string
	TokenMint1     string
	TokenVault0    string
	TokenVault1    string
	ObservationKey string
	MintDecimals0  uint8
	MintDecimals1  uint8
	TickSpacing    uint16
}

type DeserializedCAMM struct {
	Padding_00             [8]byte
	Bump                   uint8
	AmmConfig              [32]byte
	Owner                  [32]byte
	TokenMint0             [32]byte
	TokenMint1             [32]byte
	TokenVault0            [32]byte
	TokenVault1            [32]byte
	ObservationKey         [32]byte
	MintDecimals0          uint8
	MintDecimals1          uint8
	TickSpacing            uint16
	Liquidity              [16]byte // u128
	SqrtPriceX64           [16]byte // u128
	TickCurrent            int32
	Padding3               uint16
	Padding4               uint16
	FeeGrowthGlobal0X64    [16]byte // u128
	FeeGrowthGlobal1X64    [16]byte // u128
	ProtocolFeesToken0     uint64
	ProtocolFeesToken1     uint64
	SwapInAmountToken0     [16]byte // u128
	SwapOutAmountToken1    [16]byte // u128
	SwapInAmountToken1     [16]byte // u128
	SwapOutAmountToken0    [16]byte // u128
	Status                 uint8
	Padding                [7]byte
	RewardInfos            [3]CAMMRewardInfo
	TickArrayBitmap        [16]uint64
	TotalFeesToken0        uint64
	TotalFeesClaimedToken0 uint64
	TotalFeesToken1        uint64
	TotalFeesClaimedToken1 uint64
	FundFeesToken0         uint64
	FundFeesToken1         uint64
	OpenTime               uint64
	RecentEpoch            uint64
	Padding1               [24]uint64
	// padding2 not specified — add it if known
}
