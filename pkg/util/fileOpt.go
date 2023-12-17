package util

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func CalculateFileSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	sha1String := fmt.Sprintf("%x", hashBytes)

	return sha1String, nil
}

func CreateDirectoryIfNotExists(path string) error {
	// Check if the directory exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Some error occurred while checking the directory
		return err
	}

	return nil
}
