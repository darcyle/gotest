package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Purchase represents a gaming purchase in the system
type Purchase struct {
	ID                int64      `json:"id"`
	TransactionID     string     `json:"transaction_id"`
	PlayerID          string     `json:"player_id"`
	PlayerUsername    string     `json:"player_username"`
	GameTitle         string     `json:"game_title"`
	ItemType          string     `json:"item_type"`
	Genre             string     `json:"genre"`
	Platform          string     `json:"platform"`
	AmountCents       int        `json:"amount_cents"`
	Currency          string     `json:"currency"`
	PlayerLevel       int        `json:"player_level"`
	CreatedAt         time.Time  `json:"created_at"`
	Enriched          bool       `json:"enriched"`
}

// PlayerLoyalty represents a player's loyalty points
type PlayerLoyalty struct {
	PlayerID      string    `json:"player_id"`
	LoyaltyPoints int       `json:"loyalty_points"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PurchaseInput represents input data for purchase creation/update
type PurchaseInput struct {
	TransactionID     string   `json:"transaction_id"`
	PlayerID          string   `json:"player_id"`
	PlayerUsername    string   `json:"player_username"`
	GameTitle         string   `json:"game_title"`
	ItemType          string   `json:"item_type"`
	Genre             string   `json:"genre"`
	Platform          string   `json:"platform"`
	AmountCents       int      `json:"amount_cents"`
	Currency          string   `json:"currency"`
	PlayerLevel       int      `json:"player_level"`
	CreatedAt         string   `json:"created_at"` // ISO8601 string, will be parsed to time.Time
}

// UpsertResult contains the result of an upsert operation
type UpsertResult struct {
	Created bool `json:"created"`
	Updated bool `json:"updated"`
}

// Common errors
var (
	ErrBadInput      = errors.New("bad input")
	ErrNotFound      = errors.New("not found")
	ErrInvalidFormat = errors.New("invalid format")
)

// PurchaseStore defines the interface for purchase storage operations
type PurchaseStore interface {
	// AddPurchase inserts or updates a purchase by transaction_id
	AddPurchase(ctx context.Context, p Purchase) (created bool, err error)
	
	// ClaimBatchForEnrichment selects unenriched purchases using FOR UPDATE SKIP LOCKED
	// Returns up to 'batch' purchases that are locked for processing
	ClaimBatchForEnrichment(ctx context.Context, batch int) ([]Purchase, error)
	
	// MarkEnriched marks a purchase as enriched
	MarkEnriched(ctx context.Context, id int64) error
}

// AddPurchase implements PurchaseStore.AddPurchase
func (s *pgStore) AddPurchase(ctx context.Context, p Purchase) (bool, error) {
	// TODO: implement purchase
	
	return false, nil
}

// ClaimBatchForEnrichment implements PurchaseStore.ClaimBatchForEnrichment
func (s *pgStore) ClaimBatchForEnrichment(ctx context.Context, batch int) ([]Purchase, error) {
	// TODO: implement claim batch for enrichment

	return nil, errors.New("not implemented")
}

// MarkEnriched implements PurchaseStore.MarkEnriched
func (s *pgStore) MarkEnriched(ctx context.Context, id int64) error {
	// TODO: Mark purchase as enriched
	
	return errors.New("not implemented")
}

// Helper function to parse ISO8601 timestamp
func parseTimestamp(s string) (time.Time, error) {
	// TODO: Parse ISO8601 timestamp string to time.Time
	// Handle formats like "2025-08-15T10:00:00Z"
	return time.Parse(time.RFC3339, s)
}