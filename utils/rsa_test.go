package utils

import (
	"bytes"
	"crypto"
	"fmt"
	"testing"
)

func TestRSAEncryptAndDecrypt(t *testing.T) {
	// 生成公钥和私钥
	priBytes, pubBytes, err := GenerateRSAPair(1024)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	fmt.Printf("Public Key:\n%s\n", string(pubBytes))
	fmt.Printf("Private Key:\n%s\n", string(priBytes))

	// 解析公钥
	publicKey, err := ParsePublicKey(pubBytes)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}
	// 解析私钥
	privateKey, err := ParsePrivateKey(priBytes)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}

	// 加密消息-超出了1024可加密范围
	message := []byte("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
	ciphertext, err := RSAEncrypt(publicKey, message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}

	fmt.Printf("字符串:%s\n加密后:%s\n\n", string(message), string(ciphertext))

	// 解密消息
	plaintext, err := RSADecrypt(privateKey, ciphertext)
	if err != nil {
		t.Fatalf("failed to decrypt message: %v", err)
	}

	fmt.Printf("解密后:%s\n", string(plaintext))
	// 比较明文和原始消息
	if !bytes.Equal(plaintext, message) {
		t.Fatalf("decrypted message does not match original message")
	}
	ciphertextOAEP, err := RSAEncryptOAEP(publicKey, message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}
	t.Logf("OAEP加密结果:%s\n", string(ciphertextOAEP))
	plaintextOAEP, err := RSADecryptOAEP(privateKey, ciphertextOAEP)
	if err != nil {
		t.Fatalf("failed to decrypt message: %v", err)
	}
	t.Logf("OAEP解密结果:%s\n", string(plaintextOAEP))
	if !bytes.Equal(plaintextOAEP, message) {
		t.Fatalf("decrypted message does not match original message")
	}
}

func TestRSASignAndVerify(t *testing.T) {
	// 生成公钥和私钥
	priBytes, pubBytes, err := GenerateRSAPair(1024)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	fmt.Printf("Public Key:\n%s\n", string(pubBytes))
	fmt.Printf("Private Key:\n%s\n", string(priBytes))

	// 解析公钥
	publicKey, err := ParsePublicKey(pubBytes)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}
	// 解析私钥
	privateKey, err := ParsePrivateKey(priBytes)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}

	// 测试数据
	data := []byte("test data")
	// 使用私钥对数据进行签名
	signature, err := RSASign(privateKey, data, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("`%s`签名结果是:\n%s\n\n", string(data), signature)
	// 使用公钥对数据和签名进行验证
	err = RSAVerify(publicKey, data, signature, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("验签是否通过:%v\n", err == nil)
}

func TestRSAVerifyBase64Key(t *testing.T) {
	// 该函数使用base64格式的公钥
	sign := "rGWk8Y6ICV5qvJ2KhHP2QT15pNs0364wsjRmdQ9VfriQMu0izP2edzsoFPOTB34QYQCBAA7n9keeV5CzbmAAO6wmn4HLCbp0i0zlTOuw9eLHxeWp+VS/+s2tEXJb9mzGTiEHUbMx7u4xQuwnAOMnvX47I0Bv3qXEyRu1L1t7J1Y="
	pubStr := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCyldWy5iGt6jQIyOF4oxlSfg6QFEiOTPk+tDKE9l/0j4E4gEh/BQCwCjLYmemjxVyAxZoP1rdGxjXnB1V1maPZo1rQmOvEsBO//wfaV8f23aomm/pB2870YjC4sckwgsq1YKAwHnDAOglcnGOZUkJ1nKyfmCFfBlWMgLRrqsynbQIDAQAB"
	publicKey, err := ParsePublicKey([]byte(pubStr))
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}
	data := []byte("test data")
	err = RSAVerify(publicKey, data, sign, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("验签是否通过:%v\n", err == nil)
}
