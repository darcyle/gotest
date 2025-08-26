package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type PurchaseRecord struct {
	TransactionID     string `json:"transaction_id"`
	PlayerID          string `json:"player_id"`
	PlayerUsername    string `json:"player_username"`
	GameTitle         string `json:"game_title"`
	ItemType          string `json:"item_type"`
	Genre             string `json:"genre"`
	Platform          string `json:"platform"`
	AmountCents       int    `json:"amount_cents"`
	Currency          string `json:"currency"`
	PlayerLevel       int    `json:"player_level"`
	CreatedAt         string `json:"created_at"`
}

var (
	games = []string{
		"The Witcher 3", "Cyberpunk 2077", "Rocket League", "Minecraft",
		"Fortnite", "Call of Duty", "FIFA 24", "Spider-Man Miles Morales",
		"Zelda BOTW", "Halo Infinite", "Valorant", "Counter-Strike 2",
		"Apex Legends", "Overwatch 2", "Grand Theft Auto V", "Red Dead Redemption 2",
	}
	
	itemTypes = []string{"game", "dlc", "cosmetic", "currency", "season_pass"}
	
	genres = []string{"RPG", "Action", "Sports", "Sandbox", "FPS", "Adventure", "Strategy", "Racing"}
	
	platforms = []string{"steam", "epic", "xbox", "playstation", "nintendo", "mobile"}
	
	currencies = []string{"USD", "EUR", "GBP", "CAD", "AUD"}
	
	// Pool of player IDs to reuse for loyalty calculations
	playerPool = make([]string, 0, 200)
	usernamePool = make([]string, 0, 200)
)

func init() {
	// Generate a pool of players to reuse
	platformPrefixes := map[string]string{
		"steam": "steam_76561198",
		"epic": "epic_player_",
		"xbox": "xbox_user_",
		"playstation": "psn_",
		"nintendo": "nintendo_",
		"mobile": "mobile_",
	}
	
	baseUsernames := []string{
		"Gamer", "Player", "Pro", "Master", "Ninja", "Legend", "Elite", "Champion",
		"Warrior", "Hero", "Shadow", "Storm", "Fire", "Ice", "Dragon", "Phoenix",
	}
	
	for i := 0; i < 200; i++ {
		platform := platforms[rand.Intn(len(platforms))]
		prefix := platformPrefixes[platform]
		playerID := fmt.Sprintf("%s%09d", prefix, rand.Intn(1000000000))
		username := fmt.Sprintf("%s%s%d", baseUsernames[rand.Intn(len(baseUsernames))], platform, rand.Intn(9999))
		
		playerPool = append(playerPool, playerID)
		usernamePool = append(usernamePool, username)
	}
}

func generatePurchase(transactionCounter int, baseTime time.Time) PurchaseRecord {
	// 70% chance to reuse existing player, 30% chance for new player
	var playerID, username string
	if rand.Float64() < 0.7 && len(playerPool) > 0 {
		idx := rand.Intn(len(playerPool))
		playerID = playerPool[idx]
		username = usernamePool[idx]
	} else {
		// Generate new player
		platform := platforms[rand.Intn(len(platforms))]
		platformPrefixes := map[string]string{
			"steam": "steam_76561198",
			"epic": "epic_player_",
			"xbox": "xbox_user_",
			"playstation": "psn_",
			"nintendo": "nintendo_",
			"mobile": "mobile_",
		}
		prefix := platformPrefixes[platform]
		playerID = fmt.Sprintf("%s%09d", prefix, rand.Intn(1000000000))
		username = fmt.Sprintf("NewPlayer%d", transactionCounter)
	}
	
	game := games[rand.Intn(len(games))]
	itemType := itemTypes[rand.Intn(len(itemTypes))]
	
	// Adjust game title based on item type
	if itemType == "dlc" {
		game += " DLC Pack"
	} else if itemType == "cosmetic" {
		game += " Skin Pack"
	} else if itemType == "currency" {
		game += " Gems"
	} else if itemType == "season_pass" {
		game += " Season Pass"
	}
	
	// Price varies by item type
	var amountCents int
	switch itemType {
	case "game":
		amountCents = rand.Intn(5000) + 1999 // $19.99 - $69.99
	case "dlc":
		amountCents = rand.Intn(2000) + 499  // $4.99 - $24.99
	case "cosmetic":
		amountCents = rand.Intn(1500) + 199  // $1.99 - $16.99
	case "currency":
		amountCents = rand.Intn(5000) + 99   // $0.99 - $50.99
	case "season_pass":
		amountCents = rand.Intn(2000) + 999  // $9.99 - $29.99
	}
	
	// Random timestamp within a range
	timeDelta := time.Duration(rand.Intn(86400*30)) * time.Second // Within 30 days
	timestamp := baseTime.Add(-timeDelta)
	
	return PurchaseRecord{
		TransactionID:  fmt.Sprintf("TXN-GEN-%06d", transactionCounter),
		PlayerID:       playerID,
		PlayerUsername: username,
		GameTitle:      game,
		ItemType:       itemType,
		Genre:          genres[rand.Intn(len(genres))],
		Platform:       platforms[rand.Intn(len(platforms))],
		AmountCents:    amountCents,
		Currency:       currencies[rand.Intn(len(currencies))],
		PlayerLevel:    rand.Intn(99) + 1, // Level 1-100
		CreatedAt:      timestamp.Format(time.RFC3339),
	}
}

func main() {
	var (
		count = flag.Int("count", 1000, "Number of records to generate")
		output = flag.String("output", "", "Output file path (default: data/test-timestamp.ndjson)")
	)
	flag.Parse()
	
	rand.Seed(time.Now().UnixNano())
	
	// Generate output filename if not specified
	if *output == "" {
		timestamp := time.Now().Format("20060102-150405")
		*output = fmt.Sprintf("data/test-%s.ndjson", timestamp)
	}
	
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	file, err := os.Create(*output)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()
	
	baseTime := time.Now()
	duplicateChance := 0.05 // 5% chance of duplicate transaction_id for upsert testing
	
	log.Printf("Generating %d purchase records...", *count)
	
	for i := 0; i < *count; i++ {
		purchase := generatePurchase(i+1, baseTime)
		
		// Occasionally create duplicates for upsert testing
		if i > 0 && rand.Float64() < duplicateChance {
			// Reuse a previous transaction ID but update some fields
			previousTxnID := fmt.Sprintf("TXN-GEN-%06d", rand.Intn(i)+1)
			purchase.TransactionID = previousTxnID
			purchase.PlayerLevel += rand.Intn(3) + 1 // Level up
			purchase.AmountCents += rand.Intn(1000)   // Price increase
		}
		
		data, err := json.Marshal(purchase)
		if err != nil {
			log.Fatalf("Failed to marshal purchase %d: %v", i+1, err)
		}
		
		if _, err := file.Write(data); err != nil {
			log.Fatalf("Failed to write purchase %d: %v", i+1, err)
		}
		
		if _, err := file.WriteString("\n"); err != nil {
			log.Fatalf("Failed to write newline after purchase %d: %v", i+1, err)
		}
		
		if (i+1)%1000 == 0 {
			log.Printf("Generated %d records...", i+1)
		}
	}
	
	log.Printf("Successfully generated %d records in %s", *count, *output)
	log.Printf("Player pool reuse: ~70%% of records use %d existing players for loyalty testing", len(playerPool))
}