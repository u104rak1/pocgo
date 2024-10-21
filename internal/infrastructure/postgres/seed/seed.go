//go:build !production

package seed

import (
	"log"

	"github.com/uptrace/bun"
)

func InsertMasterData(db *bun.DB) {
	if err := saveTransactionTypeMaster(db); err != nil {
		log.Println("Error inserting transaction type master:", err)
	}
	if err := saveCurrencyMaster(db); err != nil {
		log.Println("Error inserting currency master:", err)
	}
}

func InsertSeedData(db *bun.DB) {
	if err := saveUser(db); err != nil {
		log.Println("Error inserting user data:", err)
	}
	if err := saveAuthentication(db); err != nil {
		log.Println("Error inserting authentication data:", err)
	}
	if err := saveAccount(db); err != nil {
		log.Println("Error inserting account data:", err)
	}
	if err := saveTransaction(db); err != nil {
		log.Println("Error inserting transaction data:", err)
	}
}
