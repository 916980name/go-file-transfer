package rsaencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"file-transfer/pkg/common"
	"file-transfer/pkg/log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	once sync.Once
)

func InitRSAKeyPair() {
	once.Do(func() {
		if privateKey == nil {
			privateKeyFilePath := viper.GetString(common.VIPER_RSA_PRI)
			publicKeyFilePath := viper.GetString(common.VIPER_RSA_PUB)
			if len(privateKeyFilePath) < 1 || len(publicKeyFilePath) < 1 {
				log.Errorw("RSA keypair read failed, some functions could not use")
			}
			err := readKeyFile(privateKeyFilePath, publicKeyFilePath)
			if err != nil {
				log.Errorw("RSA keypair parse failed, some functions could not use", "error", err)
			}
		}
	})
}

func readKeyFile(privateKeyFile string, publicKeyFile string) error {
	// Read the private key file
	privateKeyData, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return errors.New("Error reading private key file:" + privateKeyFile)
	}

	// Decode the PEM data
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return errors.New("invalid private key file")
	}

	// Parse the RSA private key
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("error parsing private key")
	}

	// Read the public key file
	publicKeyData, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return errors.New("error reading public key file:" + publicKeyFile)
	}

	// Decode the PEM data
	block, _ = pem.Decode(publicKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return errors.New("invalid public key file")
	}

	// Parse the RSA public key
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return errors.New("error parsing public key")
	}
	return nil
}

func Encrypt(data string) (string, error) {
	// Example usage: Encrypt and decrypt with the keys
	originalData := []byte(data)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, originalData)
	if err != nil {
		return "", err
	}
	return string(encryptedData), nil
}

func Decrypt(data string) (string, error) {
	encryptedData := []byte(data)
	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedData)
	if err != nil {
		return "", err
	}
	return string(decryptedData), nil
}
