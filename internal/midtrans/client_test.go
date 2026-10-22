package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestVerifyNotification(t *testing.T) {
	const (
		orderID     = "EW-1-test"
		statusCode  = "200"
		grossAmount = "25000.00"
		serverKey   = "sandbox-server-key"
	)

	hash := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	signature := hex.EncodeToString(hash[:])

	client := NewClient(serverKey)
	if !client.VerifyNotification(orderID, statusCode, grossAmount, signature) {
		t.Fatal("expected valid notification signature to pass")
	}
	if client.VerifyNotification(orderID, statusCode, grossAmount, "invalid-signature") {
		t.Fatal("expected invalid notification signature to fail")
	}
	if NewClient("").VerifyNotification(orderID, statusCode, grossAmount, signature) {
		t.Fatal("expected missing server key to reject notification")
	}
}
