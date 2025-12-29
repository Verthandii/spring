package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestAesEncryptAndDecrypt(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	key := RandStr(32)
	iv := RandStr(16)
	fmt.Printf("IV: %s\n", iv)
	fmt.Printf("Key: %s\n", key)
	text := "110101201003075293"
	fmt.Printf("加密前: %s\n", text)

	encrypted, err := AesEncrypt(text, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("加密后: %s\n", encrypted)
	decrypted, err := AesDecrypt(encrypted, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("解密后: %s\n", decrypted)
}
