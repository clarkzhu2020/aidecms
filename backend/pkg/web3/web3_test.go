package web3

import (
	"context"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	manager := GetManager()

	// Test manager initialization
	if manager == nil {
		t.Fatal("manager should not be nil")
	}

	// Test supported chains initially empty
	chains := manager.GetSupportedChains()
	if len(chains) != 0 {
		t.Errorf("expected 0 chains, got %d", len(chains))
	}
}

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		chain   Chain
		address string
		valid   bool
	}{
		// Ethereum
		{Ethereum, "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb", false}, // missing last character
		{Ethereum, "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", true},
		{Ethereum, "742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", false}, // missing 0x

		// BSC (same format as Ethereum)
		{BSC, "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", true},

		// Bitcoin
		{Bitcoin, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", true},
		{Bitcoin, "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq", true},
		{Bitcoin, "invalid", false},

		// Solana
		{Solana, "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK", true},
		{Solana, "invalid", false},
	}

	for _, tt := range tests {
		err := ValidateAddress(tt.chain, tt.address)
		if tt.valid && err != nil {
			t.Errorf("ValidateAddress(%s, %s) expected valid, got error: %v", tt.chain, tt.address, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateAddress(%s, %s) expected invalid, got no error", tt.chain, tt.address)
		}
	}
}

func TestValidateTxHash(t *testing.T) {
	tests := []struct {
		chain  Chain
		txHash string
		valid  bool
	}{
		// Ethereum
		{Ethereum, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", true},
		{Ethereum, "0x1234", false},
		{Ethereum, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", false}, // missing 0x

		// Bitcoin (64 hex characters)
		{Bitcoin, "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098", true},
		{Bitcoin, "invalid", false},
	}

	for _, tt := range tests {
		err := ValidateTxHash(tt.chain, tt.txHash)
		if tt.valid && err != nil {
			t.Errorf("ValidateTxHash(%s, %s) expected valid, got error: %v", tt.chain, tt.txHash, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateTxHash(%s, %s) expected invalid, got no error", tt.chain, tt.txHash)
		}
	}
}

func TestConfig(t *testing.T) {
	cfg := LoadConfig()

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	// Test that config has default values
	if cfg.EthereumRPC == "" {
		t.Error("EthereumRPC should have a default value")
	}

	if cfg.SolanaRPC == "" {
		t.Error("SolanaRPC should have a default value")
	}
}

func TestMultiChainAddress(t *testing.T) {
	// This is a unit test that doesn't require actual blockchain connections
	addr := MultiChainAddress{
		Bitcoin:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		Ethereum: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		BSC:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		Solana:   "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK",
	}

	// Validate all addresses
	if err := ValidateAddress(Bitcoin, addr.Bitcoin); err != nil {
		t.Errorf("Bitcoin address validation failed: %v", err)
	}

	if err := ValidateAddress(Ethereum, addr.Ethereum); err != nil {
		t.Errorf("Ethereum address validation failed: %v", err)
	}

	if err := ValidateAddress(BSC, addr.BSC); err != nil {
		t.Errorf("BSC address validation failed: %v", err)
	}

	if err := ValidateAddress(Solana, addr.Solana); err != nil {
		t.Errorf("Solana address validation failed: %v", err)
	}
}

func TestTransaction(t *testing.T) {
	tx := &Transaction{
		Hash:        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		From:        "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		To:          "0x1234567890abcdef1234567890abcdef12345678",
		Value:       "1000000000000000000",
		BlockNumber: 12345678,
		BlockHash:   "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		Status:      "success",
		GasUsed:     21000,
		GasPrice:    "20000000000",
		Timestamp:   time.Now().Unix(),
	}

	if tx.Hash == "" {
		t.Error("transaction hash should not be empty")
	}

	if tx.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", tx.Status)
	}
}

func TestWalletInfo(t *testing.T) {
	info := &WalletInfo{
		Address: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		Chain:   Ethereum,
		Balance: "1000000000000000000",
		Nonce:   5,
		TxCount: 5,
		Extra:   make(map[string]interface{}),
	}

	if info.Address == "" {
		t.Error("wallet address should not be empty")
	}

	if info.Chain != Ethereum {
		t.Errorf("expected chain 'ethereum', got '%s'", info.Chain)
	}
}

// Integration tests (require actual RPC connections)
func TestEthereumClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Use a public testnet RPC
	client, err := NewEthereumClient("https://rpc.ankr.com/eth")
	if err != nil {
		t.Fatalf("failed to create ethereum client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test GetBlockNumber
	blockNumber, err := client.GetBlockNumber(ctx)
	if err != nil {
		t.Errorf("GetBlockNumber failed: %v", err)
	} else {
		t.Logf("Current block number: %d", blockNumber)
	}

	// Test GetChainID
	chainID, err := client.GetChainID(ctx)
	if err != nil {
		t.Errorf("GetChainID failed: %v", err)
	} else {
		t.Logf("Chain ID: %s", chainID.String())
	}

	// Test GetGasPrice
	gasPrice, err := client.GetGasPrice(ctx)
	if err != nil {
		t.Errorf("GetGasPrice failed: %v", err)
	} else {
		t.Logf("Gas price: %s wei", gasPrice)
	}
}

func TestSolanaClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Use a public Solana RPC
	client := NewSolanaClient("https://api.mainnet-beta.solana.com")
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test GetBlockNumber
	slot, err := client.GetBlockNumber(ctx)
	if err != nil {
		t.Errorf("GetBlockNumber failed: %v", err)
	} else {
		t.Logf("Current slot: %d", slot)
	}

	// Test GetVersion
	version, err := client.GetVersion(ctx)
	if err != nil {
		t.Errorf("GetVersion failed: %v", err)
	} else {
		t.Logf("Solana version: %v", version)
	}
}
