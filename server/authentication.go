package server

import (
	dbutils "DistriAI-Node/utils/db_utils"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
)

func UserAuthentication(db *dbutils.DB, validityPeriod int64, signature string, message string) (bool, error) {
	buyerPublicKey, err := db.Get([]byte("buyer"))
	if err != nil {
		return false, fmt.Errorf("> Get error: %v", err)
	}

	publicKeyStr := string(buyerPublicKey)
	publicKey, err := solana.PublicKeyFromBase58(publicKeyStr)
	if err != nil {
		return false, fmt.Errorf("> PublicKeyFromBase58 error: %v", err)
	}

	out, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return false, fmt.Errorf("> SignatureFromBase58 error: %v", err)
	}

	currentTime := time.Now().Unix() / validityPeriod
	msg := fmt.Sprintf("%s/%v/%s", message, currentTime, publicKeyStr)

	if publicKey.Verify([]byte(msg), out) {
		return true, nil
	} else {
		currentTime -= 1
		msg = fmt.Sprintf("%s/%v/%s", message, currentTime, publicKeyStr)

		if publicKey.Verify([]byte(msg), out) {
			return true, nil
		}
	}

	return false, nil
}
