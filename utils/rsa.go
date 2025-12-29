package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

var (
	errPublic          = errors.New("public key is not an RSA public key")
	ErrDataToLarge     = errors.New("message too long for RSA public key size")
	ErrDataLen         = errors.New("data length error")
	ErrDataBroken      = errors.New("data broken, first byte is not zero")
	ErrKeyPairDismatch = errors.New("data is not encrypted by the private key")
	ErrDecryption      = errors.New("decryption error")
	ErrPublicKey       = errors.New("get public key error")
	ErrPrivateKey      = errors.New("get private key error")
)

// decodeKey 解码key
// 先尝试用pem解码，如果失败则使用base64解码
func decodeKey(key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		return block.Bytes, nil
	}
	return base64.StdEncoding.DecodeString(string(key))
}

// GenerateRSAPair 生成私钥和公钥
func GenerateRSAPair(bits int) ([]byte, []byte, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	// 编码私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	// 获取公钥
	publicKey := privateKey.Public()
	// 编码公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	return pri, pub, nil
}

// ParsePublicKey 解析公钥-PKIX格式
func ParsePublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	data, err := decodeKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errPublic
	}
	return publicKey, nil
}

// ParsePrivateKey 解析私钥-PKCS1格式
func ParsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	data, err := decodeKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	// 解析私钥字节
	privateKey, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	// 返回私钥实例和nil
	return privateKey, nil
}

// ParsePrivateKeyCS8 解析私钥-PKCS8格式
func ParsePrivateKeyCS8(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	data, err := base64.StdEncoding.DecodeString(string(privateKeyBytes))
	if err != nil {
		return nil, err
	}
	// 解析私钥字节
	privateKey, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	// 返回私钥实例
	return rsaPrivateKey, nil
}

// RSAEncrypt RSA加密消息-返回base64编码的结果
// 如果密文长度超过 publicKey.Size() - 11，则会进行分段加密。
func RSAEncrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	chunkSize := publicKey.Size() - 11
	// 如果消息长度小于等于分段大小，则不需要分段
	if len(message) <= chunkSize {
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(encryptedChunk)), nil
	}
	// 如果消息长度大于分段大小，则需要分段
	var ciphertext []byte
	for i := 0; i < len(message); i += chunkSize {
		end := i + chunkSize
		if end > len(message) {
			end = len(message)
		}
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message[i:end])
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, encryptedChunk...)
	}
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// RSADecrypt 解密消息-密文是base64编码
// 如果密文长度超过私钥大小，则会进行分段解密。
func RSADecrypt(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// 将密文从 base64 编码转换为字节
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	// 如果数据长度小于或等于私钥大小，则直接解密并返回结果
	if len(ciphertext) <= privateKey.Size() {
		decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
		if err != nil {
			return nil, err
		}
		return decryptedBytes, nil
	}
	// 如果数据长度大于私钥大小，则分段解密
	var decryptedMessage []byte
	for i := 0; i < len(ciphertext); i += privateKey.Size() {
		decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext[i:i+privateKey.Size()])
		if err != nil {
			return nil, err
		}
		decryptedMessage = append(decryptedMessage, decryptedBytes...)
	}
	return decryptedMessage, nil // 返回解密后的消息
}

// RSASign 使用私钥对数据进行签名
// privateKey 是 RSA 私钥
// data 是要签名的数据
// hash 是要使用的哈希函数，可以是 crypto.SHA1、crypto.SHA256 或 crypto.SHA512
func RSASign(privateKey *rsa.PrivateKey, data []byte, hash crypto.Hash) (string, error) {
	// 计算数据的哈希值
	hashed := hash.New()
	hashed.Write(data)
	// 签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed.Sum(nil))
	if err != nil {
		return "", err
	}
	// 返回 base64 编码的签名字符串
	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerify 使用公钥对数据和签名进行验证
// publicKey 是 RSA 公钥
// data 是要验证的数据
// sign 是 base64 编码的签名字符串
// hash 是要使用的哈希函数，可以是 crypto.SHA1、crypto.SHA256 或 crypto.SHA512
func RSAVerify(publicKey *rsa.PublicKey, data []byte, sign string, hash crypto.Hash) error {
	// 解码签名字符串
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	// 计算数据的哈希值
	hashed := hash.New()
	hashed.Write(data)
	// 验证签名
	err = rsa.VerifyPKCS1v15(publicKey, hash, hashed.Sum(nil), signature)
	if err != nil {
		return err
	}
	return nil
}

// PubKeyDecrypt 公钥解密
func PubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, ErrDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, ErrDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, ErrDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, ErrKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

// RSAEncryptOAEP RSA加密消息-返回base64编码的结果
func RSAEncryptOAEP(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	// 使用 SHA-256 作为哈希函数
	hash := sha256.New()

	// 计算最大可加密的块大小
	// 对于 OAEP，最大块大小是 k - 2*hashLen - 2，其中 k 是密钥大小（字节）
	chunkSize := publicKey.Size() - 2*hash.Size() - 2

	// 如果消息长度小于等于分段大小，则不需要分段
	if len(message) <= chunkSize {
		encryptedChunk, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, message, nil)
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(encryptedChunk)), nil
	}

	// 如果消息长度大于分段大小，则需要分段
	var ciphertext []byte
	for i := 0; i < len(message); i += chunkSize {
		end := i + chunkSize
		if end > len(message) {
			end = len(message)
		}
		encryptedChunk, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, message[i:end], nil)
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, encryptedChunk...)
	}
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// RSADecryptOAEP RSA解密消息-密文是base64编码
func RSADecryptOAEP(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// 使用 SHA-256 作为哈希函数
	hash := sha256.New()

	// 将密文从 base64 编码转换为字节
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	// 计算每段加密块的大小（与加密时一致）
	chunkSize := privateKey.Size() // OAEP 加密后的块大小等于密钥大小

	// 如果数据长度小于或等于块大小，则直接解密
	if len(ciphertext) <= chunkSize {
		decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext, nil)
		if err != nil {
			return nil, err
		}
		return decryptedBytes, nil
	}

	// 如果数据长度大于块大小，则分段解密
	var decryptedMessage []byte
	for i := 0; i < len(ciphertext); i += chunkSize {
		end := i + chunkSize
		if end > len(ciphertext) {
			end = len(ciphertext)
		}

		decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext[i:end], nil)
		if err != nil {
			return nil, err
		}
		decryptedMessage = append(decryptedMessage, decryptedBytes...)
	}
	return decryptedMessage, nil
}
