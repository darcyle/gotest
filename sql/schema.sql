-- GameVault Store: Database Schema
-- This schema supports gaming purchase ingestion and loyalty point enrichment

CREATE TABLE IF NOT EXISTS purchases (
  id                BIGSERIAL PRIMARY KEY,
  transaction_id    TEXT NOT NULL UNIQUE,
  player_id         TEXT NOT NULL,
  player_username   TEXT NOT NULL,
  game_title        TEXT NOT NULL,
  item_type         TEXT NOT NULL CHECK (item_type IN ('game', 'dlc', 'cosmetic', 'currency', 'season_pass')),
  genre             TEXT NOT NULL,
  platform          TEXT NOT NULL CHECK (platform IN ('steam', 'epic', 'xbox', 'playstation', 'nintendo', 'mobile')),
  amount_cents      INTEGER NOT NULL CHECK (amount_cents >= 0),
  currency          TEXT NOT NULL DEFAULT 'USD',
  player_level      INTEGER NOT NULL DEFAULT 1 CHECK (player_level >= 1),
  created_at        TIMESTAMPTZ NOT NULL,
  enriched          BOOLEAN NOT NULL DEFAULT FALSE,
  
  -- Add constraints for data integrity
  CONSTRAINT purchases_transaction_id_not_empty CHECK (length(transaction_id) > 0),
  CONSTRAINT purchases_player_id_not_empty CHECK (length(player_id) > 0),
  CONSTRAINT purchases_game_title_not_empty CHECK (length(game_title) > 0)
);

-- Player loyalty points table
CREATE TABLE IF NOT EXISTS player_loyalty (
  player_id         TEXT PRIMARY KEY,
  loyalty_points    INTEGER NOT NULL DEFAULT 0 CHECK (loyalty_points >= 0),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  -- Add constraint for data integrity
  CONSTRAINT player_loyalty_player_id_not_empty CHECK (length(player_id) > 0)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_purchases_id ON purchases(id);
CREATE INDEX IF NOT EXISTS idx_purchases_transaction_id ON purchases(transaction_id);
CREATE INDEX IF NOT EXISTS idx_purchases_player_id ON purchases(player_id);
CREATE INDEX IF NOT EXISTS idx_purchases_platform ON purchases(platform);
CREATE INDEX IF NOT EXISTS idx_purchases_enriched ON purchases(enriched) WHERE enriched = false;
CREATE INDEX IF NOT EXISTS idx_purchases_created_at ON purchases(created_at);
CREATE INDEX IF NOT EXISTS idx_purchases_game_title ON purchases(game_title);
CREATE INDEX IF NOT EXISTS idx_purchases_genre ON purchases(genre);

-- Composite indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_purchases_id_created_at ON purchases(id, created_at);
CREATE INDEX IF NOT EXISTS idx_purchases_player_platform ON purchases(player_id, platform);
CREATE INDEX IF NOT EXISTS idx_purchases_genre_amount ON purchases(genre, amount_cents);

-- Player loyalty indexes
CREATE INDEX IF NOT EXISTS idx_player_loyalty_updated_at ON player_loyalty(updated_at);