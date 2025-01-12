package main

import (
	"fmt"
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/database"
)

func main() {
	// Open the key-value store database
	kvStore, err := database.OpenKVStore("bitcoin.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer kvStore.Close()

	// User menu
	for {
		fmt.Println("\n=== Blockchain Menu ===")
		fmt.Println("1. View blockchain")
		fmt.Println("2. Create a wallet")
		fmt.Println("3. Add a transaction")
		fmt.Println("4. Mine a block")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		_, err := fmt.Scan(&choice)
		fmt.Println("")
		if err != nil {
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		switch choice {
		case 1:
		case 2:
		case 3:
		case 4:
		case 5:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
