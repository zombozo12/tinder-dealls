package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"io"
)

func ExtractUserClaims(ctx *fiber.Ctx, key string) (user User, err error) {
	jwtUser := ctx.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	userClaims := claims["sub"].(string)

	userDecrypted, err := DecryptAESWithGCM(userClaims, key)
	if err != nil {
		log.Debugf("error: %s", err)
		return User{}, err
	}

	if err := json.Unmarshal([]byte(userDecrypted), &user); err != nil {
		log.Debugf("error: %s", err)
		return User{}, err
	}

	return user, nil
}

func EncryptAESWithGCM(stringToEncrypt, keyString string) (encryptedString string, err error) {
	// Convert the key from string to bytes
	key := []byte(keyString)
	plaintext := []byte(stringToEncrypt)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data
	// The first nonce argument in Seal is the prefix
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAESWithGCM(encryptedString, keyString string) (decryptedString string, err error) {
	// Convert the key from string to bytes
	key := []byte(keyString)
	enc := []byte(encryptedString)

	decryptedEnc, err := base64.StdEncoding.DecodeString(string(enc))
	if err != nil {
		return "", err
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := decryptedEnc[:nonceSize], decryptedEnc[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
