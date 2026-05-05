package utils

import (
	"fmt"
	"log"
)

func GenerateQRValue(qrFor string, qrIds string, createdBy string) string {
	// Create the data string using qr_for, qr_ids, and created_by
	dataString := fmt.Sprintf("*()*", qrFor, qrIds, createdBy)

	// Encrypt the data string
	encryptedQrValue, err := Encrypt(dataString)
	if err != nil {
		log.Printf("Error encrypting QR value: %v", err)
		return ""
	}

	return encryptedQrValue
}
