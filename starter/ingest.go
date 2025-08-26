package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)


// StreamNDJSON parses newline-delimited JSON and calls fn for each purchase
func StreamNDJSON(ctx context.Context, r io.Reader, fn func(Purchase) error) error {
	// TODO: Implement this.
	
	return nil
}



// ValidatePurchaseInput performs basic validation on purchase input
func ValidatePurchaseInput(input PurchaseInput) error {
	// TODO: Add validation rules

	return nil
}