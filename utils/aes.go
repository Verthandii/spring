package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

// PKCS7Padding 是一个用来执行PKCS7填充的函数.
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS7UnPadding 是一个用于移除PKCS7填充的函数。
func PKCS7UnPadding(plaintext []byte) []byte {
	length := len(plaintext)
	unPadding := int(plaintext[length-1])
	return plaintext[:(length - unPadding)]
}

// AesEncrypt AES加密 接收明文、密钥和初始化向量作为输入，返回经过AES算法加密并经过hex编码后的密文
func AesEncrypt(text string, key string, iv string) (string, error) {
	plaintext := PKCS7Padding([]byte(text), aes.BlockSize)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	// 检查明文长度是否是块大小的整数倍
	if len(plaintext)%aes.BlockSize != 0 {
		return "", errors.New("明文长度必须是块大小的整数倍")
	}
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plaintext)
	return hex.EncodeToString(ciphertext), nil
}

// AesDecrypt AES解密 接收密文、密钥和初始化向量作为输入，返回解密后的明文
func AesDecrypt(encryptedText string, key string, iv string) (string, error) {
	ciphertext, _ := hex.DecodeString(encryptedText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	// 检查密文长度是否是块大小的整数倍
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("密文长度必须是块大小的整数倍")
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS7UnPadding(plaintext)
	return string(plaintext), nil
}
