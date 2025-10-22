package types

import "github.com/gagliardetto/solana-go"

const (
	MAX_TICK_INDEX                         = 443636
	MIN_TICK_INDEX                         = -443636
	TICK_ARRAY_SIZE                        = 88
	FULL_RANGE_ONLY_TICK_SPACING_THRESHOLD = 32768
)

var (
	WHIRLPOOLPROGRAMID  = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")
	RAYDIUMV4POOL       = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	RAYDIUMCONCENTRATED = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	RAYDIUMCONSTANTMM   = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")

	WSOL = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	USDC = solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2q8jG4wGkFAb9ZmXDtZ6zDP9F1")
	USDT = solana.MustPublicKeyFromBase58("Es9vMFrzaCERbJxBJ1KzjE6tGk5R9FzTz5u5Q6F6X7wA")
)
