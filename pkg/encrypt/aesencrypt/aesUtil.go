package aesencrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"file-transfer/pkg/common"
	"file-transfer/pkg/log"
	"sync"

	"github.com/spf13/viper"
)

var (
	key string
	iv  string

	once sync.Once
)

func InitAES() {
	once.Do(func() {
		if len(key) < 1 {
			key = viper.GetString(common.VIPER_AES_KEY)
			iv = viper.GetString(common.VIPER_AES_IV)
			if len(key) < 1 || len(iv) < 1 {
				log.Errorw("AES key/iv read failed, some functions could not use")
			}
		}
	})
}

// GetAESDecrypted decrypts given text in AES 256 CBC
func GetAESDecrypted(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pKCS5UnPadding(ciphertext)

	return string(ciphertext), nil
}

// pKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string) (string, error) {
	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}
