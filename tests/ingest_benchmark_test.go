package tests

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	main "gaming-purchases-system"
)

// BenchmarkStreamNDJSON tests parsing throughput
func BenchmarkStreamNDJSON(b *testing.B) {
	// Create sample NDJSON data
	sampleRecord := `{"transaction_id":"TXN-BENCH-%d","player_id":"player_bench_%d","player_username":"BenchPlayer%d","game_title":"Benchmark Game","item_type":"game","genre":"Action","platform":"steam","amount_cents":2999,"currency":"USD","player_level":25,"created_at":"2025-08-15T10:00:00Z"}`
	
	tests := []struct {
		name  string
		lines int
	}{
		{"100_records", 100},
		{"1000_records", 1000},
		{"10000_records", 10000},
	}
	
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Pre-generate test data
			var jsonData strings.Builder
			for i := 0; i < tt.lines; i++ {
				jsonData.WriteString(fmt.Sprintf(sampleRecord+"\n", i, i, i))
			}
			testData := jsonData.String()
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				reader := strings.NewReader(testData)
				processedCount := 0
				
				err := main.StreamNDJSON(context.Background(), reader, func(purchase main.Purchase) error {
					processedCount++
					return nil
				})
				
				if err != nil {
					b.Fatalf("StreamNDJSON failed: %v", err)
				}
				
				if processedCount != tt.lines {
					b.Fatalf("Expected %d records, got %d", tt.lines, processedCount)
				}
			}
			
			// Report throughput
			recordsPerOp := float64(tt.lines)
			b.ReportMetric(recordsPerOp, "records/op")
		})
	}
}

// BenchmarkIngestLargeFile tests memory efficiency with large files
func BenchmarkIngestLargeFile(b *testing.B) {
	// Generate a larger dataset for memory testing
	recordCount := 50000
	sampleRecord := `{"transaction_id":"TXN-LARGE-%d","player_id":"player_%d","player_username":"Player%d","game_title":"Large Test Game","item_type":"game","genre":"RPG","platform":"epic","amount_cents":4999,"currency":"USD","player_level":42,"created_at":"2025-08-15T15:30:00Z"}`
	
	// Pre-generate test data
	var jsonData strings.Builder
	for i := 0; i < recordCount; i++ {
		jsonData.WriteString(fmt.Sprintf(sampleRecord+"\n", i, i, i))
	}
	testData := jsonData.String()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testData)
		processedCount := 0
		
		err := main.StreamNDJSON(context.Background(), reader, func(purchase main.Purchase) error {
			processedCount++
			// Simulate minimal processing to test parsing overhead
			_ = purchase.TransactionID
			return nil
		})
		
		if err != nil {
			b.Fatalf("StreamNDJSON failed: %v", err)
		}
		
		if processedCount != recordCount {
			b.Fatalf("Expected %d records, got %d", recordCount, processedCount)
		}
	}
	
	// Calculate throughput metrics
	recordsPerOp := float64(recordCount)
	b.ReportMetric(recordsPerOp, "records/op")
	
	// Calculate approximate data size processed
	dataSizeMB := float64(len(testData)) / (1024 * 1024)
	b.ReportMetric(dataSizeMB, "MB/op")
}

// BenchmarkGeneratedLargeFile tests with realistic large files using generate-data.go
func BenchmarkGeneratedLargeFile(b *testing.B) {
	tests := []struct {
		name  string
		count int
	}{
		{"generated_5k", 5000},
		{"generated_25k", 25000},
		{"generated_100k", 100000},
	}
	
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Generate test file using our data generator
			tempDir := b.TempDir()
			testFile := filepath.Join(tempDir, fmt.Sprintf("bench-%d.ndjson", tt.count))
			
			// Run generate-data.go to create realistic test data
			cmd := exec.Command("go", "run", "../generate-data.go", 
				"-count", fmt.Sprintf("%d", tt.count),
				"-output", testFile)
			cmd.Dir = ".."
			
			output, err := cmd.CombinedOutput()
			if err != nil {
				b.Fatalf("Failed to generate test data: %v\nOutput: %s", err, output)
			}
			
			// Verify file was created
			if _, err := os.Stat(testFile); os.IsNotExist(err) {
				b.Fatalf("Test file was not created: %s", testFile)
			}
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				file, err := os.Open(testFile)
				if err != nil {
					b.Fatalf("Failed to open test file: %v", err)
				}
				
				processedCount := 0
				err = main.StreamNDJSON(context.Background(), file, func(purchase main.Purchase) error {
					processedCount++
					// Simulate realistic processing
					_ = purchase.TransactionID
					_ = purchase.PlayerID
					_ = purchase.AmountCents
					return nil
				})
				
				file.Close()
				
				if err != nil {
					b.Fatalf("StreamNDJSON failed: %v", err)
				}
				
				if processedCount != tt.count {
					b.Fatalf("Expected %d records, got %d", tt.count, processedCount)
				}
			}
			
			// Get file size for reporting
			fileInfo, err := os.Stat(testFile)
			if err == nil {
				fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
				b.ReportMetric(fileSizeMB, "MB/op")
			}
			
			recordsPerOp := float64(tt.count)
			b.ReportMetric(recordsPerOp, "records/op")
		})
	}
}

// BenchmarkValidatePurchaseInput tests validation performance
func BenchmarkValidatePurchaseInput(b *testing.B) {
	validInput := main.PurchaseInput{
		TransactionID:  "TXN-VALID-BENCH",
		PlayerID:       "player_bench_001",
		PlayerUsername: "BenchmarkUser",
		GameTitle:      "Benchmark Game Title",
		ItemType:       "game",
		Genre:          "Action",
		Platform:       "steam",
		AmountCents:    2999,
		Currency:       "USD",
		PlayerLevel:    50,
		CreatedAt:      "2025-08-15T10:00:00Z",
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		err := main.ValidatePurchaseInput(validInput)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}