package chain

import (
	"encoding/hex"
	"fmt"
)

// Display отображает все блоки
func (bc *Blockchain) Display() {
	for i, b := range bc.Blocks {
		fmt.Printf("--- Block ID: %d ---\n", i)
		fmt.Printf("Timestamp: 	%d\n", b.Header.Timestamp)
		fmt.Printf("Transactions: 	%d\n", len(b.Transaction))
		fmt.Printf("Prev Hash:  %s\n", hex.EncodeToString(b.Header.PreviousHash))
		fmt.Printf("Size: %d bytes\n", b.Size)
		fmt.Printf("Hash: 		%s\n", hex.EncodeToString(b.Hash))
		fmt.Println("  --- ฿ ---  ")
	}
}
