package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	// The UUIDs from the database query
	uuidStrings := []string{
		"01973569-f7d7-7bd5-b217-dc29dd412877",
		"01973569-f7d7-76ae-86f9-3956f51ef45e",
		"01973569-f7d7-722f-9f40-06602a203a7e",
	}

	fmt.Println("=== Verifying UUIDv7 from Finance Import ===")

	for i, uuidStr := range uuidStrings {
		parsedUuid, err := uuid.Parse(uuidStr)
		if err != nil {
			fmt.Printf("UUID %d: ERROR parsing %s: %v\n", i+1, uuidStr, err)
			continue
		}

		fmt.Printf("UUID %d: %s\n", i+1, uuidStr)
		fmt.Printf("  Version: %d\n", parsedUuid.Version())
		fmt.Printf("  Is UUIDv7: %t\n", parsedUuid.Version() == 7)
		fmt.Println()
	}
}
