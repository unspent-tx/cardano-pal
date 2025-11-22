package handlers

import (
	"context"
	"fmt"
	"go-cardano-address-safety/templates"
	"net/http"
	"os"

	"github.com/blockfrost/blockfrost-go"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}

func HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr1 := r.Form.Get("addr1")

	if addr1 == "" {
		http.Error(w, "address required", 400)
		return
	}

	// Create Blockfrost client
	key := os.Getenv("BLOCKFROST_KEY")
	api := blockfrost.NewAPIClient(blockfrost.APIClientOptions{
		ProjectID: key,
		Server:    blockfrost.CardanoPreProd,
	})
	ctx := context.Background()

	// Get address info
	fmt.Printf("Attempting to fetch address: %s\n", addr1)
	fmt.Printf("Using API key: %s...\n", key[:15]) // Print first 15 chars only

	resultCh := api.AddressTransactionsAll(ctx, addr1)
	var allTxs []blockfrost.AddressTransactions
	for result := range resultCh {
		if result.Err != nil {
			fmt.Printf("error: %v/n", result.Err)
			continue
		}
		for _, tx := range result.Res {
			fmt.Printf("Tx hash: %s, Block height: %d\n", tx.TxHash, tx.BlockHeight)
		}
		allTxs = append(allTxs, result.Res...)

	}

	result := templates.AnalysisResult{
		Addr: addr1,
		Info: allTxs,
	}

	templates.Results(result).Render(r.Context(), w)
}
