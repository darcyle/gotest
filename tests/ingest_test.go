package tests

import (
	"context"
	"strings"
	"testing"

	main "gaming-purchases-system"
)


// TestStreamNDJSON tests NDJSON parsing functionality
func TestStreamNDJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		wantCount int
		wantErr   bool
	}{
		{
			name: "valid ndjson",
			jsonData: `{"transaction_id":"TXN-001","player_id":"player_001","player_username":"GamerAlice","game_title":"Cyberpunk 2077","item_type":"game","genre":"RPG","platform":"steam","amount_cents":5999,"currency":"USD","player_level":15,"created_at":"2025-08-15T10:00:00Z"}
{"transaction_id":"TXN-002","player_id":"player_002","player_username":"BobTheBuilder","game_title":"Minecraft","item_type":"game","genre":"Sandbox","platform":"epic","amount_cents":2699,"currency":"USD","player_level":8,"created_at":"2025-08-15T11:00:00Z"}`,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "empty ndjson",
			jsonData:  "",
			wantCount: 0,
			wantErr:   false,
		},
		// TODO: Add more test cases
		// - Invalid JSON
		// - Missing fields
		// - Invalid timestamps
		// - Empty lines (should be skipped)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.jsonData)
			
			var purchases []main.Purchase
			err := main.StreamNDJSON(context.Background(), reader, func(purchase main.Purchase) error {
				purchases = append(purchases, purchase)
				return nil
			})
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if len(purchases) != tt.wantCount {
				t.Errorf("Got %d purchases, want %d", len(purchases), tt.wantCount)
			}
			
			// TODO: Validate purchase contents
		})
	}
}


// TestValidatePurchaseInput tests input validation
func TestValidatePurchaseInput(t *testing.T) {
	tests := []struct {
		name    string
		input   main.PurchaseInput
		wantErr bool
	}{
		{
			name: "valid purchase",
			input: main.PurchaseInput{
				TransactionID:  "TXN-001",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "game",
				Genre:          "RPG",
				Platform:       "steam",
				AmountCents:    5999,
				Currency:       "USD",
				PlayerLevel:    15,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "empty transaction_id",
			input: main.PurchaseInput{
				TransactionID:  "",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "game",
				Genre:          "RPG",
				Platform:       "steam",
				AmountCents:    5999,
				Currency:       "USD",
				PlayerLevel:    15,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "invalid item_type",
			input: main.PurchaseInput{
				TransactionID:  "TXN-001",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "invalid",
				Genre:          "RPG",
				Platform:       "steam",
				AmountCents:    5999,
				Currency:       "USD",
				PlayerLevel:    15,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "invalid platform",
			input: main.PurchaseInput{
				TransactionID:  "TXN-001",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "game",
				Genre:          "RPG",
				Platform:       "invalid",
				AmountCents:    5999,
				Currency:       "USD",
				PlayerLevel:    15,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			input: main.PurchaseInput{
				TransactionID:  "TXN-001",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "game",
				Genre:          "RPG",
				Platform:       "steam",
				AmountCents:    -100,
				Currency:       "USD",
				PlayerLevel:    15,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "invalid player_level",
			input: main.PurchaseInput{
				TransactionID:  "TXN-001",
				PlayerID:       "player_001",
				PlayerUsername: "GamerAlice",
				GameTitle:      "Cyberpunk 2077",
				ItemType:       "game",
				Genre:          "RPG",
				Platform:       "steam",
				AmountCents:    5999,
				Currency:       "USD",
				PlayerLevel:    0,
				CreatedAt:      "2025-08-15T10:00:00Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := main.ValidatePurchaseInput(tt.input)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

